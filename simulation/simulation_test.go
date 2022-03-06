package simulation

import (
	"testing"
	"time"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
	"github.com/HenryGeorgist/go-wat/plugin"
)

func TestDeterministicSimulation(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPlugin: spa}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPlugin: spb}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPlugin: ap}
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

	modelLinks := make(map[component.InputDataLocation]component.OutputDataLocation)
	modelLinks[aminputs[0]] = smaoutput[0]
	modelLinks[aminputs[1]] = smboutput[0]
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
		programOrder:      programOrder,
		models:            models,
		TimeWindow:        tw,
		outputDestination: "/workspaces/go-wat/testdata/",
		inputSource:       "/workspaces/go-wat/testdata/",
	}
	//compute
	Compute(deterministicconfig)
}
func TestStochasticSimulation(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPlugin: spa}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPlugin: spb}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPlugin: ap}
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

	modelLinks := make(map[component.InputDataLocation]component.OutputDataLocation)
	modelLinks[aminputs[0]] = smaoutput[0]
	modelLinks[aminputs[1]] = smboutput[0]
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
		programOrder:             programOrder,
		models:                   models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        1,
		LifecyclesPerRealization: 1,
		InitialRealizationSeed:   1234,
		InitialEventSeed:         1234,
		outputDestination:        "/workspaces/go-wat/testdata/",
		inputSource:              "/workspaces/go-wat/testdata/",
	}
	//compute
	Compute(stochasticconfig)
}
