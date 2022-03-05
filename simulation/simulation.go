package simulation

import (
	"errors"

	"github.com/HenryGeorgist/go-wat/compute"
)

type Configuration interface {
	ProgramOrder() compute.ProgramOrder
	InputSource() string
	OutputDestination() string
}
type DeterministicCompute struct {
	ProgramOrder      compute.ProgramOrder
	TimeWindow        compute.TimeWindow
	OutputDestination string
	InputSource       string
}
type StochasticConfiguration struct {
	ProgramOrder           compute.ProgramOrder
	LifecycleTimeWindow    compute.TimeWindow
	InitialRealizationSeed int64
	InitialEventSeed       int64
	OutputDestination      string
	InputSource            string
}

func Compute(config Configuration) error {
	stochastic, ok := config.(StochasticConfiguration)
	var coptions compute.ComputeOptions
	if ok {
		//loop for realizations
		//loop for lifecycles
		//loop for events

		seo := compute.StochasticEventOptions{
			RealizationNumber: 1,
			LifecycleNumber:   1,
			EventNumber:       1,
			RealizationSeed:   stochastic.InitialRealizationSeed,
			EventSeed:         stochastic.InitialEventSeed,
			TimeWindow:        stochastic.LifecycleTimeWindow,
		}
		coptions = compute.ComputeOptions{
			InputSource:       config.InputSource(),
			OutputDestination: config.OutputDestination(),
			EventOptions:      seo,
		}
		computeEvent(config, coptions)
	} else {
		//assume deterministic
		deterministic, _ := config.(DeterministicConfiguration)
		deo := compute.DeterministicEventOptions{
			TimeWindow: deterministic.TimeWindow,
		}
		coptions = compute.ComputeOptions{
			InputSource:       config.InputSource(),
			OutputDestination: config.OutputDestination(),
			EventOptions:      deo,
		}
		return computeEvent(config, coptions)
	}
	return errors.New("something bad happened")
}
func computeEvent(config Configuration, options compute.ComputeOptions) error {
	for _, p := range config.ProgramOrder().Plugins {
		err := p.Compute(options)
		if err != nil {
			return err
		}
	}
	return nil
}
