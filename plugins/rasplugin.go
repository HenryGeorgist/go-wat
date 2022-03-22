package plugins

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-wat/component"
	"go-wat/compute"
	"go-wat/config"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/filestore"
	"github.com/USACE/mcat-ras/tools"
	// "github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

// Need some converter to pull this from text
// Using this as a place holder, which will fail on any model that has anything other than 1Hour
var rasIntervals map[string]float64 = map[string]float64{"1HOUR": 1}

// 	Hard coded (pulled) from Muncie model for demo purposes
var flows []float64 = []float64{13500, 14000, 14500, 15000, 15500, 16000, 16500, 17000, 17500, 18000, 18500, 19000,
	19500, 20000, 20500, 21000, 21000, 21000, 21000, 21000, 21000, 20500, 20000, 19500,
	19000, 18500, 18000, 17500, 17000, 16500, 16000, 15500, 15000, 14583.33, 14166.67, 13750,
	13333.33, 12916.67, 1250012083.33, 11666.67, 11250, 10833.33, 10416.67, 10000, 9666.67, 9333.33,
	9000, 8666.67, 8333.33, 8000, 7666.67, 7333.33, 7000, 6666.67, 6333.33, 6000, 5875, 5750, 5625,
	5500, 5375, 5250, 5125, 5000}

type RasPlugin struct {
}

type RasBoundaryConditions struct {
	BCLine   string    `json:"bc_line"`
	Interval float64   `json:"interval"`
	Steps    int       `json:"steps"`
	Flows    []float64 `json:"flows"`
}

// HecRasBCs is a placeholder utility funciton for reading data from models
func HecRasBCs(rm config.RasModelInfo) (RasBoundaryConditions, error) {

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
				// Todo: make this work on any model, not just Muncy!
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

						rbc.Flows = flows
						rbcs = append(rbcs, rbc)
					} else {
						continue
					}
				}
			}
		}

	}

	// Configured to use the muncy model for demo purposes only
	// grabbing the firstBC info from u01 | p04, using hard coded flows
	return rbcs[0], nil

}

// Name of the model (ras - bc-1)
// ParentPluginName (hyrdoscalaar plugin)
// Flows shape set
// Flow frequency LPIII (or other BootstrappableDistribution)
// Links are output links (what can I produce)
type RasModel struct {
	Name             string                                `json:"name"`
	ParentPluginName string                                `json:"parent_plugin_name"`
	Flows            []float64                             `json:"flows"`
	TimeStep         time.Duration                         `json:"timestep"`
	FlowFrequency    statistics.BootstrappableDistribution `json:"flow_frequency"`
	Links            component.ModelLinks                  `json:"-"`
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

// Get input list of BC's from a plan from HDF...
func (rp RasPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	// no links needed here, this serves as a generator in this context
	ret := make([]component.InputDataLocation, 0)
	return ret
}

func (rp RasPlugin) Name() string {
	return "RAS Plugin"
}

// List output file (p*.hdf)
func (rp RasPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Name:                 "Hydrograph",
		Parameter:            "Flow",
		Format:               "Array",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: rp.Name(),
	}
	ret = append(ret, output)
	return ret
}

func (rp RasPlugin) Compute(model component.Model, options compute.Options) error {
	return nil
}
