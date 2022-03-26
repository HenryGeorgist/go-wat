package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"go-wat/component"
	"go-wat/option"
)

type RasPlugin struct {
}

// Name of the model (Muncie)
// ParentPluginName (ras plugin)
// Links are output links (what can I produce)
type RasModel struct {
	Name             string               `json:"name"`
	BasePath         string               `json:"basepath"`
	ProjectFilePath  string               `json:"projectfile"`
	UfilePath        string               `json:"unsteadyfile"`
	ParentPluginName string               `json:"parent_plugin_name"`
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

	for _, link := range links.Links {

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
			// write ufile
			// Todo: remove this hardcoded file ext
			inputFlowFile := rm.UfilePath //"/home/slawler/workbench/repos/go-wat/sample-data/Muncie/Muncie.u01"
			outputFile := options.OutputDestination + "/Muncie.u01"

			linestart := 9
			linestop := 14

			// temp solution
			fileBytes, err := ioutil.ReadFile(inputFlowFile)

			if err != nil {
				return err
			}

			flowDataLines := strings.Split(string(fileBytes), "\n")
			newFlowFileData := ""
			for i, line := range flowDataLines {
				if i < linestart-1 {
					newFlowFileData += line
				} else if i == linestart {
					newFlowFileData += formattedRasData
				} else if i > linestop {
					newFlowFileData += line
				} else {
					continue
				}

			}

			fmt.Println("Link ok....", lcsv)
			f, err := os.Create(outputFile)
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = f.WriteString(newFlowFileData)

			if err != nil {
				return err

			}

		}

	}
	return nil
}
