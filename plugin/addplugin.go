package plugin

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

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
	valueA := 0.0
	valueB := 0.0

	inputs := a.InputLinks(model)
	links := model.ModelLinkages()
	link1 := true
	for _, i := range inputs {
		input, ok := links.Links[i]
		if !ok {
			fmt.Println("could not find input link")
			return errors.New("couldnt find link")
		}
		inputdest := options.InputSource + input.LinkInfo
		f, err := os.Open(inputdest)
		defer f.Close()
		if err != nil {
			fmt.Println("could not find input link")
			return err
		}
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
	//add them together
	result := valueA + valueB
	//write out the result.
	outputs := a.OutputLinks(model)
	for _, o := range outputs {
		outputdest := options.OutputDestination + o.LinkInfo
		w, err := os.OpenFile(outputdest, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
		}
		defer w.Close()
		fmt.Fprint(w, result)
	}
	return errors.New("under construction")
}
