package plugin

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

type HydrographScalerPlugin struct {
}
type HydrographScalerModel struct {
	Name             string                            `json:"name"`
	ParentPluginName string                            `json:"parent_plugin_name"`
	Flows            []float64                         `json:"flows"`
	TimeStep         time.Duration                     `json:"timestep"`
	FlowFrequency    statistics.ContinuousDistribution `json:"flow_frequency"`
	Links            component.ModelLinks              `json:"-"`
}

//model implementation
func (hsm HydrographScalerModel) ModelName() string {
	return hsm.Name
}
func (hsm HydrographScalerModel) PluginName() string {
	return hsm.ParentPluginName
}
func (hsm HydrographScalerModel) ModelLinkages() component.ModelLinks {
	return hsm.Links
}

//plugin implementation
//plugin helper function.
func (hsp HydrographScalerPlugin) MarshalJSON() ([]byte, error) {
	ret := "{\"plugin_name\":\"" + hsp.Name() + "\"}"
	return []byte(ret), nil
}
func (hsp HydrographScalerPlugin) Name() string {
	return "Hydrograph Scaling Plugin"
}
func (hsp HydrographScalerPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	return ret
}
func (hsp HydrographScalerPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Name:                 "Hydrograph",
		Parameter:            "Flow",
		Format:               "Array",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: hsp.Name(),
	}
	ret = append(ret, output)
	return ret
}
func (hsp HydrographScalerPlugin) Compute(model component.Model, options compute.Options) error {
	hsm, hsmok := model.(HydrographScalerModel)
	if hsmok {
		value := 1.0
		stochastic, ok := options.EventOptions.(compute.StochasticEventOptions)
		if ok {
			//use a seed!
			r := rand.New(rand.NewSource(stochastic.EventSeed))
			value = hsm.FlowFrequency.InvCDF(r.Float64())
		}
		//write it to the output destination in some agreed upon format?
		outputs := hsp.OutputLinks(model)
		for _, o := range outputs {
			lcsv, _ := o.LinkInfo.(component.LocalCSVLink)
			outputdest := options.OutputDestination + lcsv.Path
			w, err := os.OpenFile(outputdest, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
			if err != nil {
				fmt.Println(err)
			}
			defer w.Close()
			currentTime := options.TimeWindow().StartTime
			fmt.Fprintln(w, "Time,Flow")
			for _, flow := range hsm.Flows {
				if options.TimeWindow().EndTime.After(currentTime) {
					fmt.Fprintln(w, fmt.Sprintf("%v,%v", currentTime, flow*value))
					currentTime = currentTime.Add(hsm.TimeStep)
				} else {
					fmt.Println("encountered more flows than the time window.")
				}
			}
			//what if the number of flows is not big enough for the whole time window? add zeros?
			//
		}
		return nil
	}
	return errors.New("could not convert model into a hydrograph scaling model")
}
