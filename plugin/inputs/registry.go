package inputs

import (
	"fmt"

	"github.com/baudtime/agent/plugin"
)

var (
	factories    = make(map[string]func() (plugin.Input, error))
	enableStates = make(map[string]bool)
)

func register(name string, isDefaultEnabled bool, factory func() (plugin.Input, error)) {
	enableStates[name] = isDefaultEnabled
	factories[name] = factory
}

func Filter(filters ...string) ([]plugin.Input, []string, error) {
	f := make(map[string]bool)

	for _, filter := range filters {
		if _, exist := enableStates[filter]; !exist {
			return nil, nil, fmt.Errorf("missing input plugin: %s", filter)
		}
		f[filter] = true
	}

	filters = filters[:0]
	var inputs []plugin.Input

	for key, enabled := range enableStates {
		if (len(f) == 0 && enabled) || f[key] {
			if input, err := factories[key](); err != nil {
				return nil, nil, err
			} else {
				inputs = append(inputs, input)
				filters = append(filters, key)
			}
		}
	}
	return inputs, filters, nil
}
