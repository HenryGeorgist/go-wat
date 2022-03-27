package simulation

import (
	"fmt"
	"testing"
	"time"

	"go-wat/component"
	"go-wat/config"
	"go-wat/option"
	"go-wat/plugins"
)

func TestDeterministicSimulation(t *testing.T) {

	// Load Configuration data
	testSettings, err := config.LoadTestSettings()
	fmt.Println("testSettings", testSettings)
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

	// Create simulation models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 1.0, ParentPluginName: spb.Name()}
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	// Assign links to pair models with plugins
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)
	aminputs := ap.InputLinks(am)

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
	tw := option.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local)}

	// Assign configuration
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
