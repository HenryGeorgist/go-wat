package simulation

import (
	"errors"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type Configuration interface {
	ProgramOrder() component.ProgramOrder
	Models() []component.Model
	InputSource() string
	OutputDestination() string
}
type DeterministicConfiguration struct {
	programOrder      component.ProgramOrder
	models            []component.Model
	TimeWindow        compute.TimeWindow
	outputDestination string
	inputSource       string
}

func (d DeterministicConfiguration) ProgramOrder() component.ProgramOrder {
	return d.programOrder
}
func (d DeterministicConfiguration) Models() []component.Model {
	return d.models
}
func (d DeterministicConfiguration) InputSource() string {
	return d.inputSource
}
func (d DeterministicConfiguration) OutputDestination() string {
	return d.outputDestination
}

type StochasticConfiguration struct {
	programOrder             component.ProgramOrder
	models                   []component.Model
	EventGenerator           component.EventGenerator
	LifecycleTimeWindow      compute.TimeWindow
	TotalRealizations        int
	LifecyclesPerRealization int
	InitialRealizationSeed   int64
	InitialEventSeed         int64
	outputDestination        string
	inputSource              string
}

func (s StochasticConfiguration) ProgramOrder() component.ProgramOrder {
	return s.programOrder
}
func (s StochasticConfiguration) Models() []component.Model {
	return s.models
}
func (s StochasticConfiguration) InputSource() string {
	return s.inputSource
}
func (s StochasticConfiguration) OutputDestination() string {
	return s.outputDestination
}

func Compute(config Configuration) error {
	stochastic, ok := config.(StochasticConfiguration)
	var coptions compute.ComputeOptions
	if ok {
		//loop for realizations
		for realization := 0; realization < stochastic.TotalRealizations; realization++ {
			//loop for lifecycles
			for lifecycle := 0; lifecycle < stochastic.LifecyclesPerRealization; lifecycle++ {
				//loop for events
				//event generator create events
				events := stochastic.EventGenerator.GenerateTimeWindows()
				for eventid, event := range events {
					seo := compute.StochasticEventOptions{
						RealizationNumber: realization,
						LifecycleNumber:   lifecycle,
						EventNumber:       eventid,
						RealizationSeed:   stochastic.InitialRealizationSeed,
						EventSeed:         stochastic.InitialEventSeed,
						TimeWindow:        event,
					}
					coptions = compute.ComputeOptions{
						InputSource:       config.InputSource(),
						OutputDestination: config.OutputDestination(),
						EventOptions:      seo,
					}
					computeEvent(config, coptions)
				}

			}

		}

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
	for idx, p := range config.ProgramOrder().Plugins {
		err := p.Compute(config.Models()[idx], options)
		if err != nil {
			return err
		}
	}
	return nil
}
