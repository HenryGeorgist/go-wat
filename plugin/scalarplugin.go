package plugin

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/compute"
)

type ScalarPlugin struct {
}
type ScalarModel struct {
	Name         string               `json:"name"`
	ParentPlugin component.Computable `json:"parent_plugin"`
	Links        component.ModelLinks `json:"-"`
	DefaultValue float64
}

func (sm ScalarModel) ModelName() string {
	return sm.Name
}
func (sm ScalarModel) Plugin() component.Computable {
	return sm.ParentPlugin
}
func (sm ScalarModel) ModelLinkages() component.ModelLinks {
	return sm.Links
}
func (s ScalarPlugin) Name() string {
	return "Scalar Plugin"
}
func (s ScalarPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	return ret
}
func (s ScalarPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Parameter:            "float64",
		Format:               "scalar",
		LinkInfo:             fmt.Sprintf("/%v.csv", model.ModelName()),
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: s.Name(),
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
		outputs := s.OutputLinks(model)
		for _, o := range outputs {
			outputdest := options.OutputDestination + o.LinkInfo
			w, err := os.OpenFile(outputdest, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
			if err != nil {
				fmt.Println(err)
			}
			defer w.Close()
			fmt.Fprint(w, value)
		}
		fmt.Println(value)
		return nil
	}
	return errors.New("could not cast the model to a scalar model")
	//not sure how we do linkages yet.

}
