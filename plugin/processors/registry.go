package processors

import (
	"fmt"

	"github.com/baudtime/agent/plugin"
)

var (
	factories    = make(map[string]func() (plugin.Processor, error))
	enableStates = make(map[string]bool)
)

func register(name string, isDefaultEnabled bool, factory func() (plugin.Processor, error)) {
	enableStates[name] = isDefaultEnabled
	factories[name] = factory
}

func Filter(filters ...string) ([]plugin.Processor, []string, error) {
	f := make(map[string]bool)

	for _, filter := range filters {
		if _, exist := enableStates[filter]; !exist {
			return nil, nil, fmt.Errorf("missing processor plugin: %s", filter)
		}
		f[filter] = true
	}

	filters = filters[:0]
	var processors []plugin.Processor

	for key, enabled := range enableStates {
		if (len(f) == 0 && enabled) || f[key] {
			if input, err := factories[key](); err != nil {
				return nil, nil, err
			} else {
				processors = append(processors, input)
				filters = append(filters, key)
			}
		}
	}
	return processors, filters, nil
}
