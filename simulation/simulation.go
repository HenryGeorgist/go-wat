package simulation

import (
	"fmt"
	"os"

	"go-wat/component"
	"go-wat/compute"

	"github.com/HydrologicEngineeringCenter/go-statistics/data"
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
	Programorder                 component.ProgramOrder   `json:"programorder"`
	ModelList                    []component.Model        `json:"models"`
	EventGenerator               component.EventGenerator `json:"eventgenerator"`
	LifecycleTimeWindow          compute.TimeWindow       `json:"timewindow"`
	TotalRealizations            int                      `json:"totalrealizations"`
	LifecyclesPerRealization     int                      `json:"lifecyclesperrealization"`
	InitialRealizationSeed       int64                    `json:"initialrealizationseed"`
	InitialEventSeed             int64                    `json:"intitaleventseed"`
	Outputdestination            string                   `json:"outputdestination"`
	Inputsource                  string                   `json:"inputsource"`
	DeleteOutputAfterRealization bool                     `json:"delete_after_realization"`
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

	var coptions compute.Options

	stochastic, ok := config.(StochasticConfiguration)
	if ok {
		//develop map of map of inline histograms
		outputvariableHistograms := make(map[string]map[string]*data.InlineHistogram)

		//get output variables by plugin
		outputvariablesMap := make(map[string][]string)
		for idx, model := range stochastic.ModelList {
			ovp, ovpok := stochastic.Programorder.Plugins[idx].(component.OutputReporter)
			if ovpok {
				histogramMap := make(map[string]*data.InlineHistogram)
				variables := ovp.OutputVariables(model)
				for _, v := range variables {
					histogramMap[v] = data.Init(.05, 0, .05)
				}
				outputvariablesMap[model.ModelName()] = variables
				outputvariableHistograms[model.ModelName()] = histogramMap
			}
		}

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
			realizationInputPath := fmt.Sprintf("%srealization-%v", rootinputPath, realization)
			realizationOutputPath := fmt.Sprintf("%srealization-%v", rootOutputPath, realization)

			//loop for lifecycles
			realizationSeeds := realizationRandom.GeneratePluginSeeds() //probably make one per model
			message := fmt.Sprintf("Computing realization %v", realization)
			fmt.Println(message)

			//each lifecycle can be a job run concurrently
			for lifecycle := 0; lifecycle < stochastic.LifecyclesPerRealization; lifecycle++ {
				//events in a lifecycle are dependent on earlier events,
				//events should not be run concurrently
				//event generator create events

				lifecycleInputPath := fmt.Sprintf("%s/lifecycle-%v", realizationInputPath, lifecycle)
				lifecycleOutputPath := fmt.Sprintf("%s/lifecycle-%v", realizationOutputPath, lifecycle)

				//loop for events
				err := stochastic.LifecycleTimeWindow.IsValid()
				if err != nil {
					return err
				}

				events := stochastic.EventGenerator.GenerateTimeWindows(stochastic.LifecycleTimeWindow)

				for eventid, event := range events {

					stochastic.Inputsource = fmt.Sprintf("%s/event-%v", lifecycleInputPath, eventid)
					stochastic.Outputdestination = fmt.Sprintf("%s/event-%v", lifecycleOutputPath, eventid)

					err := os.MkdirAll(stochastic.InputSource(), os.ModePerm)
					if err != nil {
						return err
					}

					err = os.MkdirAll(stochastic.OutputDestination(), os.ModePerm)
					if err != nil {
						return err
					}

					//probably make one per model
					eventSeeds := eventRandom.GeneratePluginSeeds()
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
						OutputVariables:   outputvariablesMap,
					}
					outputvars, err := computeEvent(config, coptions)
					if err != nil {
						return err
					}

					for k, v := range outputvars {
						histos, ok := outputvariableHistograms[k]
						if ok {
							//load em up.
							for hk, hv := range histos {
								hv.AddObservation(v[0])
								histos[hk] = hv
							}
						}
						outputvariableHistograms[k] = histos
					}
				} //events

			} //lifecycles
			if stochastic.DeleteOutputAfterRealization {
				fmt.Println("Deleting " + realizationOutputPath)
				os.RemoveAll(realizationOutputPath)
			}

		} //realizations
		fmt.Println(outputvariableHistograms)
		return nil

	} else {
		//assume deterministic
		outputvariablesMap := make(map[string][]string)
		deterministic, _ := config.(DeterministicConfiguration)

		deo := compute.DeterministicEventOptions{EventTimeWindow: deterministic.TimeWindow}

		coptions = compute.Options{
			InputSource:       config.InputSource(),
			OutputDestination: config.OutputDestination(),
			EventOptions:      deo,
			OutputVariables:   outputvariablesMap,
		}

		err := os.MkdirAll(config.OutputDestination(), os.ModePerm)
		if err != nil {
			return err
		}
		// syscall.Umask(0)

		err = os.MkdirAll(config.OutputDestination(), os.ModePerm)
		if err != nil {
			return err
		}
		// syscall.Umask(0)

		outputs, err := computeEvent(config, coptions)
		fmt.Println(outputs)
		return err
	}

}

//computeEvent iterates over the program order and requests each plugin to compute the associated model in the model list.
func computeEvent(config Configuration, options compute.Options) (map[string][]float64, error) {
	outputvariablemap := make(map[string][]float64)
	for idx, p := range config.ProgramOrder().Plugins {
		err := p.Compute(config.Models()[idx], options)
		if err != nil {
			return nil, err
		}
		options.EventOptions = options.IncrementModelIndex()

		// if the plugin implements the outputrecorder interface, get outputs
		outputvariables, ok := options.OutputVariables[config.Models()[idx].ModelName()]
		if ok {
			ovp, ovpok := p.(component.OutputReporter)
			if ovpok {
				outputs, err := ovp.ComputeOutputVariables(outputvariables, config.Models()[idx])
				if err != nil {
					panic(err)
				}
				outputvariablemap[config.Models()[idx].ModelName()] = outputs
			}
		}
	}
	return outputvariablemap, nil
}
