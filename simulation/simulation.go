package simulation

import (
	"errors"
	"math/rand"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

//Configuration defines the interface for a simulation compute configuration
type Configuration interface {
	ProgramOrder() component.ProgramOrder
	Models() []component.Model
	InputSource() string
	OutputDestination() string
}

//DeterministicConfiguration implement the Configuration interface for a DeterministicCompute
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

//StochasticConfiguration implements the Configuration interface for a Stochastic Simulation
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
		eventRandom := rand.NewSource(stochastic.InitialEventSeed)
		realizationRandom := rand.NewSource(stochastic.InitialRealizationSeed)
		//each realization can be run conccurrently
		for realization := 0; realization < stochastic.TotalRealizations; realization++ {
			//loop for lifecycles
			realizationSeed := realizationRandom.Int63() //probably make one per model
			//each lifecycle can be a job run concurrently
			for lifecycle := 0; lifecycle < stochastic.LifecyclesPerRealization; lifecycle++ {
				//events in a lifecycle are dependent on earlier events,
				//events should not be run concurrently
				//event generator create events
				//loop for events
				events := stochastic.EventGenerator.GenerateTimeWindows(stochastic.LifecycleTimeWindow)
				for eventid, event := range events {
					eventSeed := eventRandom.Int63() //probably make one per model
					seo := compute.StochasticEventOptions{
						RealizationNumber: realization,
						LifecycleNumber:   lifecycle,
						EventNumber:       eventid,
						RealizationSeed:   realizationSeed,
						EventSeed:         eventSeed,
					}
					seo.UpdateTimeWindow(event)
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
		deo := compute.DeterministicEventOptions{}
		deo.UpdateTimeWindow(deterministic.TimeWindow)
		coptions = compute.ComputeOptions{
			InputSource:       config.InputSource(),
			OutputDestination: config.OutputDestination(),
			EventOptions:      deo,
		}
		return computeEvent(config, coptions)
	}
	return errors.New("something bad happened")
}

//computeEvent iterates over the program order and requests each plugin to compute the associated model in the model list.
func computeEvent(config Configuration, options compute.ComputeOptions) error {
	for idx, p := range config.ProgramOrder().Plugins {
		err := p.Compute(config.Models()[idx], options)
		if err != nil {
			return err
		}
		//if the plugin implements the outputrecorder interface, get outputs
	}
	return nil
}
