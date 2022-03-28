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
	Pfile            string               `json:"planfile"`
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

// for a ras plugin the major output is the hdf file, many post processes require the terrain and the hdf.
//down the road, we may want to offer additional outputs, like stage or flow at a cross section from the unsteady plan specification
//or stored map options from the rasmapper specification.
func (rp RasPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	rm, rmok := model.(RasModel)
	if rmok {
		_, planFile := filepath.Split(rm.Pfile)
		terrain := component.OutputDataLocation{
			Name:                 "Terrain",
			Parameter:            "Elevation",
			Format:               "geotif",
			LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/Terrain/%v.tif", model.ModelName())}, //this is not quite right
			GeneratingModelName:  model.ModelName(),
			GeneratingPluginName: rp.Name(),
		}
		rasHDFOutput := component.OutputDataLocation{
			Name:                 "HEC-RAS HDF output",
			Parameter:            "Depth",
			Format:               "HDF",
			LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.hdf", planFile)},
			GeneratingModelName:  model.ModelName(),
			GeneratingPluginName: rp.Name(),
		}
		ret = append(ret, terrain)
		ret = append(ret, rasHDFOutput)
	}

	return ret
}

func (rp RasPlugin) Compute(model component.Model, options option.Options) error {

	rm, rmok := model.(RasModel)
	if !rmok {
		return fmt.Errorf("not a ras model")
	}

	links := model.ModelLinkages()
	// requiredSims := make([]ContainerParams, 0)

	for _, link := range links.Links {

		var simContainerParams ContainerParams

		lcsv, linkok := link.OutputDataLocation.LinkInfo.(component.LocalCSVLink)
		if linkok {

			// read Hydrologic Sampler output provided by link
			hsmOutputCSVFile := options.InputSource + lcsv.Path
			lineBytes, err := ioutil.ReadFile(hsmOutputCSVFile)
			if err != nil {
				return err
			}

			formattedRasData, err := hydroArrayToRasFormat(lineBytes)
			if err != nil {
				return err
			}

			// Parse data
			inputFlowFile := rm.Ufile.Path
			_, flowFile := filepath.Split(rm.Ufile.Path)
			_, planFile := filepath.Split(rm.Pfile)
			outputFlowFile := fmt.Sprintf("%v/%v", options.OutputDestination, flowFile)

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
			simContainerParams.InputRasModelDir = rm.BasePath
			simContainerParams.ModelName = rm.Name
			simContainerParams.UpdatedFlowFile = outputFlowFile
			simContainerParams.PlanFile = planFile
			simContainerParams.OutputHDF = fmt.Sprintf("%v/%v.hdf", options.OutputDestination, planFile)

		}

		// // for testing without ras sim
		// _, err := runSimInContainerPreview(simContainerParams)
		// if err != nil {
		// 	return err
		// }

		_, err := RunSimInContainer(simContainerParams)
		//TODO: ensure the terrain and the hdf are both stored in the outputdestination.
		if err != nil {
			return err
		}

	}

	return nil
}
