package inputs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/baudtime/agent/plugin"
	. "github.com/baudtime/agent/vars"
)

const (
	loadAvgMetric        = "load_avg"
	processRunningMetric = "proc_running_num"
)

func init() {
	register("loadavg", plugin.DefaultDisabled, newLoadavgCollector)
}

type loadavgCollector struct{}

func newLoadavgCollector() (plugin.Input, error) {
	return &loadavgCollector{}, nil
}

func (c *loadavgCollector) Collect(ch chan<- plugin.Metric) error {
	loadAvg, err := readLoadAvg(procFilePath("loadavg"))
	if err != nil {
		return fmt.Errorf("couldn't get load: %s", err)
	}

	ch <- plugin.Metric{Name: loadAvgMetric, Value: loadAvg.Last1Min}.With(HostLabels...)
	ch <- plugin.Metric{Name: processRunningMetric, Value: float64(loadAvg.ProcessRunning)}.With(HostLabels...)
	return nil
}

type LoadAvg struct {
	Last1Min       float64 `json:"last1min"`
	Last5Min       float64 `json:"last5min"`
	Last15Min      float64 `json:"last15min"`
	ProcessRunning uint64  `json:"process_running"`
	ProcessTotal   uint64  `json:"process_total"`
	LastPID        uint64  `json:"last_pid"`
}

func readLoadAvg(path string) (*LoadAvg, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	content := strings.TrimSpace(string(b))
	fields := strings.Fields(content)

	if len(fields) < 5 {
		return nil, errors.New("Cannot parse loadavg: " + content)
	}

	process := strings.Split(fields[3], "/")

	if len(process) != 2 {
		return nil, errors.New("Cannot parse loadavg: " + content)
	}

	loadavg := LoadAvg{}

	if loadavg.Last1Min, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return nil, err
	}

	if loadavg.Last5Min, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return nil, err
	}

	if loadavg.Last15Min, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return nil, err
	}

	if loadavg.ProcessRunning, err = strconv.ParseUint(process[0], 10, 64); err != nil {
		return nil, err
	}

	if loadavg.ProcessTotal, err = strconv.ParseUint(process[1], 10, 64); err != nil {
		return nil, err
	}

	if loadavg.LastPID, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return nil, err
	}

	return &loadavg, nil
}
