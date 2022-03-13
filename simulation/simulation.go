package simulation

import (
	"errors"
	"fmt"
	"os"

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
	Programorder      component.ProgramOrder `json:"programorder"`
	ModelList         []component.Model      `json:"models"`
	TimeWindow        compute.TimeWindow     `json:"timewindow"`
	Outputdestination string                 `json:"outputdestination"`
	Inputsource       string                 `json:"inputsource"`
}

func (d DeterministicConfiguration) ProgramOrder() component.ProgramOrder {
	return d.Programorder
}
func (d DeterministicConfiguration) Models() []component.Model {
	return d.ModelList
}
func (d DeterministicConfiguration) InputSource() string {
	return d.Inputsource
}
func (d DeterministicConfiguration) OutputDestination() string {
	return d.Outputdestination
}

//StochasticConfiguration implements the Configuration interface for a Stochastic Simulation
type StochasticConfiguration struct {
	Programorder             component.ProgramOrder   `json:"programorder"`
	ModelList                []component.Model        `json:"models"`
	EventGenerator           component.EventGenerator `json:"eventgenerator"`
	LifecycleTimeWindow      compute.TimeWindow       `json:"timewindow"`
	TotalRealizations        int                      `json:"totalrealizations"`
	LifecyclesPerRealization int                      `json:"lifecyclesperrealization"`
	InitialRealizationSeed   int64                    `json:"initialrealizationseed"`
	InitialEventSeed         int64                    `json:"intitaleventseed"`
	Outputdestination        string                   `json:"outputdestination"`
	Inputsource              string                   `json:"inputsource"`
}

func (s StochasticConfiguration) ProgramOrder() component.ProgramOrder {
	return s.Programorder
}
func (s StochasticConfiguration) Models() []component.Model {
	return s.ModelList
}
func (s StochasticConfiguration) InputSource() string {
	return s.Inputsource
}
func (s StochasticConfiguration) OutputDestination() string {
	return s.Outputdestination
}

func Compute(config Configuration) error {
	stochastic, ok := config.(StochasticConfiguration)
	var coptions compute.Options
	if ok {
		//loop for realizations
		rootOutputPath := config.OutputDestination()
		rootinputPath := config.InputSource()
		eventRandom := component.SeedManager{
			Seed:        stochastic.InitialEventSeed,
			PluginCount: len(config.Models()),
		}
		eventRandom.Init()
		realizationRandom := component.SeedManager{
			Seed:        stochastic.InitialRealizationSeed,
			PluginCount: len(config.Models()),
		}
		realizationRandom.Init()
		//each realization can be run conccurrently
		for realization := 0; realization < stochastic.TotalRealizations; realization++ {
			realizationInputPath := rootinputPath + "realization " + fmt.Sprint(realization) + "/"
			realizationOutputPath := rootOutputPath + "realization " + fmt.Sprint(realization) + "/"
			//loop for lifecycles
			realizationSeeds := realizationRandom.GeneratePluginSeeds() //probably make one per model
			//each lifecycle can be a job run concurrently
			for lifecycle := 0; lifecycle < stochastic.LifecyclesPerRealization; lifecycle++ {
				//events in a lifecycle are dependent on earlier events,
				//events should not be run concurrently
				//event generator create events
				lifecycleInputPath := realizationInputPath + "lifecycle " + fmt.Sprint(lifecycle) + "/"
				lifecycleOutputPath := realizationOutputPath + "lifecycle " + fmt.Sprint(lifecycle) + "/"
				//loop for events
				events := stochastic.EventGenerator.GenerateTimeWindows(stochastic.LifecycleTimeWindow)
				for eventid, event := range events {
					stochastic.Inputsource = lifecycleInputPath + "event " + fmt.Sprint(eventid)
					stochastic.Outputdestination = lifecycleOutputPath + "event " + fmt.Sprint(eventid)
					_ = os.MkdirAll(stochastic.InputSource(), os.ModeTemporary)
					_ = os.MkdirAll(stochastic.OutputDestination(), os.ModeTemporary)
					eventSeeds := eventRandom.GeneratePluginSeeds() //probably make one per model
					seo := compute.StochasticEventOptions{
						RealizationNumber: realization,
						LifecycleNumber:   lifecycle,
						EventNumber:       eventid,
						RealizationSeeds:  realizationSeeds,
						EventSeeds:        eventSeeds,
						EventTimeWindow:   event,
						CurrentModel:      0,
					}
					coptions = compute.Options{
						InputSource:       stochastic.InputSource(),
						OutputDestination: stochastic.OutputDestination(),
						EventOptions:      seo,
					}
					computeEvent(config, coptions)
				}

			}

		}

	} else {
		//assume deterministic
		deterministic, _ := config.(DeterministicConfiguration)
		deo := compute.DeterministicEventOptions{EventTimeWindow: deterministic.TimeWindow}
		coptions = compute.Options{
			InputSource:       config.InputSource(),
			OutputDestination: config.OutputDestination(),
			EventOptions:      deo,
		}
		return computeEvent(config, coptions)
	}
	return errors.New("something bad happened")
}

//computeEvent iterates over the program order and requests each plugin to compute the associated model in the model list.
func computeEvent(config Configuration, options compute.Options) error {
	for idx, p := range config.ProgramOrder().Plugins {
		err := p.Compute(config.Models()[idx], options)
		if err != nil {
			return err
		}
		options.EventOptions = options.IncrementModelIndex()
		//if the plugin implements the outputrecorder interface, get outputs
	}
	return nil
}
