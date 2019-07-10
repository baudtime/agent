package aggregators

import "github.com/baudtime/agent/plugin"

func init() {
	register("minmax", plugin.DefaultDisabled, func() (plugin.Aggregator, error) {
		return &MinMax{cache: make(map[uint64]aggregate)}, nil
	})
}

type MinMax struct {
	cache map[uint64]aggregate
}

type aggregate struct {
	metric plugin.Metric
	min    float64
	max    float64
}

func (m *MinMax) Add(in plugin.Metric) (error, bool) {
	id := in.Hash()
	if _, ok := m.cache[id]; !ok {
		m.cache[id] = aggregate{
			metric: in,
			min:    in.Value,
			max:    in.Value,
		}
	} else {
		if in.Value < m.cache[id].min {
			tmp := m.cache[id]
			tmp.min = in.Value
			m.cache[id] = tmp
		} else if in.Value > m.cache[id].max {
			tmp := m.cache[id]
			tmp.max = in.Value
			m.cache[id] = tmp
		}
	}
	return nil, false
}

func (m *MinMax) Push(ch chan<- plugin.Metric) error {
	for _, aggregate := range m.cache {
		ch <- plugin.Metric{aggregate.metric.Name + "_min", aggregate.metric.Labels, aggregate.min}
		ch <- plugin.Metric{aggregate.metric.Name + "_max", aggregate.metric.Labels, aggregate.max}
	}

	return nil
}

func (m *MinMax) Reset() {
	for k, _ := range m.cache {
		delete(m.cache, k)
	}
}
