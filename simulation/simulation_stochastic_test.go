package simulation

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-wat/component"
	"go-wat/config"
	"go-wat/option"
	"go-wat/plugins"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func TestStochasticSimulation(t *testing.T) {

	// Load Configuration data
	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(testSettings)

	// Instantiate required plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}
	ap := plugins.AddPlugin{}

	// Create a program execution order
	activePlugins := make([]component.Computable, 3)
	activePlugins[0] = spa
	activePlugins[1] = spb
	activePlugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: activePlugins}

	// Create simulation models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	// Assign links to pair models with plugins
	smaoutput := spa.OutputLinks(sma) // Outputs from ValueA model
	smboutput := spa.OutputLinks(smb) // Outputs from ValueB model
	aminputs := ap.InputLinks(am)     // Inputs for APlusB model

	//  Associate model dependency as needed
	modelLinks := make([]component.Link, 2)

	// Independent Models
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}

	// Dependent Models
	am.Links = ml

	// Create list of models (in order of dependency)
	models := make([]component.Model, 3)
	models[0] = sma
	models[1] = smb
	models[2] = am

	// Options

	// Use a timewindow
	tw := option.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}

	// Use an event generator
	eg := option.AnnualEventGenerator{}

	// Assign configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:                 programOrder,
		ModelList:                    models,
		EventGenerator:               eg,
		LifecycleTimeWindow:          tw,
		TotalRealizations:            2,
		LifecyclesPerRealization:     1,
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

	// Load Configuration data
	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Instantiate required plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}
	ap := plugins.AddPlugin{}

	// Create a program execution order
	activePlugins := make([]component.Computable, 3)
	activePlugins[0] = spa
	activePlugins[1] = spb
	activePlugins[2] = ap
	programOrder := component.ProgramOrder{Plugins: activePlugins}

	// Create simulation models with test data
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	// Assign links to pair models with plugins
	smaoutput := spa.OutputLinks(sma) // Outputs from ValueA model
	smboutput := spa.OutputLinks(smb) // Outputs from ValueB model
	aminputs := ap.InputLinks(am)     // Inputs for APlusB model

	//  Associate model dependency as needed
	modelLinks := make([]component.Link, 2)

	// Independent Models
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}

	// Dependent Models
	am.Links = ml

	// Create list of models (in order of dependency)
	models := make([]component.Model, 3)
	models[0] = sma
	models[1] = smb
	models[2] = am

	// Options

	// Use a timewindow
	tw := option.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}

	// Use an event generator
	eg := option.AnnualEventGenerator{}

	// Assign configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        2,
		LifecyclesPerRealization: 1,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		Outputdestination:        testSettings.OutputDataDir,
		Inputsource:              testSettings.InputDataDir,
	}

	bytes, err := json.Marshal(stochasticconfig)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(bytes))
	}

	// Compute ?
	// err = Compute(stochasticconfig)
	// if err != nil {
	// 	t.Fatal(err)
	// }

}

func TestStochasticSimulation_withHydrograph(t *testing.T) {

	// Load Configuration data
	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Instantiate required plugins
	hsp := plugins.HydrographScalerPlugin{}

	// Create a program execution order
	activePlugins := make([]component.Computable, 1)
	activePlugins[0] = hsp
	programOrder := component.ProgramOrder{Plugins: activePlugins}

	// Create simulation models with test data
	hsm := plugins.HydrographScalerModel{
		Name:             "RASBoundary",
		ParentPluginName: hsp.Name(),
		TimeStep:         time.Hour,
		Flows:            []float64{1.0, 5.0, 2.0, 15.0},
		FlowFrequency:    statistics.LogPearsonIIIDistribution{Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10},
	}

	// Assign links to pair models with plugins
	// no assignment needed...

	//  Associate model dependency as needed
	modelLinks := make([]component.Link, 0)

	// Independent Models
	ml := component.ModelLinks{Links: modelLinks}

	// Dependent Models
	hsm.Links = ml

	// Create list of models (in order of dependency)
	models := make([]component.Model, 1)
	models[0] = hsm

	// Options

	// Use a timewindow
	tw := option.TimeWindow{StartTime: time.Date(2006, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2068, time.December, 31, 1, 1, 1, 1, time.Local)}

	// Use an event generator
	eg := option.AnnualEventGenerator{}

	// Assign configuration
	stochasticconfig := StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        2,
		LifecyclesPerRealization: 1,
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

	// Load Configuration data
	testSettings, err := config.LoadTestSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Instantiate required plugins
	hsp := plugins.HydrographScalerPlugin{}
	rp := plugins.RasPlugin{}

	// Create a program execution order
	activePlugins := make([]component.Computable, 2)
	activePlugins[0] = hsp
	activePlugins[1] = rp
	programOrder := component.ProgramOrder{Plugins: activePlugins}

	rm, err := plugins.NewRasModel(testSettings.RasModel)
	if err != nil {
		t.Fatal(err)
	}

	hsm, err := plugins.NewHydrographScalerModel(testSettings.HydroModel)
	if err != nil {
		t.Fatal(err)
	}

	// Assign links to pair models with plugins
	rminputs := rp.InputLinks(rm)
	hsmoutputs := hsp.OutputLinks(hsm)
	// rmoutputs := rp.OutputLinks(rm)

	//  Associate model dependency as needed
	modelLinks := make([]component.Link, 1)

	// Independent Models
	modelLinks[0] = component.Link{InputDataLocation: rminputs[0], OutputDataLocation: hsmoutputs[0]}
	ml := component.ModelLinks{Links: modelLinks}

	// Dependent Models
	rm.Links = ml

	// Create List of models (in order of dependency)
	models := make([]component.Model, 2)
	models[0] = hsm
	models[1] = rm

	// Options

	// Use a timewindow
	tw := option.TimeWindow{StartTime: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}

	// Use an event generator
	eg := option.AnnualEventGenerator{}

	// Use a configuration: Deterministic --or-- stochastic
	stochasticconfig := StochasticConfiguration{
		Programorder:                 programOrder,
		ModelList:                    models,
		EventGenerator:               eg,
		LifecycleTimeWindow:          tw,
		TotalRealizations:            1,
		LifecyclesPerRealization:     2,
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
