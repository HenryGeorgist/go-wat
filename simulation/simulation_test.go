package simulation

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-wat/component"
	"go-wat/compute"
	"go-wat/config"
	"go-wat/plugins"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func TestDeterministicSimulation(t *testing.T) {

	testSettings, err := config.LoadTestSettings()
	fmt.Println("testSettings", testSettings)
	if err != nil {
		t.Fatal(err)
	}

	//create two scalar plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}

	//create two scalar models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 1.0, ParentPluginName: spb.Name()}

	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

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
		Outputdestination: testSettings.OutputDataDir,
		Inputsource:       testSettings.InputDataDir,
	}

	// Compute
	err = Compute(deterministicconfig)
	if err != nil {
		t.Fatal(err)
	}
}
func TestStochasticSimulation(t *testing.T) {

	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	//create two scalar plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}

	//create two scalar models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}

	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	//create a program order
	activePlugins := make([]component.Computable, 3)
	activePlugins[0] = spa
	activePlugins[1] = spb
	activePlugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: activePlugins}

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
	eg := plugins.AnnualEventGenerator{}

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
		Outputdestination:            testSettings.OutputDataDir,
		Inputsource:                  testSettings.InputDataDir,
		DeleteOutputAfterRealization: false,
	}

	// Compute
	err = Compute(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	}

}
func TestStochasticSimulation_serialization(t *testing.T) {

	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	//create two scalar plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}

	//create two scalar models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}

	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	//create a program order
	activePlugins := make([]component.Computable, 3)
	activePlugins[0] = spa
	activePlugins[1] = spb
	activePlugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: activePlugins}

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
	eg := plugins.AnnualEventGenerator{}

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
		Outputdestination:        testSettings.OutputDataDir,
		Inputsource:              testSettings.InputDataDir,
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

	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	//create a hydrograph scaler plugin
	hsp := plugins.HydrographScalerPlugin{}

	//create a hydrograph scaler model
	flows := []float64{1.0, 5.0, 2.0, 15.0}
	flow_frequency := statistics.LogPearsonIIIDistribution{Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10}

	hsm := plugins.HydrographScalerModel{
		Name:             "RASBoundary",
		ParentPluginName: hsp.Name(),
		TimeStep:         time.Hour,
		Flows:            flows,
		FlowFrequency:    flow_frequency,
	}

	//create a program order
	activePlugins := make([]component.Computable, 1)
	activePlugins[0] = hsp

	programOrder := component.ProgramOrder{Plugins: activePlugins}

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
	eg := plugins.AnnualEventGenerator{}

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
		Outputdestination:        testSettings.OutputDataDir,
		Inputsource:              testSettings.InputDataDir,
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

func TestStochasticSimulation_withRAS(t *testing.T) {
	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Get RAS BC from the Munice model
	rasBCs, err := plugins.HecRasBCs(testSettings.RasModel)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("rasBCs.........", rasBCs)

	// //create a hydrograph scaler plugin
	// hsp := plugins.HydrographScalerPlugin{}
	// rp := plugins.RasPlugin{}

	// flow_frequency := statistics.LogPearsonIIIDistribution{Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10}

	// hsm := plugins.HydrographScalerModel{
	// 	Name:             "RASBoundary",
	// 	ParentPluginName: hsp.Name(),
	// 	TimeStep:         time.Hour,
	// 	Flows:            flows,
	// 	FlowFrequency:    flow_frequency,
	// }

	// rm := plugins.RasModel{
	// 	Name:             "Muncie",
	// 	ParentPluginName: rp.Name(),
	// 	TimeStep:         time.Hour,
	// 	Flows:            flows,
	// 	FlowFrequency:    flow_frequency,
	// }

	// //create a program order
	// activePlugins := make([]component.Computable, 2)
	// activePlugins[0] = hsp
	// activePlugins[1] = rp

	// programOrder := component.ProgramOrder{Plugins: activePlugins}

	// //model link
	// rasinputs := rp.InputLinks(rm)
	// hsmoutputs := hsp.OutputLinks(hsm)
	// rasoutputs := rp.OutputLinks(rm)

	// numBoundaries := len(rasinputs)
	// modelLinks := make([]component.Link, numBoundaries)
	// modelLinks[0] = component.Link{InputDataLocation: rasinputs[0], OutputDataLocation: hsmoutputs[0]}
	// // uncomment and add depneding on # of BC's
	// modelLinks[1] = component.Link{InputDataLocation: rasinputs[0], OutputDataLocation: rasoutputs[0]}

	// ml := component.ModelLinks{Links: modelLinks}
	// rm.Links = ml
	// hsm.Links = ml

	// //set up a model list
	// models := make([]component.Model, 2)
	// models[0] = hsm
	// models[1] = rm

	// //set up a timewindow
	// tw := compute.TimeWindow{StartTime: time.Date(2006, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2068, time.December, 31, 1, 1, 1, 1, time.Local)}

	// //create an event generator
	// eg := plugins.AnnualEventGenerator{}

	// //set up a configuration
	// stochasticconfig := StochasticConfiguration{
	// 	Programorder:             programOrder,
	// 	ModelList:                models,
	// 	EventGenerator:           eg,
	// 	LifecycleTimeWindow:      tw,
	// 	TotalRealizations:        1,
	// 	LifecyclesPerRealization: 1,
	// 	InitialRealizationSeed:   1234,
	// 	InitialEventSeed:         1234,
	// 	Outputdestination:        testSettings.OutputDataDir,
	// 	Inputsource:              testSettings.InputDataDir,
	// }

	// bytes, err := json.Marshal(stochasticconfig)
	// if err != nil {
	// 	t.Fatal(err)
	// } else {
	// 	fmt.Println("StochasticConfiguration: ", string(bytes))
	// }

	// // Compute
	// err = Compute(stochasticconfig)
	// if err != nil {
	// 	t.Fatal(err)
	// }

}
