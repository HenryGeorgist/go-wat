package plugin

import (
	"errors"
	"fmt"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type AddPlugin struct {
}
type AddModel struct {
	name   string
	plugin *component.Computable
	links  component.ModelLinks
}

func (am AddModel) ModelName() string {
	return am.name
}
func (sm AddModel) Plugin() *component.Computable {
	return sm.plugin
}
func (am AddModel) ModelLinkages() component.ModelLinks {
	return am.links
}
func (a AddPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	valueA := component.InputDataLocation{
		Parameter: "float64",
		Format:    "scalar",
	}
	valueB := component.InputDataLocation{
		Parameter: "float64",
		Format:    "scalar",
	}
	ret = append(ret, valueA)
	ret = append(ret, valueB)
	return ret
}
func (a AddPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Parameter:       "float64",
		Format:          "scalar",
		LinkInfo:        fmt.Sprintf("/%v.csv", model.ModelName()),
		GeneratingModel: &model,
	}
	ret = append(ret, output)
	return ret
}
func (a AddPlugin) Compute(model component.Model, options compute.Options) error {
	//model.ModelLinkages[]
	//not sure how we do linkages yet.
	//get the values from the links
	//add them together
	//write out the result.
	return errors.New("under construction")
}
