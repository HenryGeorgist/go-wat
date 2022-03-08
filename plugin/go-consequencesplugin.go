package plugin

import (
	"fmt"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type GoConsequencesPlugin struct {
}
type GoConsequencesModel struct {
	Name                   string               `json:"name"`
	ParentPluginName       string               `json:"parent_plugin_name"`
	StructureInventoryPath string               `json:"structure_inventory_path"`
	StructureInventoryType string               `json:"structure_inventory_type"`
	OutputType             string               `json:"output_type"`
	Links                  component.ModelLinks `json:"-"`
}

func (gcm GoConsequencesModel) ModelName() string {
	return gcm.Name
}
func (gcm GoConsequencesModel) PluginName() string {
	return gcm.ParentPluginName
}
func (gcp GoConsequencesPlugin) MarshalJSON() ([]byte, error) {
	ret := "{\"plugin_name\":\"" + gcp.Name() + "\"}"
	return []byte(ret), nil
}
func (gcm GoConsequencesModel) ModelLinkages() component.ModelLinks {
	return gcm.Links
}
func (gcp GoConsequencesPlugin) Name() string {
	return "go-consequences Plugin"
}
func (gcp GoConsequencesPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	terrain := component.InputDataLocation{
		Name:      "Terrain",
		Parameter: "Elevation",
		Format:    "geotif",
	}
	rasHDFOutput := component.InputDataLocation{
		Name:      "HEC-RAS HDF output",
		Parameter: "Depth",
		Format:    "HDF",
	}
	ret = append(ret, terrain)
	ret = append(ret, rasHDFOutput)
	return ret
}
func (gcp GoConsequencesPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Name:                 "Study EAD",
		Parameter:            "float64",
		Format:               "scalar",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: gcp.Name(),
	}
	ret = append(ret, output)
	return ret
}
func (gcp GoConsequencesPlugin) Compute(model component.Model, options compute.Options) error {
	//get terrain and hdf from HEC-RAS
	//load inventory
	//set output location based on compute options.
	return nil
}
