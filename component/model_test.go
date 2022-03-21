package component_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-wat/component"
	"go-wat/plugin"
)

func TestMarshalModelLinks(t *testing.T) {
	//create two scalar plugins
	spa := plugin.ScalarPlugin{}
	spb := plugin.ScalarPlugin{}
	//create two scalar models
	sma := plugin.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugin.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}
	//create an add plugin
	ap := plugin.AddPlugin{}
	//create an add model
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
	//model link
	aminputs := ap.InputLinks(am)
	smaoutput := spa.OutputLinks(sma)
	smboutput := spa.OutputLinks(smb)

	modelLinks := make([]component.Link, 2)
	modelLinks[0] = component.Link{InputDataLocation: aminputs[0], OutputDataLocation: smaoutput[0]}
	modelLinks[1] = component.Link{InputDataLocation: aminputs[1], OutputDataLocation: smboutput[0]}
	ml := component.ModelLinks{Links: modelLinks}
	am.Links = ml
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
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
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
	am := plugin.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}
	//model link
	amoutputs := ap.OutputLinks(am)

	bytes, err := json.Marshal(amoutputs[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(bytes))
}
