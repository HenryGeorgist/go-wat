package plugins

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"go-wat/component"
	"go-wat/option"
)

type ScalarPlugin struct {
}

type ScalarModel struct {
	Name             string               `json:"name"`
	ParentPluginName string               `json:"parent_plugin_name"`
	Links            component.ModelLinks `json:"-"`
	DefaultValue     float64
}

func (sm ScalarModel) ModelName() string {
	return sm.Name
}

func (sm ScalarModel) PluginName() string {
	return sm.ParentPluginName
}

func (sm ScalarModel) ModelLinkages() component.ModelLinks {
	return sm.Links
}

func (s ScalarPlugin) Name() string {
	return "Scalar Plugin"
}

func (s ScalarPlugin) MarshalJSON() ([]byte, error) {
	ret := "{\"plugin_name\":\"" + s.Name() + "\"}"
	return []byte(ret), nil
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
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: s.Name(),
	}
	ret = append(ret, output)
	return ret
}

func (s ScalarPlugin) Compute(model component.Model, options option.Options) error {
	sm, smok := model.(ScalarModel)
	if smok {
		value := sm.DefaultValue
		stochastic, ok := options.EventOptions.(option.StochasticEventOptions)
		if ok {
			//use a seed!
			r := rand.New(rand.NewSource(stochastic.EventSeeds[options.CurrentModelIndex()]))
			value = r.NormFloat64()
		}
		//write it to the output destination in some agreed upon format?
		outputs := s.OutputLinks(model)
		for _, o := range outputs {
			lcsv, _ := o.LinkInfo.(component.LocalCSVLink)
			outputdest := options.OutputDestination + lcsv.Path
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
