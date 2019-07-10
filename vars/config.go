package vars

import (
	"time"

	"github.com/baudtime/baudtime/util/toml"
)

type InputsConfig struct {
	CollectFrequency toml.Duration `toml:"collect_frequency"`
	Enabled          []string      `toml:"enabled,omitempty"`
}

type ProcessorsConfig struct {
	Enabled []string `toml:"enabled,omitempty"`
}

type AggregatorsConfig struct {
	AggregateFrequency toml.Duration `toml:"aggregate_frequency"`
	Enabled            []string      `toml:"enabled,omitempty"`
}

type OutputsConfig struct {
	BatchFlush int      `toml:"batch_flush"`
	Enabled    []string `toml:"enabled,omitempty"`
}

type Config struct {
	Inputs      InputsConfig      `toml:"inputs"`
	Processors  ProcessorsConfig  `toml:"processors"`
	Aggregators AggregatorsConfig `toml:"aggregators"`
	Outputs     OutputsConfig     `toml:"outputs"`
	ServicePort *int              `toml:"service_port,omitempty"`
}

var Cfg = &Config{
	Inputs: InputsConfig{
		CollectFrequency: toml.Duration(1 * time.Second),
		Enabled:          []string{"cpu", "diskio", "tcp"},
	},
	Processors: ProcessorsConfig{},
	Aggregators: AggregatorsConfig{
		AggregateFrequency: toml.Duration(1 * time.Second),
		Enabled:            []string{"minmax"},
	},
	Outputs: OutputsConfig{
		BatchFlush: 16,
		Enabled:    []string{"logger"},
	},
}

func LoadConfig(tomlFile string) error {
	err := toml.LoadFromToml(tomlFile, Cfg)
	if err != nil {
		return err
	}
	return nil
}
