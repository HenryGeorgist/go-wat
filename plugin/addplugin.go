package plugin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"go-wat/component"
	"go-wat/compute"
)

var output float64 = 0.0 // bad form

type AddPlugin struct {
}

type AddModel struct {
	Name             string               `json:"name"`
	ParentPluginName string               `json:"parent_plugin_name"`
	Links            component.ModelLinks `json:"-"`
	output           float64              `json:"-"`
}

func (am AddModel) ModelName() string {
	return am.Name
}

func (am AddModel) PluginName() string {
	return am.ParentPluginName
}

func (ap AddPlugin) MarshalJSON() ([]byte, error) {
	ret := "{\"plugin_name\":\"" + ap.Name() + "\"}"
	return []byte(ret), nil
}

func (am AddModel) ModelLinkages() component.ModelLinks {
	return am.Links
}

func (ap AddPlugin) Name() string {
	return "Add Plugin"
}

func (ap AddPlugin) InputLinks(model component.Model) []component.InputDataLocation {
	ret := make([]component.InputDataLocation, 0)
	valueA := component.InputDataLocation{
		Name:      "valueA",
		Parameter: "float64",
		Format:    "scalar",
	}
	valueB := component.InputDataLocation{
		Name:      "valueB",
		Parameter: "float64",
		Format:    "scalar",
	}
	ret = append(ret, valueA)
	ret = append(ret, valueB)
	return ret
}

func (ap AddPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
	ret := make([]component.OutputDataLocation, 0)
	output := component.OutputDataLocation{
		Parameter:            "float64",
		Format:               "scalar",
		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
		GeneratingModelName:  model.ModelName(),
		GeneratingPluginName: ap.Name(),
	}
	ret = append(ret, output)
	return ret
}

func (ap AddPlugin) Compute(model component.Model, options compute.Options) error {
	//model.ModelLinkages[]
	valueA := 0.0
	valueB := 0.0

	//inputs := a.InputLinks(model)
	links := model.ModelLinkages()
	link1 := true
	for _, i := range links.Links {
		lcsv, linkok := i.OutputDataLocation.LinkInfo.(component.LocalCSVLink)
		if linkok {
			inputdest := options.InputSource + lcsv.Path
			f, err := os.Open(inputdest)
			if err != nil {
				fmt.Println("could not find input link")
				return err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			scanner.Scan()
			s := scanner.Text()
			if link1 {
				valueA, err = strconv.ParseFloat(s, 64)
				if err != nil {
					fmt.Println("could not parse the first file")
				}
				link1 = false
			} else {
				valueB, err = strconv.ParseFloat(s, 64)
				if err != nil {
					fmt.Println("could not parse the second file")
				}
			}
		}

	}
	//add them together
	result := valueA + valueB
	output = result
	//write out the result.
	outputs := ap.OutputLinks(model)
	for _, o := range outputs {
		lcsv, _ := o.LinkInfo.(component.LocalCSVLink)
		outputdest := options.OutputDestination + lcsv.Path
		w, err := os.OpenFile(outputdest, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
		}
		defer w.Close()
		fmt.Fprint(w, result)
	}
	return nil
}

//implement the output variable interface.
func (ap AddPlugin) OutputVariables(model component.Model) []string {
	ret := make([]string, 1)
	ret[0] = ap.Name() + " output"
	return ret
}

func (ap AddPlugin) ComputeOutputVariables(selectedVariables []string, model component.Model) ([]float64, error) {
	ret := make([]float64, 0)
	if len(selectedVariables) > 0 {
		if selectedVariables[0] == ap.Name()+" output" {
			return append(ret, output), nil
		}
	}
	return ret, fmt.Errorf("no output variable found")
}
