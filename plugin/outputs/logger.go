package outputs

import (
	"log"

	"github.com/baudtime/agent/plugin"
)

func init() {
	register("logger", plugin.DefaultEnabled, newLogger)
}

type logger struct{}

func newLogger() (plugin.Output, error) {
	return &logger{}, nil
}

func (l *logger) Connect() error {
	return nil
}

func (l *logger) Write(metrics []plugin.Metric) error {
	for _, metric := range metrics {
		log.Println(metric)
	}
	return nil
}

func (l *logger) Close() error {
	return nil
}
