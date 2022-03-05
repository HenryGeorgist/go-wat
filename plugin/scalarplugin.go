package plugin

import (
	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type ScalarPlugin struct {
	Model component.Model
}

func (s ScalarPlugin) InputLinks() []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	return ret
}
func (s ScalarPlugin) OutputLinks() []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Parameter:       "float64",
		Format:          "scalar",
		LinkInfo:        "on disk?",
		GeneratingModel: &s.Model,
	}
	ret = append(ret, output)
	return ret
}
func (s ScalarPlugin) Compute(model component.Model, options compute.Options) error {
	//model.ModelLinkages[]
	//not sure how we do linkages yet.
	return nil
}
