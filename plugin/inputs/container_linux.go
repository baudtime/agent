package inputs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/baudtime/agent/plugin"
	"github.com/baudtime/agent/vars"
	"github.com/baudtime/baudtime/msg"
	"github.com/go-kit/kit/log/level"
	"github.com/mdlayher/taskstats"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	containerCpuUsagePercentMetric      = "cpu_usage_percent"
	containerMemoryUsagePercentMetric   = "mem_usage_percent"
	containerSwapUsagePercentMetric     = "swap_usage_percent"
	containerSwapUsageBytesMetric       = "swap_usage_bytes"
	containerProcessRunningMetric       = "proc_running_num"
	containerLoadAvgMetric              = "load_avg"
	containerRxBytesMetric              = "net_dev_rx_bytes"
	containerRxPacketsMetric            = "net_dev_rx_packets"
	containerRxErrsMetric               = "net_dev_rx_errs"
	containerRxDropMetric               = "net_dev_rx_drop"
	containerTxBytesMetric              = "net_dev_tx_bytes"
	containerTxPacketsMetric            = "net_dev_tx_packets"
	containerTxErrsMetric               = "net_dev_tx_errs"
	containerTxDropMetric               = "net_dev_tx_drop"
	containerTcpConnNumMetric           = "tcp_conn_num"
	containerTcpRetransMetric           = "tcp_retrans"
	containerTcpActiveOpensMetric       = "tcp_active_opens"
	containerReadCountMetric            = "disk_read_count"
	containerWriteCountMetric           = "disk_write_count"
	containerReadKBytesMetric           = "disk_read_kb"
	containerWriteKBytesMetric          = "disk_write_kb"
	containerIopsInProgressMetric       = "disk_iops_in_progress"
	containerIoTimeMetric               = "disk_io_time"
	containerBusyPercentMetric          = "disk_busy_percent"
	containerFsUsagePercentMetric       = "fs_usage_percent"
	containerFsUsageBytesMetric         = "fs_usage_bytes"
	containerFsInodesUsagePercentMetric = "fs_inodes_usage_percent"
	containerFsInodesUsageNumMetric     = "fs_inodes_usage_num"
	containerNtpClockOffsetMilliMetric  = "ntp_clk_offset_ms"
)

var (
	runtimeRoot        = kingpin.Flag("input.container.oci-root", "root path of OCI runtime").Default("/var/run/docker/execdriver/native").String()
	containerPlgLogger = plugin.Logger("container")
)

func init() {
	register("container", plugin.DefaultDisabled, newContainerCollector)
}

type hostInfo struct {
	IP       string
	HostName string
}

type containerStats struct { // each container has a containerStats
	basicStats  *cgroups.Stats
	netDevStats map[string]NetDevStats
	snmpStats   *Snmp
	diskIOStats *DiskIOStats
	loadAvg     float64
	timestamp   time.Time
}

type containerInfo struct {
	hostInfo
	containerStats
}

type containerCollector struct {
	factory                     libcontainer.Factory
	containerInfos              map[string]containerInfo
	ignoredNetDevicesPattern    *regexp.Regexp
	ignoredDiskDevicesPattern   *regexp.Regexp
	ignoredFsMountPointsPattern *regexp.Regexp
	ignoredFsTypesPattern       *regexp.Regexp
}

func newContainerCollector() (plugin.Input, error) {
	if _, err := os.Stat(*runtimeRoot); err != nil {
		return nil, err
	}
	factory, err := libcontainer.New(*runtimeRoot, libcontainer.Cgroupfs)
	if err != nil {
		return nil, err
	}
	return &containerCollector{
		factory:                     factory,
		containerInfos:              make(map[string]containerInfo),
		ignoredNetDevicesPattern:    regexp.MustCompile(*netDevIgnoredDevices),
		ignoredDiskDevicesPattern:   regexp.MustCompile("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme)\\.+$"),
		ignoredFsMountPointsPattern: regexp.MustCompile(*ignoredFsMountPoints),
		ignoredFsTypesPattern:       regexp.MustCompile(*ignoredFsTypes),
	}, nil
}

func (c *containerCollector) getAllContainers() ([]libcontainer.Container, error) {
	entries, err := ioutil.ReadDir(*runtimeRoot)
	if err != nil {
		return nil, err
	}

	containers := make([]libcontainer.Container, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			container, err := c.factory.Load(entry.Name())
			if err != nil {
				containerPlgLogger(level.Error).Log("err", err)
				continue
			}
			containers = append(containers, container)
		}
	}

	return containers, nil
}

