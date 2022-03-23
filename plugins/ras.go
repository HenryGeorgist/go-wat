package plugins

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"go-wat/component"
	"go-wat/config"
	"go-wat/option"

	"github.com/USACE/filestore"
	"github.com/USACE/mcat-ras/tools"
	// "github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

// Need some converter to pull this from text
// Using this as a place holder, which will fail on any model that has anything other than 1Hour
var rasIntervals map[string]float64 = map[string]float64{"1HOUR": 1}

type RasPlugin struct {
}

type RasBoundaryConditions struct {
	BCLine   string  `json:"bc_line"`
	Interval float64 `json:"interval"`
	Steps    int     `json:"steps"`
	//Flows    []float64 `json:"flows"`
}

// HecRasBCs is a placeholder utility funciton for reading data from models
func hecRasBCs(rm config.RasModelInfo) (RasBoundaryConditions, error) {

	var rbc RasBoundaryConditions
	fs, err := filestore.NewFileStore(filestore.BlockFSConfig{})
	if err != nil {
		return rbc, err
	}

	modelData, err := tools.NewRasModel(rm.ProjectFilePath, fs)
	if err != nil {
		return rbc, err
	}

	var rbcs []RasBoundaryConditions

	for _, file := range modelData.Metadata.PlanFiles {

		lineBytes, err := ioutil.ReadFile(rm.BasePath + "." + file.FlowFile)
		if err != nil {
			return rbc, err
		}

		lines := strings.Split(string(lineBytes), "\n")

		for i, line := range lines {
			match, err := regexp.MatchString("=", line)
			if err != nil {
				return rbc, err
			}

			if match {
				lineData := strings.Split(line, "=")

				switch lineData[0] {
				// Todo: make this work on any model, not just muncie!
				case "Boundary Location":
					nextLine := strings.Split(lines[i+1], "=")[0]

					if nextLine == "Interval" {

						rbc.BCLine = strings.Trim(line, "\r")

						intervalText := strings.Trim(strings.Split(lines[i+1], "=")[1], " \r")
						numericInterval := rasIntervals[intervalText]
						rbc.Interval = numericInterval

						stepsText := strings.Trim(strings.Split(lines[i+2], "=")[1], " \r")
						stepsNumeric, err := strconv.Atoi(stepsText)
						if err != nil {
							return rbc, err
						}
						rbc.Steps = stepsNumeric

						//rbc.Flows = flows
						rbcs = append(rbcs, rbc)
					} else {
						continue
					}
				}
			}
		}

	}

	// Configured to use the muncie model for demo purposes only
	// grabbing the firstBC info from u01 | p04, using hard coded flows
	return rbcs[0], nil

}

// Name of the model (Muncie)
// ParentPluginName (ras plugin)
// Links are output links (what can I produce)
type RasModel struct {
	Name             string               `json:"name"`
	BasePath         string               `json:"basepath"`
	ProjectFilePath  string               `json:"projectfile"`
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

// Get input list of BC's from a plan u file...
func (rp RasPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 1)
	rm, rmok := model.(RasModel)
	if rmok {
		rbcs, err := hecRasBCs(config.RasModelInfo{
			BasePath:        rm.BasePath,
			ProjectFilePath: rm.ProjectFilePath,
		})
		if err != nil {
			panic(err)
		}
		idl := component.InputDataLocation{
			Name:      rm.Name + " " + rbcs.BCLine,
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

// List output file (p*.hdf)
func (rp RasPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Name:                 model.ModelName() + " output hdf file",
		Parameter:            "RAS output",
		Format:               "HDF",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.hdf", model.ModelName())}, //this is not quite right
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: rp.Name(),
	}
	ret = append(ret, output)
	return ret
}

func (rp RasPlugin) Compute(model component.Model, options option.Options) error {
	return nil
}
