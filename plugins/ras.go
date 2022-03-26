package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go-wat/component"
	"go-wat/option"
)

type RasPlugin struct {
}

type RasModel struct {
	Name             string               `json:"name"`
	BasePath         string               `json:"basepath"`
	ProjectFilePath  string               `json:"projectfile"`
	ParentPluginName string               `json:"parent_plugin_name"`
	Ufile            RasFlowFile          `json:"unsteadyfile"`
	Links            component.ModelLinks `json:"-"`
}

//model implementation
func (rm RasModel) ModelName() string {
	return rm.Name
}

func (rm RasModel) PluginName() string {
	return rm.ParentPluginName
}

// Todo: Write this first
func (rm RasModel) ModelLinkages() component.ModelLinks {
	return rm.Links
}

// Input data locations: boundary conditions names
func (rp RasPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 1)
	rm, rmok := model.(RasModel)
	if rmok {
		idl := component.InputDataLocation{
			Name:      rm.BasePath,
			Parameter: "flow",
			Format:    "csv",
		}
		ret[0] = idl
	}
	return ret
}

func (rp RasPlugin) Name() string {
	return "RAS Plugin"
}

// Update to U file
func (rp RasPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Name:                 model.ModelName() + " output u file",
		Parameter:            "RAS output",
		Format:               "txt",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.txt", model.ModelName())}, //this is not quite right
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: rp.Name(),
	}
	ret = append(ret, output)
	return ret
}

func (rp RasPlugin) Compute(model component.Model, options option.Options) error {

	rm, rmok := model.(RasModel)
	if !rmok {
		return fmt.Errorf("not a ras model")
	}

	links := model.ModelLinkages()
	requiredSims := make([]ContainerParams, 0)

	for _, link := range links.Links {

		var simContainerParams ContainerParams

		lcsv, linkok := link.OutputDataLocation.LinkInfo.(component.LocalCSVLink)
		if linkok {

			// read Hydrologic Sampler output provided by link
			hsmOutputFile := options.InputSource + lcsv.Path
			lineBytes, err := ioutil.ReadFile(hsmOutputFile)
			if err != nil {
				return err
			}

			formattedRasData, err := hydroArrayToRasFormat(lineBytes)
			if err != nil {
				return err
			}

			// Parse data
			inputFlowFile := rm.Ufile.Path //"/home/slawler/workbench/repos/go-wat/sample-data/Muncie/Muncie.u01"
			_, flowFile := filepath.Split(rm.Ufile.Path)
			outputFlowFile := options.OutputDestination + flowFile

			//Todo: Make this dynamic, BCLines[0] works for Muncie, where theres is only 1 bcline
			// Should be able to add an iterator for these cases, replacing 0 with i.
			newFlowFileData, err := newUfile(inputFlowFile, formattedRasData,
				rm.Ufile.BCLines[0].OridnatesLocation.Start,
				rm.Ufile.BCLines[0].OridnatesLocation.Stop)

			f, err := os.Create(outputFlowFile)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.WriteString(newFlowFileData)
			if err != nil {
				return err
			}

			// Todo: read these from json...
			simContainerParams.InputRasModelDir = "/home/slawler/workbench/repos/ras-container/sample-model"
			simContainerParams.ModelName = "Muncie"
			simContainerParams.PlanFile = "Muncie.p04"
			simContainerParams.OutputHDF = "/home/slawler/workbench/repos/go-wat/test-data/realization-0/lifecycle-0/event-0/Muncie.p04.hdf"

			requiredSims = append(requiredSims, simContainerParams)
		}

		// Todo: Switch this to parallel Compute here, start with one requiredSims[0]
		// for testing, then add buffered channel
		_, err := RunSimInContainer(requiredSims[0])
		if err != nil {
			return err
		}
	}
	return nil
}