func (c *containerCollector) Collect(ch chan<- plugin.Metric) error {
	containers, err := c.getAllContainers()
	if err != nil {
		return err
	}

	var multiErr error

	netlinkCli, err := taskstats.New()
	if err != nil {
		multiErr = multierr.Append(multiErr, err)
	}
	defer netlinkCli.Close()

	now := time.Now()

	ntpClockOffset, err := getNtpClockOffset()
	if err != nil {
		multiErr = multierr.Append(multiErr, err)
	}

	allFileSystemStats, err := getFileSystemStats(c.ignoredFsMountPointsPattern, c.ignoredFsTypesPattern)
	if err != nil {
		containerPlgLogger(level.Error).Log("err", err)
		multiErr = multierr.Append(multiErr, err)
	}

	allDiskIOStats, err := readDiskIOStats(procFilePath("diskstats"), c.ignoredDiskDevicesPattern)
	if err != nil {
		containerPlgLogger(level.Error).Log("err", err)
		multiErr = multierr.Append(multiErr, err)
	}

	for _, container := range containers {
		id := container.ID()
		mountPoint := getMountPoint(container)

		state, err := container.State()
		if err != nil {
			multiErr = multierr.Append(multiErr, err)
			continue
		}

		nativeStats, err := container.Stats()
		if err != nil {
			multiErr = multierr.Append(multiErr, err)
			continue
		}

		netDevStats, err := readNetDevStat(fmt.Sprintf("/proc/%d/net/dev", state.InitProcessPid), c.ignoredNetDevicesPattern)
		if err != nil {
			containerPlgLogger(level.Error).Log("err", err)
			multiErr = multierr.Append(multiErr, err)
		}

		snmpStats, err := readSnmp(fmt.Sprintf("/proc/%d/net/snmp", state.InitProcessPid))
		if err != nil {
			containerPlgLogger(level.Error).Log("err", err)
			multiErr = multierr.Append(multiErr, err)
		}

		var diskIOStats *DiskIOStats
		filesystemStats, found := allFileSystemStats[mountPoint]
		if found {
			if devShortName, err := getDevShortName(filesystemStats.labels.device); err != nil {
				containerPlgLogger(level.Error).Log("err", err)
				multiErr = multierr.Append(multiErr, err)
			} else {
				diskIOStats = allDiskIOStats[devShortName]
			}
		}

		stats := containerStats{
			basicStats:  nativeStats.CgroupStats,
			netDevStats: netDevStats,
			snmpStats:   snmpStats,
			diskIOStats: diskIOStats,
			loadAvg:     -1,
			timestamp:   now,
		}

		info, found := c.containerInfos[id]
		if !found {
			if hInfo, err := getHostInfo(container); err != nil {
				containerPlgLogger(level.Error).Log("err", err)
				multiErr = multierr.Append(multiErr, err)
			} else {
				c.containerInfos[id] = containerInfo{
					hostInfo:       hInfo,
					containerStats: stats,
				}
			}
			continue
		}

		preBasicStats := info.basicStats
		preNetStats := info.netDevStats
		preSnmpStats := info.snmpStats
		preDiskIOStats := info.diskIOStats
		elapsedSec := stats.timestamp.Sub(info.timestamp).Seconds()

		labels := []msg.Label{
			{"ip", info.hostInfo.IP},
			{"hostname", info.hostInfo.HostName},
			msg.Label{"parent", vars.LocalIP},
		}

		collectCPUStats(ch, stats.basicStats, preBasicStats, elapsedSec, state, labels...)
		collectMemStats(ch, stats.basicStats, labels...)
		if processes, err := container.Processes(); err == nil {
			collectProcessNum(ch, len(processes), labels...)
		}
		collectNetDevStats(ch, netDevStats, preNetStats, elapsedSec, labels...)
		collectTcpStats(ch, snmpStats, preSnmpStats, elapsedSec, labels...)
		collectDiskIOStats(ch, diskIOStats, preDiskIOStats, elapsedSec, labels...)
		collectFilesystemStats(ch, filesystemStats, labels...)
		collectNtpClockOffsetStats(ch, ntpClockOffset, labels...)

		if netlinkStats, err := netlinkCli.CGroupStats(state.CgroupPaths["cpu"]); err == nil {
			if info.loadAvg < 0 {
				stats.loadAvg = float64(netlinkStats.Running)
			} else {
				loadDecay := math.Exp(float64(-elapsedSec / 60))
				stats.loadAvg = info.loadAvg*loadDecay + float64(netlinkStats.Running)*(1.0-loadDecay)
			}
			collectLoadAvg(ch, stats.loadAvg, labels...)
		}

		c.containerInfos[id] = containerInfo{
			hostInfo:       info.hostInfo,
			containerStats: stats,
		}
	}

	for id, lastStats := range c.containerInfos {
		if time.Since(lastStats.timestamp) > 2*time.Minute {
			delete(c.containerInfos, id)
		}
	}

	return multiErr
}

