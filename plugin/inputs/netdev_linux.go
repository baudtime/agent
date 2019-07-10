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
	rxBytesMetric   = "net_dev_rx_bytes"
	rxPacketsMetric = "net_dev_rx_packets"
	rxErrsMetric    = "net_dev_rx_errs"
	rxDropMetric    = "net_dev_rx_drop"
	txBytesMetric   = "net_dev_tx_bytes"
	txPacketsMetric = "net_dev_tx_packets"
	txErrsMetric    = "net_dev_tx_errs"
	txDropMetric    = "net_dev_tx_drop"
)

var (
	netDevIgnoredDevices = kingpin.Flag("input.netdev.ignored-devices", "Regexp of net devices to ignore for netdev collector.").
				Default("^(lo|virtual|(docker|br|cane|ovs|ifb).+)$").String()
	netDevPlgLogger = plugin.Logger("netdev")
)

func init() {
	register("netdev", plugin.DefaultDisabled, newNetDevCollector)
}

type netDevCollector struct {
	ignoredDevicesPattern *regexp.Regexp
	lastStats             map[string]NetDevStats
	lastCollectTime       time.Time
}

func newNetDevCollector() (plugin.Input, error) {
	return &netDevCollector{
		ignoredDevicesPattern: regexp.MustCompile(*netDevIgnoredDevices),
		lastStats:             make(map[string]NetDevStats),
	}, nil
}

func (c *netDevCollector) Collect(ch chan<- plugin.Metric) error {
	netDevStats, err := readNetDevStat(procFilePath("net/dev"), c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}

	now := time.Now()
	elapsedSec := now.Sub(c.lastCollectTime).Seconds()

	for i, devStat := range netDevStats {
		if lastStat, found := c.lastStats[devStat.Iface]; found {
			devLabel := msg.Label{"iface", netDevStats[i].Iface}

			rxBytes := uint64(float64(devStat.RxBytes-lastStat.RxBytes) / elapsedSec)
			rxPackets := uint64(float64(devStat.RxPackets-lastStat.RxPackets) / elapsedSec)
			rxErrs := uint64(float64(devStat.RxErrs-lastStat.RxErrs) / elapsedSec)
			rxDrop := uint64(float64(devStat.RxDrop-lastStat.RxDrop) / elapsedSec)
			txBytes := uint64(float64(devStat.TxBytes-lastStat.TxBytes) / elapsedSec)
			txPackets := uint64(float64(devStat.TxPackets-lastStat.TxPackets) / elapsedSec)
			txErrs := uint64(float64(devStat.TxErrs-lastStat.TxErrs) / elapsedSec)
			txDrop := uint64(float64(devStat.TxDrop-lastStat.TxDrop) / elapsedSec)

			ch <- plugin.Metric{Name: rxBytesMetric, Value: float64(rxBytes)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: rxPacketsMetric, Value: float64(rxPackets)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: rxErrsMetric, Value: float64(rxErrs)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: rxDropMetric, Value: float64(rxDrop)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: txBytesMetric, Value: float64(txBytes)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: txPacketsMetric, Value: float64(txPackets)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: txErrsMetric, Value: float64(txErrs)}.With(HostLabels...).With(devLabel)
			ch <- plugin.Metric{Name: txDropMetric, Value: float64(txDrop)}.With(HostLabels...).With(devLabel)
		}
	}

	c.lastStats = netDevStats
	c.lastCollectTime = now

	return nil
}

type NetDevStats struct {
	Iface        string `json:"iface"`
	RxBytes      uint64 `json:"rxbytes"`
	RxPackets    uint64 `json:"rxpackets"`
	RxErrs       uint64 `json:"rxerrs"`
	RxDrop       uint64 `json:"rxdrop"`
	RxFifo       uint64 `json:"rxfifo"`
	RxFrame      uint64 `json:"rxframe"`
	RxCompressed uint64 `json:"rxcompressed"`
	RxMulticast  uint64 `json:"rxmulticast"`
	TxBytes      uint64 `json:"txbytes"`
	TxPackets    uint64 `json:"txpackets"`
	TxErrs       uint64 `json:"txerrs"`
	TxDrop       uint64 `json:"txdrop"`
	TxFifo       uint64 `json:"txfifo"`
	TxColls      uint64 `json:"txcolls"`
	TxCarrier    uint64 `json:"txcarrier"`
	TxCompressed uint64 `json:"txcompressed"`
}

func readNetDevStat(path string, ignore *regexp.Regexp) (map[string]NetDevStats, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// lines[2:] remove /proc/net/dev header
	results := make(map[string]NetDevStats)

	for _, line := range lines[2:] {
		// patterns
		// <iface>: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
		// or
		// <iface>:0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 (without space after colon)
		colon := strings.Index(line, ":")

		if colon > 0 {
			metrics := line[colon+1:]
			fields := strings.Fields(metrics)

			iface := strings.Replace(line[0:colon], " ", "", -1)
			if ignore.MatchString(iface) {
				netDevPlgLogger(level.Debug).Log("msg", "ignoring device", "dev", iface)
				continue
			}

			var netDevStat NetDevStats

			netDevStat.Iface = iface
			netDevStat.RxBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			netDevStat.RxPackets, _ = strconv.ParseUint(fields[1], 10, 64)
			netDevStat.RxErrs, _ = strconv.ParseUint(fields[2], 10, 64)
			netDevStat.RxDrop, _ = strconv.ParseUint(fields[3], 10, 64)
			netDevStat.RxFifo, _ = strconv.ParseUint(fields[4], 10, 64)
			netDevStat.RxFrame, _ = strconv.ParseUint(fields[5], 10, 64)
			netDevStat.RxCompressed, _ = strconv.ParseUint(fields[6], 10, 64)
			netDevStat.RxMulticast, _ = strconv.ParseUint(fields[7], 10, 64)
			netDevStat.TxBytes, _ = strconv.ParseUint(fields[8], 10, 64)
			netDevStat.TxPackets, _ = strconv.ParseUint(fields[9], 10, 64)
			netDevStat.TxErrs, _ = strconv.ParseUint(fields[10], 10, 64)
			netDevStat.TxDrop, _ = strconv.ParseUint(fields[11], 10, 64)
			netDevStat.TxFifo, _ = strconv.ParseUint(fields[12], 10, 64)
			netDevStat.TxColls, _ = strconv.ParseUint(fields[13], 10, 64)
			netDevStat.TxCarrier, _ = strconv.ParseUint(fields[14], 10, 64)
			netDevStat.TxCompressed, _ = strconv.ParseUint(fields[15], 10, 64)

			results[iface] = netDevStat
		}
	}

	return results, nil
}
