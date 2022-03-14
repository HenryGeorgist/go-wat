package simulation

import (
	"fmt"
	"os"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
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
	stochastic, ok := config.(StochasticConfiguration)
	var coptions compute.Options
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
			realizationInputPath := rootinputPath + "realization " + fmt.Sprint(realization) + "/"
			realizationOutputPath := rootOutputPath + "realization " + fmt.Sprint(realization) + "/"
			//loop for lifecycles
			realizationSeeds := realizationRandom.GeneratePluginSeeds() //probably make one per model
			fmt.Println(fmt.Sprintf("Computing realization %v", realization))
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
						OutputVariables:   outputvariablesMap,
					}
					outputvars, err := computeEvent(config, coptions)
					if err != nil {
						panic(err)
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
		//if the plugin implements the outputrecorder interface, get outputs
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
