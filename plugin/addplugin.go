package plugin

import (
	"errors"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type AddPlugin struct {
	Inputs  []component.InputDataLocation
	Outputs []component.OutputDataLocation
}

func (a AddPlugin) InputLinks() []component.InputDataLocation {
	return a.Inputs
}
func (a AddPlugin) OutputLinks() []component.OutputDataLocation {
	return a.Outputs
}
func (a AddPlugin) Compute(model component.Model, options compute.Options) error {
	//model.ModelLinkages[]
	//not sure how we do linkages yet.
	return errors.New("under construction")
}