func collectCPUStats(ch chan<- plugin.Metric, stats, preStats *cgroups.Stats, elapsedSec float64, state *libcontainer.State, labels ...msg.Label) error {
	if stats == nil || preStats == nil {
		return errors.New("cpu stats is nil")
	}
	cpuLimits := len(stats.CpuStats.CpuUsage.PercpuUsage)

	cpuQuota, err := getCgroupParamUint(state.CgroupPaths["cpu"], "cpu.cfs_quota_us")
	if err != nil {
		return err
	}

	cpuPeriod, err := getCgroupParamUint(state.CgroupPaths["cpu"], "cpu.cfs_period_us")
	if err != nil {
		return err
	}

	if cpuQuota != 0 && cpuPeriod != 0 {
		cpuLimits = int(cpuQuota / cpuPeriod)
	} else if len(state.Config.Cgroups.CpusetCpus) > 0 {
		cpuLimits = len(state.Config.Cgroups.CpusetCpus)
	}

	cpuTotalDelta := float64(stats.CpuStats.CpuUsage.TotalUsage - preStats.CpuStats.CpuUsage.TotalUsage)

	ch <- plugin.Metric{Name: containerCpuUsagePercentMetric, Value: cpuTotalDelta * 100.0 / float64(time.Second) / elapsedSec / float64(cpuLimits)}.With(labels...)
	return nil
}

func collectMemStats(ch chan<- plugin.Metric, stats *cgroups.Stats, labels ...msg.Label) error {
	if stats == nil {
		return errors.New("memory stats is nil")
	}

	if stats != nil {
		if stats.MemoryStats.Usage.Limit > 0 {
			memPercent := float64(stats.MemoryStats.Stats["active_anon"]+stats.MemoryStats.Stats["inactive_anon"]) /
				float64(stats.MemoryStats.Usage.Limit) * 100.0
			ch <- plugin.Metric{Name: containerMemoryUsagePercentMetric, Value: memPercent}.With(labels...)
		}
		if stats.MemoryStats.SwapUsage.Limit > 0 {
			swapPercent := float64(stats.MemoryStats.Stats["swap"]) / float64(stats.MemoryStats.SwapUsage.Limit) * 100.0
			ch <- plugin.Metric{Name: containerSwapUsagePercentMetric, Value: swapPercent}.With(labels...)
		}
		ch <- plugin.Metric{Name: containerSwapUsageBytesMetric, Value: float64(stats.MemoryStats.Stats["swap"])}.With(labels...)
	}
	return nil
}

func collectNetDevStats(ch chan<- plugin.Metric, stats, preStats map[string]NetDevStats, elapsedSec float64, labels ...msg.Label) error {
	if stats == nil || preStats == nil {
		return errors.New("net dev stats is nil")
	}
	for _, devStat := range stats {
		if lastStat, found := preStats[devStat.Iface]; found {
			ifaceLabel := msg.Label{"iface", devStat.Iface}

			rxBytes := uint64(float64(devStat.RxBytes-lastStat.RxBytes) / elapsedSec)
			rxPackets := uint64(float64(devStat.RxPackets-lastStat.RxPackets) / elapsedSec)
			rxErrs := uint64(float64(devStat.RxErrs-lastStat.RxErrs) / elapsedSec)
			rxDrop := uint64(float64(devStat.RxDrop-lastStat.RxDrop) / elapsedSec)
			txBytes := uint64(float64(devStat.TxBytes-lastStat.TxBytes) / elapsedSec)
			txPackets := uint64(float64(devStat.TxPackets-lastStat.TxPackets) / elapsedSec)
			txErrs := uint64(float64(devStat.TxErrs-lastStat.TxErrs) / elapsedSec)
			txDrop := uint64(float64(devStat.TxDrop-lastStat.TxDrop) / elapsedSec)

			ch <- plugin.Metric{Name: containerRxBytesMetric, Value: float64(rxBytes)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerRxPacketsMetric, Value: float64(rxPackets)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerRxErrsMetric, Value: float64(rxErrs)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerRxDropMetric, Value: float64(rxDrop)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerTxBytesMetric, Value: float64(txBytes)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerTxPacketsMetric, Value: float64(txPackets)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerTxErrsMetric, Value: float64(txErrs)}.With(labels...).With(ifaceLabel)
			ch <- plugin.Metric{Name: containerTxDropMetric, Value: float64(txDrop)}.With(labels...).With(ifaceLabel)
		}
	}
	return nil
}

