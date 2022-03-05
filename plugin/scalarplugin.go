package plugin

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type ScalarPlugin struct {
	Model component.Model
}
type ScalarModel struct {
	name         string
	plugin       *component.Computable
	links        component.ModelLinks
	DefaultValue float64
}

func (sm ScalarModel) ModelName() string {
	return sm.name
}
func (sm ScalarModel) Plugin() *component.Computable {
	return sm.plugin
}
func (sm ScalarModel) ModelLinkages() component.ModelLinks {
	return sm.links
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
	sm, smok := model.(ScalarModel)
	if smok {
		value := sm.DefaultValue
		stochastic, ok := options.EventOptions.(compute.StochasticEventOptions)
		if ok {
			//use a seed!
			r := rand.New(rand.NewSource(stochastic.EventSeed))
			value = r.NormFloat64()
		}
		//write it to the output destination in some agreed upon format?
		fmt.Println(value)
		return nil
	}
	return errors.New("could not cast the model to a scalar model")
	//not sure how we do linkages yet.

}
