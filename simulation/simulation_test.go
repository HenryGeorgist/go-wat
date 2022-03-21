package simulation

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-wat/component"
	"go-wat/compute"
	"go-wat/plugin"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

const (
	inputDataDir  string = "/workspaces/test-data/"
	outputDataDir string = "/workspaces/test-data/"
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
		Outputdestination: outputDataDir,
		Inputsource:       inputDataDir,
	}

	// Compute
	err := Compute(deterministicconfig)
	if err != nil {
		t.Fatal(err)
	}
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
	tw := compute.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2068, time.December, 31, 1, 1, 1, 1, time.Local)}

	//create an event generator
	eg := plugin.AnnualEventGenerator{}

	//set up a configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:                 programOrder,
		ModelList:                    models,
		EventGenerator:               eg,
		LifecycleTimeWindow:          tw,
		TotalRealizations:            5,
		LifecyclesPerRealization:     3,
		InitialRealizationSeed:       1234,
		InitialEventSeed:             1234,
		Outputdestination:            outputDataDir,
		Inputsource:                  inputDataDir,
		DeleteOutputAfterRealization: true,
	}

	// Compute
	err := Compute(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	}

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
		Outputdestination:        outputDataDir,
		Inputsource:              inputDataDir,
	}

	//compute
	bytes, err := json.Marshal(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(bytes))
	}

}
func TestStochasticSimulation_withHydrograph(t *testing.T) {
	//create a hydrograph scaler plugin
	hsp := plugin.HydrographScalerPlugin{}

	//create a hydrograph scaler model
	flows := []float64{1.0, 5.0, 2.0, 15.0}
	flow_frequency := statistics.LogPearsonIIIDistribution{Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10}

	hsm := plugin.HydrographScalerModel{
		Name:             "RASBoundary",
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
	tw := compute.TimeWindow{StartTime: time.Date(2006, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2068, time.December, 31, 1, 1, 1, 1, time.Local)}

	//create an event generator
	eg := plugin.AnnualEventGenerator{}

	//set up a configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        3,
		LifecyclesPerRealization: 2,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		Outputdestination:        outputDataDir,
		Inputsource:              inputDataDir,
	}

	bytes, err := json.Marshal(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("StochasticConfiguration: ", string(bytes))
	}

	// Compute
	err = Compute(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	}

}
