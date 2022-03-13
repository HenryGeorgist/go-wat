package simulation

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
	"github.com/HenryGeorgist/go-wat/plugin"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func TestDeterministicSimulation(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 1.0, ParentPluginName: spb.Name()}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
	//create a program order
	plugins := make([]component.Computable, 3)
	plugins[0] = spa
	plugins[1] = spb
	plugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: plugins}
	//model link
	aminputs := ap.InputLinks(am)
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)

	modelLinks := make([]component.Link, 2)
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}
	am.Links = ml
	//set up a model list
	models := make([]component.Model, 3)
	models[0] = sma
	models[1] = smb
	models[2] = am
	//set up a timewindow
	tw := compute.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local)}
	//set up a configuration
	deterministicconfig := DeterministicConfiguration{
		Programorder:      programOrder,
		ModelList:         models,
		TimeWindow:        tw,
		Outputdestination: "/workspaces/go-wat/testdata/",
		Inputsource:       "/workspaces/go-wat/testdata/",
	}
	//compute
	Compute(deterministicconfig)
}
func TestStochasticSimulation(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
	//create a program order
	plugins := make([]component.Computable, 3)
	plugins[0] = spa
	plugins[1] = spb
	plugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: plugins}
	//model link
	aminputs := ap.InputLinks(am)
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)

	modelLinks := make([]component.Link, 2)
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}
	am.Links = ml
	//set up a model list
	models := make([]component.Model, 3)
	models[0] = sma
	models[1] = smb
	models[2] = am
	//set up a timewindow
	tw := compute.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local)}
	//create an event generator
	eg := plugin.AnnualEventGenerator{}
	//set up a configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        5,
		LifecyclesPerRealization: 3,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		Outputdestination:        "/workspaces/go-wat/testdata/",
		Inputsource:              "/workspaces/go-wat/testdata/",
	}
	//compute
	Compute(stochasticconfig)
}
func TestStochasticSimulation_serialization(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
	//create a program order
	plugins := make([]component.Computable, 3)
	plugins[0] = spa
	plugins[1] = spb
	plugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: plugins}
	//model link
	aminputs := ap.InputLinks(am)
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)

	modelLinks := make([]component.Link, 2)
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}
	am.Links = ml
	//set up a model list
	models := make([]component.Model, 3)
	models[0] = sma
	models[1] = smb
	models[2] = am
	//set up a timewindow
	tw := compute.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local)}
	//create an event generator
	eg := plugin.AnnualEventGenerator{}
	//set up a configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        1,
		LifecyclesPerRealization: 1,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		Outputdestination:        "/workspaces/go-wat/testdata/",
		Inputsource:              "/workspaces/go-wat/testdata/",
	}
	//compute
	bytes, err := json.Marshal(stochasticconfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
}
func TestStochasticSimulation_withHydrograph(t *testing.T) {
	//create a hydrograph scaler plugin
	hsp := plugin.HydrographScalerPlugin{}
	//create a hydrograph scaler model
	flows := []float64{1.0, 5.0, 2.0}
	flow_frequency := statistics.LogPearsonIIIDistribution{Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10}
	hsm := plugin.HydrographScalerModel{
		Name:             "RAS Boundary",
		ParentPluginName: hsp.Name(),
		TimeStep:         time.Hour,
		Flows:            flows,
		FlowFrequency:    flow_frequency,
	}

	//create a program order
	plugins := make([]component.Computable, 1)
	plugins[0] = hsp
	programOrder := component.ProgramOrder{Plugins: plugins}
	//model link
	modelLinks := make([]component.Link, 0)
	ml := component.ModelLinks{Links: modelLinks}
	hsm.Links = ml
	//set up a model list
	models := make([]component.Model, 1)
	models[0] = hsm
	//set up a timewindow
	tw := compute.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2018, 1, 1, 4, 1, 1, 1, time.Local)}
	//create an event generator
	eg := plugin.AnnualEventGenerator{}
	//set up a configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        1,
		LifecyclesPerRealization: 1,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		Outputdestination:        "/workspaces/go-wat/testdata/",
		Inputsource:              "/workspaces/go-wat/testdata/",
	}
	bytes, _ := json.Marshal(stochasticconfig)
	fmt.Println(string(bytes))
	//compute
	Compute(stochasticconfig)
}
