package inputs

import (
	"fmt"
	"net"
	"time"

	"github.com/baudtime/agent/plugin"
	. "github.com/baudtime/agent/vars"
	"github.com/beevik/ntp"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	ntpClockOffsetMilliMetric = "ntp_clk_offset_ms"
)

var (
	ntpServer          = kingpin.Flag("input.ntp.server", "NTP server to use for ntp collector").Default("127.0.0.1").String()
	ntpProtocolVersion = kingpin.Flag("input.ntp.protocol-version", "NTP protocol version").Default("4").Int()
	ntpServerIsLocal   = kingpin.Flag("input.ntp.server-is-local", "Certify that collector.ntp.server address is the same local host as this collector.").Default("false").Bool()
)

type ntpCollector struct{}

func init() {
	register("ntp", plugin.DefaultDisabled, newNtpCollector)
}

// newNtpCollector returns a new Collector exposing sanity of local NTP server.
// Default definition of "local" is:
// - collector.ntp.server address is a loopback address (or collector.ntp.server-is-mine flag is turned on)
// - the server is reachable with outgoin IP_TTL = 1
func newNtpCollector() (plugin.Input, error) {
	ipaddr := net.ParseIP(*ntpServer)
	if !*ntpServerIsLocal && (ipaddr == nil || !ipaddr.IsLoopback()) {
		return nil, fmt.Errorf("only IP address of local NTP server is valid for --collector.ntp.server")
	}

	if *ntpProtocolVersion < 2 || *ntpProtocolVersion > 4 {
		return nil, fmt.Errorf("invalid NTP protocol version %d; must be 2, 3, or 4", *ntpProtocolVersion)
	}

	return &ntpCollector{}, nil
}

func (c *ntpCollector) Collect(ch chan<- plugin.Metric) error {
	clkOffset, err := getNtpClockOffset()
	if err != nil {
		return err
	}

	ch <- plugin.Metric{Name: ntpClockOffsetMilliMetric, Value: float64(clkOffset)}.With(HostLabels...)
	return nil
}

func getNtpClockOffset() (int64, error) {
	resp, err := ntp.QueryWithOptions(*ntpServer, ntp.QueryOptions{
		Version: *ntpProtocolVersion,
		Timeout: 2 * time.Second, // default `ntpdate` timeout
	})

	if err != nil {
		return -1, fmt.Errorf("couldn't get SNTP reply: %s", err)
	}

	return resp.ClockOffset.Milliseconds(), nil
}
