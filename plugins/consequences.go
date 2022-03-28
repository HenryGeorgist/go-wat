package plugins

import (
	"fmt"
	"log"

	"go-wat/component"
	"go-wat/option"

	"github.com/USACE/go-consequences/hazardproviders"
)

type ConsequencesPlugin struct {
}
type StructureSource string

const (
	NsiApi StructureSource = "NSI API"
	Gpkg   StructureSource = "Geopasckage"
	Shp    StructureSource = "Shapefile"
)

type OutputType string

const (
	JsonOutput    OutputType = "JSON"
	GpkgOutput    OutputType = "Geopackage"
	ShpOutput     OutputType = "Shapefile"
	SummaryOutput OutputType = "Summary CSV"
)

type ConsequencesModel struct {
	Name             string               `json:"name"`
	ParentPluginName string               `json:"parent_plugin_name"`
	Links            component.ModelLinks `json:"-"`
	StructureSource  StructureSource      `json:"structuresource"`
	StructurePath    string               `json:"structure_path"`
	OutputType       OutputType           `json:"output_type"`
}

func (cm ConsequencesModel) ModelName() string {
	return cm.Name
}

func (cm ConsequencesModel) PluginName() string {
	return cm.ParentPluginName
}
func (cm ConsequencesModel) ModelLinkages() component.ModelLinks {
	return cm.Links
}
func (cp ConsequencesPlugin) MarshalJSON() ([]byte, error) {
	ret := "{\"plugin_name\":\"" + cp.Name() + "\"}"
	return []byte(ret), nil
}

func (cp ConsequencesPlugin) Name() string {
	return "Consequences Plugin"
}

func (cp ConsequencesPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	hp := component.InputDataLocation{
		Name:      "HydraulicProvider",
		Parameter: "depth",
		Format:    "hdf",
	}
	terrain := component.InputDataLocation{
		Name:      "Terrain Grid",
		Parameter: "elev",
		Format:    "grid",
	}
	ret = append(ret, hp)
	ret = append(ret, terrain)
	return ret
}

func (cp ConsequencesPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	cm, cmok := model.(ConsequencesModel)
	if cmok {
		switch cm.OutputType {
		case JsonOutput:
			output := component.OutputDataLocation{
				Parameter:            "consequences",
				Format:               "json",
				LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.json", model.ModelName())},
				GeneratingModelName:  model.ModelName(),
				GeneratingPluginName: cp.Name(),
			}
			ret = append(ret, output)
			break
		case GpkgOutput:
			output := component.OutputDataLocation{
				Parameter:            "consequences",
				Format:               "geopackage",
				LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.gpkg", model.ModelName())},
				GeneratingModelName:  model.ModelName(),
				GeneratingPluginName: cp.Name(),
			}
			ret = append(ret, output)
			break
		case ShpOutput:
			output := component.OutputDataLocation{
				Parameter:            "consequences",
				Format:               "shapefile",
				LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.shp", model.ModelName())},
				GeneratingModelName:  model.ModelName(),
				GeneratingPluginName: cp.Name(),
			}
			ret = append(ret, output)
			break
		case SummaryOutput:
			output := component.OutputDataLocation{
				Parameter:            "consequences",
				Format:               "summary",
				LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
				GeneratingModelName:  model.ModelName(),
				GeneratingPluginName: cp.Name(),
			}
			ret = append(ret, output)
			break
		}

	}
	return ret
}

func (cp ConsequencesPlugin) Compute(model component.Model, options option.Options) error {
	//model.ModelLinkages[]
	hdfPath := ""
	terrainPath := ""
	links := model.ModelLinkages()
	linksAreGood := 0
	for _, i := range links.Links {
		param := i.OutputDataLocation.Parameter
		switch param {
		case "depth":
			link, linkok := i.OutputDataLocation.LinkInfo.(component.LocalCSVLink)
			if linkok {
				linksAreGood += 1
				hdfPath = link.Path
			}
			break
		case "elev":
			link, linkok := i.OutputDataLocation.LinkInfo.(component.LocalCSVLink)
			if linkok {
				linksAreGood += 1
				terrainPath = link.Path
			}
			break
		}
	}
	cm, cmok := model.(ConsequencesModel)
	if cmok {

		if linksAreGood == 2 {
			_, err := hazardproviders.Init(terrainPath)
			if err != nil {
				log.Printf("Error loading terrain:%s\n", err)
				return err
			}
			log.Print(terrainPath)
			log.Print(hdfPath)
			log.Print(cm)
			/*
				rdh := RasDepthHazard{
					structurePath: cm.StructurePath,
					tcr:           &tcr,
					outputPath:    options.OutputDestination + cm.ModelName() + ".gpkg",
					filePath:      hdfPath,
					terrainPath:   terrainPath,
					inputEpsg:     int(4326),
				}
				rdh.Run()
			*/
		}
	}
	return nil
}