func collectTcpStats(ch chan<- plugin.Metric, stats, preStats *Snmp, elapsedSec float64, labels ...msg.Label) error {
	if stats == nil || preStats == nil {
		return errors.New("tcp stats is nil")
	}
	tcpConnNum := stats.TcpCurrEstab
	tcpRetrans := uint64(float64(stats.TcpRetransSegs-preStats.TcpRetransSegs) / elapsedSec)
	tcpActiveOpens := uint32(float64(stats.TcpActiveOpens-preStats.TcpActiveOpens) / elapsedSec)

	ch <- plugin.Metric{Name: containerTcpConnNumMetric, Value: float64(tcpConnNum)}.With(labels...)
	ch <- plugin.Metric{Name: containerTcpRetransMetric, Value: float64(tcpRetrans)}.With(labels...)
	ch <- plugin.Metric{Name: containerTcpActiveOpensMetric, Value: float64(tcpActiveOpens)}.With(labels...)
	return nil
}

func collectDiskIOStats(ch chan<- plugin.Metric, stats, preStats *DiskIOStats, elapsedSec float64, labels ...msg.Label) error {
	if stats == nil || preStats == nil {
		return errors.New("disk io stats is nil")
	}
	readCount := uint32(float64(stats.ReadCount-preStats.ReadCount) / elapsedSec)
	writeCount := uint32(float64(stats.WriteCount-preStats.WriteCount) / elapsedSec)
	readKBytes := uint32((float64(stats.ReadBytes-preStats.ReadBytes) / elapsedSec) / 1024)
	writeKBytes := uint32((float64(stats.WriteBytes-preStats.WriteBytes) / elapsedSec) / 1024)
	iopsInProgress := uint32(stats.IopsInProgress)
	ioTime := uint32(float64(stats.IoTime-preStats.IoTime) / elapsedSec)
	busyPercent := uint32(ioTime / 10)
	if busyPercent > 100 {
		busyPercent = 100
	}

	ch <- plugin.Metric{Name: containerReadCountMetric, Value: float64(readCount)}.With(labels...)
	ch <- plugin.Metric{Name: containerWriteCountMetric, Value: float64(writeCount)}.With(labels...)
	ch <- plugin.Metric{Name: containerReadKBytesMetric, Value: float64(readKBytes)}.With(labels...)
	ch <- plugin.Metric{Name: containerWriteKBytesMetric, Value: float64(writeKBytes)}.With(labels...)
	ch <- plugin.Metric{Name: containerIopsInProgressMetric, Value: float64(iopsInProgress)}.With(labels...)
	ch <- plugin.Metric{Name: containerIoTimeMetric, Value: float64(ioTime)}.With(labels...)
	ch <- plugin.Metric{Name: containerBusyPercentMetric, Value: float64(busyPercent)}.With(labels...)

	return nil
}

