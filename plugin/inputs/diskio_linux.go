package inputs

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/baudtime/agent/plugin"
	. "github.com/baudtime/agent/vars"
	"github.com/baudtime/baudtime/msg"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	diskSectorSize       = 512
	readCountMetric      = "disk_read_count"
	writeCountMetric     = "disk_write_count"
	readKBytesMetric     = "disk_read_kb"
	writeKBytesMetric    = "disk_write_kb"
	iopsInProgressMetric = "disk_iops_in_progress"
	ioTimeMetric         = "disk_io_time"
	busyPercentMetric    = "disk_busy_percent"
)

var (
	diskIOIgnoredDevices = kingpin.Flag("input.diskio.ignored-devices", "Regexp of devices to ignore for diskstats.").
				Default("^(ram|loop|fd|(v|xv)d[a-z]|nvme\\d+n\\d+p|dm\\-)\\d+$").String() //^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p|dm\\-)\\d+$
	diskioPlgLogger = plugin.Logger("diskio")
)

func init() {
	register("diskio", plugin.DefaultDisabled, newDiskIOStatsCollector)
}

type diskIOStatsCollector struct {
	ignoredDevicesPattern *regexp.Regexp
	lastStats             map[string]*DiskIOStats
	lastCollectTime       time.Time
}

func newDiskIOStatsCollector() (plugin.Input, error) {
	return &diskIOStatsCollector{
		ignoredDevicesPattern: regexp.MustCompile(*diskIOIgnoredDevices),
		lastStats:             make(map[string]*DiskIOStats),
	}, nil
}

func (c *diskIOStatsCollector) Collect(ch chan<- plugin.Metric) error {
	diskIOStats, err := readDiskIOStats(procFilePath("diskstats"), c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %s", err)
	}

	now := time.Now()
	elapsedSec := now.Sub(c.lastCollectTime).Seconds()

	for i, diskIOStat := range diskIOStats {
		devLabel := msg.Label{"device", diskIOStats[i].Name}

		if lastStat, found := c.lastStats[diskIOStat.Name]; found {
			readCount := int64(float64(diskIOStat.ReadCount-lastStat.ReadCount) / elapsedSec)
			writeCount := int64(float64(diskIOStat.WriteCount-lastStat.WriteCount) / elapsedSec)
			readKBytes := float64(diskIOStat.ReadBytes-lastStat.ReadBytes) / elapsedSec / 1024
			writeKBytes := float64(diskIOStat.WriteBytes-lastStat.WriteBytes) / elapsedSec / 1024
			ioTime := float64(diskIOStat.IoTime-lastStat.IoTime) / elapsedSec
			busyPercent := uint32(ioTime / 10)
			if busyPercent > 100 {
				busyPercent = 100
			}

			ch <- plugin.Metric{Name: readCountMetric, Value: float64(readCount)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: writeCountMetric, Value: float64(writeCount)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: readKBytesMetric, Value: float64(readKBytes)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: writeKBytesMetric, Value: float64(writeKBytes)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: iopsInProgressMetric, Value: float64(diskIOStat.IopsInProgress)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: ioTimeMetric, Value: float64(ioTime)}.With(devLabel).With(HostLabels...)
			ch <- plugin.Metric{Name: busyPercentMetric, Value: float64(busyPercent)}.With(devLabel).With(HostLabels...)
		}
	}

	c.lastStats = diskIOStats
	c.lastCollectTime = now

	return nil
}

type DiskIOStats struct {
	Name             string `json:"name"`
	ReadCount        uint64 `json:"readCount"`
	MergedReadCount  uint64 `json:"mergedReadCount"`
	WriteCount       uint64 `json:"writeCount"`
	MergedWriteCount uint64 `json:"mergedWriteCount"`
	ReadBytes        uint64 `json:"readBytes"`
	WriteBytes       uint64 `json:"writeBytes"`
	ReadTime         uint64 `json:"readTime"`
	WriteTime        uint64 `json:"writeTime"`
	IopsInProgress   uint64 `json:"iopsInProgress"`
	IoTime           uint64 `json:"ioTime"`
	WeightedIO       uint64 `json:"weightedIO"`
}

func readDiskIOStats(path string, ignore *regexp.Regexp) (map[string]*DiskIOStats, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	results := make(map[string]*DiskIOStats, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			// malformed line in /proc/diskstats, avoid panic by ignoring.
			continue
		}

		name := fields[2]
		if ignore.MatchString(name) {
			diskioPlgLogger(level.Debug).Log("msg", "ignoring device", "dev", name)
			continue
		}

		reads, err := strconv.ParseUint((fields[3]), 10, 64)
		if err != nil {
			return results, err
		}
		mergedReads, err := strconv.ParseUint((fields[4]), 10, 64)
		if err != nil {
			return results, err
		}
		rbytes, err := strconv.ParseUint((fields[5]), 10, 64)
		if err != nil {
			return results, err
		}
		rtime, err := strconv.ParseUint((fields[6]), 10, 64)
		if err != nil {
			return results, err
		}
		writes, err := strconv.ParseUint((fields[7]), 10, 64)
		if err != nil {
			return results, err
		}
		mergedWrites, err := strconv.ParseUint((fields[8]), 10, 64)
		if err != nil {
			return results, err
		}
		wbytes, err := strconv.ParseUint((fields[9]), 10, 64)
		if err != nil {
			return results, err
		}
		wtime, err := strconv.ParseUint((fields[10]), 10, 64)
		if err != nil {
			return results, err
		}
		iopsInProgress, err := strconv.ParseUint((fields[11]), 10, 64)
		if err != nil {
			return results, err
		}
		iotime, err := strconv.ParseUint((fields[12]), 10, 64)
		if err != nil {
			return results, err
		}
		weightedIO, err := strconv.ParseUint((fields[13]), 10, 64)
		if err != nil {
			return results, err
		}

		results[name] = &DiskIOStats{
			Name:             name,
			ReadBytes:        rbytes * diskSectorSize,
			WriteBytes:       wbytes * diskSectorSize,
			ReadCount:        reads,
			WriteCount:       writes,
			MergedReadCount:  mergedReads,
			MergedWriteCount: mergedWrites,
			ReadTime:         rtime,
			WriteTime:        wtime,
			IopsInProgress:   iopsInProgress,
			IoTime:           iotime,
			WeightedIO:       weightedIO,
		}
	}

	return results, nil
}
