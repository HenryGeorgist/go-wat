package component_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/HenryGeorgist/go-wat/component"
	"github.com/HenryGeorgist/go-wat/plugin"
)

func TestMarshalModelLinks(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPlugin: spa}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPlugin: spb}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPlugin: ap}
	//model link
	aminputs := ap.InputLinks(am)
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)

	modelLinks := make(map[component.InputDataLocation]component.OutputDataLocation)
	modelLinks[aminputs[0]] = smaoutput[0]
	modelLinks[aminputs[1]] = smboutput[0]
	ml := component.ModelLinks{Links: modelLinks}
	am.Links = ml
	for i, o := range ml.Links {
		bytesi, err := json.Marshal(i)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bytesi))
		byteso, err := json.Marshal(o)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(byteso))
	}
	bytes, err := json.Marshal(ml)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(bytes))
}
func TestMarshalInputDataLocation(t *testing.T) {
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPlugin: ap}
	//model link
	aminputs := ap.InputLinks(am)

	bytes, err := json.Marshal(aminputs[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(bytes))
}
func TestMarshalOutputDataLocation(t *testing.T) {
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPlugin: ap}
	//model link
	amoutputs := ap.OutputLinks(am)

	bytes, err := json.Marshal(amoutputs[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(bytes))
}
