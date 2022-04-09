package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go-wat/component"
	"go-wat/config"
	"go-wat/option"
	"go-wat/plugins"
	"go-wat/simulation"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func LoadConfig(configFile string) (config.EventSettings, error) {

	var ts config.EventSettings
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return ts, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ts, err
	}

	json.Unmarshal(jsonData, &ts)
	userRootDir := filepath.FromSlash(ts.UserHomeDir)
	ts.UserHomeDir = userRootDir
	ts.InputDataDir = filepath.FromSlash(fmt.Sprintf("%v/%v/", userRootDir, ts.InputDataDir))
	ts.OutputDataDir = filepath.FromSlash(fmt.Sprintf("%v/%v/", userRootDir, ts.OutputDataDir))
	ts.HydroModel = filepath.FromSlash(fmt.Sprintf("%v/%v", userRootDir, ts.HydroModel))

	return ts, nil

}

func main() {
	fmt.Println("welcome to go-wat")

	var configPath string
	flag.StringVar(&configPath, "config", "", "please specify an input file using `-config=myconfig.json`")
	flag.Parse()

	if configPath == "" {
		fmt.Println("given path...", string(configPath))
		fmt.Println("please specify an input file using `-config=myconfig.json`")
		return
	} else {
		if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
			fmt.Println("input file does not exist or is inaccessible")
			return
		}

	}

	// Load Configuration data
	eventSettings, err := LoadConfig(configPath)
	if err != nil {
		fmt.Println("error", eventSettings)
	}

	// Instantiate required plugins
	hsp := plugins.HydrographScalerPlugin{}

	// Create a program execution order
	activePlugins := make([]component.Computable, 1)
	activePlugins[0] = hsp
	programOrder := component.ProgramOrder{Plugins: activePlugins}

	hsm, err := plugins.NewHydrographScalerModel(eventSettings.HydroModel)
	if err != nil {
		fmt.Println("error", eventSettings)
	}

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
	tw := option.TimeWindow{StartTime: time.Date(2006, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2007, time.December, 31, 1, 1, 1, 1, time.Local)}

	// Use an event generator
	eg := option.AnnualEventGenerator{}

	// Assign configuration
	stochasticconfig := simulation.StochasticConfiguration{
		Programorder:             programOrder,
		ModelList:                models,
		EventGenerator:           eg,
		LifecycleTimeWindow:      tw,
		TotalRealizations:        eventSettings.TotalRealizations,
		LifecyclesPerRealization: eventSettings.LifecyclesPerRealization,
		InitialRealizationSeed:   eventSettings.InitialRealizationSeed,
		InitialEventSeed:         eventSettings.InitialEventSeed,
		Outputdestination:        eventSettings.OutputDataDir,
		Inputsource:              eventSettings.InputDataDir,
	}

	bytes, err := json.Marshal(stochasticconfig)
	if err != nil {
		fmt.Println("error....StochasticConfiguration: ", string(bytes))
	}

	// Compute
	fileOutputs, err := simulation.Compute(stochasticconfig)
	if err != nil {
		fmt.Println("Compute: ", err)
	}

	fmt.Println(fileOutputs)

}
