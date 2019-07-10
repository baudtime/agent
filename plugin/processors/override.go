package processors

import (
	"github.com/baudtime/agent/plugin"
	"github.com/baudtime/baudtime/msg"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	override = kingpin.Flag("processor.override.override", "override the whole name").String()
	prefix   = kingpin.Flag("processor.override.prefix", "prefix of the original name").String()
	suffix   = kingpin.Flag("processor.override.suffix", "suffix of the original name").String()
)

func init() {
	register("override", plugin.DefaultDisabled, func() (plugin.Processor, error) {
		return &Override{
			NameOverride: *override,
			NamePrefix:   *prefix,
			NameSuffix:   *suffix,
		}, nil
	})
}

type Override struct {
	NameOverride string
	NamePrefix   string
	NameSuffix   string
	Labels       map[string]string
}

func (p *Override) Apply(metric plugin.Metric) plugin.Metric {
	if len(p.NameOverride) > 0 {
		metric.Name = p.NameOverride
	}
	if len(p.NamePrefix) > 0 {
		metric.Name = p.NamePrefix + metric.Name
	}
	if len(p.NameSuffix) > 0 {
		metric.Name = metric.Name + p.NameSuffix
	}
	for key, value := range p.Labels {
		metric.With(msg.Label{key, value})
	}
	return metric
}
