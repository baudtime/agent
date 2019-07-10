package manager

import (
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/baudtime/agent/plugin"
	"github.com/baudtime/agent/plugin/aggregators"
	"github.com/baudtime/agent/plugin/inputs"
	"github.com/baudtime/agent/plugin/outputs"
	"github.com/baudtime/agent/plugin/processors"
	"github.com/baudtime/agent/vars"
	"github.com/go-kit/kit/log/level"
)

type runningAggregator struct {
	Aggregator
	sync.Mutex
}

func (agg *runningAggregator) Add(metric Metric) (err error, dropOriginal bool) {
	agg.Lock()
	err, dropOriginal = agg.Aggregator.Add(metric)
	agg.Unlock()
	return
}

func (agg *runningAggregator) Push(ch chan<- Metric) (err error) {
	agg.Lock()
	err = agg.Aggregator.Push(ch)
	agg.Unlock()
	return
}

func (agg *runningAggregator) Reset() {
	agg.Lock()
	agg.Aggregator.Reset()
	agg.Unlock()
}

type runningOutput struct {
	Output
	buf []Metric
}

func (o *runningOutput) Add(metric Metric) {
	o.buf = append(o.buf, metric.Clone())
	if len(o.buf) > vars.Cfg.Outputs.BatchFlush {
		o.Write(o.buf)
		o.buf = o.buf[:0]
	}
}

type Manager struct {
	inputs      []Input
	processors  []Processor
	aggregators []*runningAggregator
	outputs     []*runningOutput
	wg          sync.WaitGroup
	exitC       chan struct{}
}

func New(inputFilters []string, processorsFilter []string, aggregatorsFilter []string, outputsFilter []string) (*Manager, error) {
	inputs, inputFilters, err := inputs.Filter(inputFilters...)
	if err != nil {
		return nil, err
	}
	level.Info(vars.Logger).Log("###", fmt.Sprintf("### INPUTS ENABLED ###: [%s]", strings.Join(inputFilters, ",")))

	processors, processorsFilter, err := processors.Filter(processorsFilter...)
	if err != nil {
		return nil, err
	}
	level.Info(vars.Logger).Log("###", fmt.Sprintf("### PROCESSORS ENABLED ###: [%s]", strings.Join(processorsFilter, ",")))

	aggregators, aggregatorsFilter, err := aggregators.Filter(aggregatorsFilter...)
	if err != nil {
		return nil, err
	}
	level.Info(vars.Logger).Log("###", fmt.Sprintf("### AGGREGATORS ENABLED ###: [%s]", strings.Join(aggregatorsFilter, ",")))

	outputs, outputsFilter, err := outputs.Filter(outputsFilter...)
	if err != nil {
		return nil, err
	}
	level.Info(vars.Logger).Log("###", fmt.Sprintf("### OUTPUTS ENABLED ###: [%s]", strings.Join(outputsFilter, ",")))

	manager := &Manager{
		inputs:     inputs,
		processors: processors,
		exitC:      make(chan struct{}),
	}
	for _, aggregator := range aggregators {
		manager.aggregators = append(manager.aggregators, &runningAggregator{Aggregator: aggregator})
	}
	for _, output := range outputs {
		manager.outputs = append(manager.outputs, &runningOutput{Output: output})
	}
	return manager, nil
}

func (manager *Manager) Start() {
	var src, dest chan Metric

	dest = make(chan Metric, 512)
	manager.runInputs(dest)
	src = dest

	if len(manager.processors) > 0 {
		dest = make(chan Metric, 512)
		manager.runProcessors(src, dest)
		src = dest
	}

	if len(manager.aggregators) > 0 {
		dest = make(chan Metric, 512)
		manager.runAggregators(src, dest)
		src = dest
	}

	manager.runOutputs(src)
}

func (manager *Manager) Stop() {
	close(manager.exitC)
	manager.wg.Wait()
	level.Info(vars.Logger).Log("msg", "plugins manager exits")
}

func (manager *Manager) runInputs(dest chan Metric) {
	manager.wg.Add(1)
	go func() {
		defer func() {
			close(dest)
			level.Info(vars.Logger).Log("msg", "inputs do not provide data any more")
			manager.wg.Done()
		}()

		ticker := time.NewTicker(time.Duration(vars.Cfg.Inputs.CollectFrequency))
		defer ticker.Stop()

		var wg sync.WaitGroup

		for {
			select {
			case <-ticker.C:
				for _, input := range manager.inputs {
					wg.Add(1)
					go func(input Input) {
						defer wg.Done()

						if err := input.Collect(dest); err != nil {
							level.Error(vars.Logger).Log("err", err)
						}
					}(input)
				}
				wg.Wait()
			case <-manager.exitC:
				return
			}
		}
	}()
}

func (manager *Manager) runProcessors(src, dest chan Metric) {
	manager.wg.Add(1)
	go func() {
		defer manager.wg.Done()

		for metric := range src {
			for _, processor := range manager.processors {
				metric = processor.Apply(metric)
			}
			dest <- metric
		}

		level.Info(vars.Logger).Log("msg", "processors do not consume data any more")
	}()
}

func (manager *Manager) runAggregators(src, dest chan Metric) {
	manager.wg.Add(1)
	go func() {
		defer manager.wg.Done()

		for metric := range src {
			for _, agg := range manager.aggregators {
				if _, dropOriginal := agg.Add(metric.Clone()); !dropOriginal {
					dest <- metric
				}
			}
		}

		level.Info(vars.Logger).Log("msg", "aggregators do not consume data any more")
	}()

	manager.wg.Add(1)
	go func() {
		defer func() {
			close(dest)
			level.Info(vars.Logger).Log("msg", "aggregators do not privide data any more")
			manager.wg.Done()
		}()

		ticker := time.NewTicker(time.Duration(vars.Cfg.Aggregators.AggregateFrequency))
		defer ticker.Stop()

		var wg sync.WaitGroup

		for {
			select {
			case <-ticker.C:
				for _, agg := range manager.aggregators {
					wg.Add(1)
					go func(agg *runningAggregator) {
						defer wg.Done()

						if err := agg.Push(dest); err != nil {
							level.Error(vars.Logger).Log("err", err)
						} else {
							agg.Reset()
						}
					}(agg)
				}
				wg.Wait()
			case <-manager.exitC:
				return
			}
		}
	}()
}

func (manager *Manager) runOutputs(src chan Metric) {
	manager.wg.Add(1)
	go func() {
		defer manager.wg.Done()

		for _, output := range manager.outputs {
			output.Connect()
		}

		var wg sync.WaitGroup

		for metric := range src {
			for _, output := range manager.outputs {
				wg.Add(1)
				go func(output *runningOutput) {
					defer wg.Done()

					output.Add(metric)
				}(output)
			}

			wg.Wait()
		}

		for _, output := range manager.outputs {
			output.Close()
		}

		level.Info(vars.Logger).Log("msg", "outputs do not consume data any more")
	}()
}
