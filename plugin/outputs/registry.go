package outputs

import (
	"fmt"

	"github.com/baudtime/agent/plugin"
)

var (
	factories    = make(map[string]func() (plugin.Output, error))
	enableStates = make(map[string]bool)
)

func register(name string, isDefaultEnabled bool, factory func() (plugin.Output, error)) {
	enableStates[name] = isDefaultEnabled
	factories[name] = factory
}

func Filter(filters ...string) ([]plugin.Output, []string, error) {
	f := make(map[string]bool)

	for _, filter := range filters {
		if _, exist := enableStates[filter]; !exist {
			return nil, nil, fmt.Errorf("missing output plugin: %s", filter)
		}
		f[filter] = true
	}

	filters = filters[:0]
	var outputs []plugin.Output

	for key, enabled := range enableStates {
		if (len(f) == 0 && enabled) || f[key] {
			if input, err := factories[key](); err != nil {
				return nil, nil, err
			} else {
				outputs = append(outputs, input)
				filters = append(filters, key)
			}
		}
	}
	return outputs, filters, nil
}