func collectFilesystemStats(ch chan<- plugin.Metric, stats *FilesystemStats, labels ...msg.Label) error {
	if stats == nil || stats.statfs == nil {
		return errors.New("fs stats is nil")
	}

	fsLabels := []msg.Label{
		{"device", stats.labels.device},
		//{"mount", stats.labels.mountPoint},
		{"type", stats.labels.fsType},
	}

	totalRoot := stats.statfs.Blocks * uint64(stats.statfs.Bsize)
	free := stats.statfs.Bfree * uint64(stats.statfs.Bsize)
	avail := stats.statfs.Bavail * uint64(stats.statfs.Bsize)
	used := totalRoot - free //so, totalRoot = used + free, totalUser = used + avail

	totalUser := used + avail
	if totalUser != 0 {
		usagePercent := (used * 100) / totalUser
		if used%totalUser != 0 {
			usagePercent += 1
		}
		if usagePercent > 100 {
			usagePercent = 100
		}
		ch <- plugin.Metric{Name: containerFsUsagePercentMetric, Value: float64(usagePercent)}.With(labels...).With(fsLabels...)
	}
	ch <- plugin.Metric{Name: containerFsUsageBytesMetric, Value: float64(used)}.With(labels...).With(fsLabels...)

	inodesTotal := stats.statfs.Files
	inodesUsed := stats.statfs.Files - stats.statfs.Ffree
	if inodesTotal != 0 {
		inodesUsagePercent := (inodesUsed * 100) / inodesTotal
		if inodesUsed%inodesTotal != 0 {
			inodesUsagePercent += 1
		}
		if inodesUsagePercent > 100 {
			inodesUsagePercent = 100
		}
		ch <- plugin.Metric{Name: containerFsInodesUsagePercentMetric, Value: float64(inodesUsagePercent)}.With(labels...).With(fsLabels...)
	}
	ch <- plugin.Metric{Name: containerFsInodesUsageNumMetric, Value: float64(inodesUsed)}.With(labels...).With(fsLabels...)

	return nil
}

func collectNtpClockOffsetStats(ch chan<- plugin.Metric, offset int64, labels ...msg.Label) error {
	ch <- plugin.Metric{Name: containerNtpClockOffsetMilliMetric, Value: float64(offset)}.With(labels...)
	return nil
}

func collectProcessNum(ch chan<- plugin.Metric, ProcessNum int, labels ...msg.Label) error {
	ch <- plugin.Metric{Name: containerProcessRunningMetric, Value: float64(ProcessNum)}.With(labels...)
	return nil
}

func collectLoadAvg(ch chan<- plugin.Metric, loadAvg float64, labels ...msg.Label) error {
	ch <- plugin.Metric{Name: containerLoadAvgMetric, Value: loadAvg}.With(labels...)
	return nil
}

func getHostInfo(container libcontainer.Container) (hostInfo, error) {
	state, err := container.State()
	if err != nil {
		return hostInfo{}, err
	}

	hostName := state.Config.Hostname

	out, err := nsenter(state.InitProcessPid, "-F -- ip -o -f inet -4 addr show scope global", true, false, false)
	if err != nil {
		return hostInfo{}, err
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		ipDesc := strings.Fields(scanner.Text())
		if len(ipDesc) < 4 {
			continue
		}

		ipAndMask := strings.Split(ipDesc[3], "/")
		ip := ipAndMask[0]

		return hostInfo{
			IP:       ip,
			HostName: hostName,
		}, nil
	}

	return hostInfo{}, errors.New("can not get ip address")
}

func getMountPoint(container libcontainer.Container) string {
	mounts := container.Config().Mounts
	for _, mount := range mounts {
		if mount.Destination == "/export" {
			return mount.Source
		}
	}
	return ""
}

func getDevShortName(devMountPoint string) (string, error) {
	devPath, err := os.Readlink(devMountPoint)
	if err != nil {
		return "", err
	}
	return path.Base(devPath), nil
}

// Gets a single uint64 value from the specified cgroup file.
func getCgroupParamUint(cgroupPath, cgroupFile string) (uint64, error) {
	fileName := filepath.Join(cgroupPath, cgroupFile)
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	res, err := parseUint(strings.TrimSpace(string(contents)), 10, 64)
	if err != nil {
		return res, fmt.Errorf("unable to parse %q as a uint from Cgroup file %q", string(contents), fileName)
	}
	return res, nil
}

func parseUint(s string, base, bitSize int) (uint64, error) {
	value, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil && intErr.(*strconv.NumError).Err == strconv.ErrRange && intValue < 0 {
			return 0, nil
		}

		return value, err
	}

	return value, nil
}

func nsenter(pid int, arg string, withNet, withDisk, withPid bool) (*bytes.Buffer, error) {
	var args []string

	if withNet {
		args = append(args, fmt.Sprintf("--net=/proc/%d/ns/net", pid))
	}
	if withDisk {
		args = append(args, fmt.Sprintf("--mount=/proc/%d/ns/mnt", pid))
	}
	if withPid {
		args = append(args, fmt.Sprintf("--pid=/proc/%d/ns/pid", pid))
	}
	args = append(args, strings.Fields(arg)...)

	cmd := exec.Command("/usr/bin/nsenter", args...)

	out := bytes.NewBuffer(make([]byte, 0, 64))
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return out, nil
}
