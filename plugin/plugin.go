package plugin

import (
	"strconv"
	"strings"

	"github.com/baudtime/agent/vars"
	"github.com/baudtime/baudtime/msg"
	"github.com/cespare/xxhash/v2"
	"github.com/go-kit/kit/log"
)

const sep = '\xff'

type Metric struct {
	Name   string
	Labels []msg.Label
	Value  float64
}

func (metric Metric) With(labels ...msg.Label) Metric {
	metric.Labels = append(metric.Labels, labels...)
	return metric
}

func (metric Metric) Hash() uint64 {
	var buf []byte
	for _, v := range metric.Labels {
		buf = append(buf, v.Name...)
		buf = append(buf, sep)
		buf = append(buf, v.Value...)
		buf = append(buf, sep)
	}
	return xxhash.Sum64(buf)
}

func (metric Metric) Clone() Metric {
	copy := Metric{
		Name:   metric.Name,
		Labels: make([]msg.Label, 0, len(metric.Labels)),
		Value:  metric.Value,
	}
	copy.Labels = append(copy.Labels, metric.Labels...)
	return copy
}

func (metric Metric) String() string {
	var sb strings.Builder

	sb.WriteString(metric.Name)
	sb.WriteString("{")
	for i, label := range metric.Labels {
		if i != 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(label.Name)
		sb.WriteByte('=')
		sb.WriteString(label.Value)
	}
	sb.WriteString("} ")
	sb.WriteString(strconv.FormatFloat(metric.Value, 'f', -1, 64))

	return sb.String()
}

const (
	DefaultEnabled  = true
	DefaultDisabled = false
)

type Input interface {
	Collect(chan<- Metric) error
}

type Output interface {
	Connect() error
	Close() error
	Write([]Metric) error
}

type Processor interface {
	Apply(Metric) Metric
}

type Aggregator interface {
	Add(Metric) (err error, dropOriginal bool)
	Push(chan<- Metric) error
	Reset()
}

func Logger(plugin string) func(level func(logger log.Logger) log.Logger) log.Logger {
	return func(level func(logger log.Logger) log.Logger) log.Logger {
		return log.With(level(vars.Logger), "plugin", plugin)
	}
}
