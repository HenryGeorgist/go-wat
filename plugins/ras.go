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
)

type RasPlugin struct {
}

type RasBoundaryConditions struct {
	BCLine   string  `json:"bc_line"`
	Interval float64 `json:"interval"`
	Steps    int     `json:"steps"`
	//Flows    []float64 `json:"flows"`
}

// Need some converter to pull this from text
// Using this as a place holder, which will fail on any model that has anything other than 1Hour
var rasIntervals map[string]float64 = map[string]float64{"1HOUR": 1}

func extractTimeInterval(s string) (float64, error) {

	rawText := strings.Trim(s, "\r")
	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return 0, fmt.Errorf("extractTimeInterval error: insufficient data from text file line")
	}

	if _, found := rasIntervals[textLineParts[1]]; !found {
		return 0, fmt.Errorf("extractTimeInterval error: unknown timestep, please add to `rasIntervals` in `ras.go`")
	}

	numericInterval := rasIntervals[textLineParts[1]]
	return numericInterval, nil

}

func extractNumberTimeSteps(s string) (int, error) {

	rawText := strings.Trim(s, "\r")

	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return 0, fmt.Errorf("extractNumberTimeSteps error: unrecognized data from text file line")
	}

	stepsNumeric, err := strconv.Atoi(strings.Trim(textLineParts[1], " "))
	if err != nil {
		return 0, err
	}

	return stepsNumeric, nil
}

func extractBCName(s string) (string, error) {

	rawText := strings.Trim(s, "\r")
	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return "", fmt.Errorf("extractBCName error: unrecognized data from text file line")
	}

	lineValues := strings.Split(textLineParts[1], ",")
	var fullBCName string = strings.Trim(lineValues[0], " ")
	for i, text := range lineValues {
		// Todo: Need to verify how BC's  are stored / nomenclature convention for this
		// currently this function strips white space from `Boundary Location` line and concatenates values with textData (skipping empty spaces)
		// i.e. `Boundary Location=White           ,Muncie          ,15696.24,        ,                ,                ,                , `
		// is returned as `White-Muncie-15696.24`
		textData := strings.Trim(text, " ")
		if i > 0 && textData != "" {
			fullBCName += "-" + textData
		}
	}
	return fullBCName, nil
}

// hecRasBCs is a placeholder utility funciton for reading data from models
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

						bcLineName, err := extractBCName(line)
						if err != nil {
							return rbc, err
						}

						numericInterval, err := extractTimeInterval(lines[i+1])
						if err != nil {
							return rbc, err
						}

						stepsNumeric, err := extractNumberTimeSteps(lines[i+2])
						if err != nil {
							return rbc, err
						}

						rbc.BCLine = bcLineName
						rbc.Interval = numericInterval
						rbc.Steps = stepsNumeric

						rbcs = append(rbcs, rbc)

						fmt.Println("bcLineName", bcLineName)

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
// Input data locations: boundary conditions names
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
			Name:      rbcs.BCLine,
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
	// get model link, read file, pull foats write to model...
	return nil
}
