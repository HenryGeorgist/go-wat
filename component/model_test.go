package component_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-wat/component"
	"go-wat/plugins"
)

func TestMarshalModelLinks(t *testing.T) {
	//create two scalar plugins
	spa := plugins.ScalarPlugin{}
	spb := plugins.ScalarPlugin{}

	//create two scalar models
	sma := plugins.ScalarModel{Name: "ValueA", DefaultValue: 2.0, ParentPluginName: spa.Name()}
	smb := plugins.ScalarModel{Name: "ValueB", DefaultValue: 2.0, ParentPluginName: spb.Name()}

	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	//model links
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
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}
func TestMarshalInputDataLocation(t *testing.T) {
	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	//model link
	aminputs := ap.InputLinks(am)

	bytes, err := json.Marshal(aminputs[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))

}
func TestMarshalOutputDataLocation(t *testing.T) {
	//create an add plugin
	ap := plugins.AddPlugin{}

	//create an add model
	am := plugins.AddModel{Name: "APlusB", ParentPluginName: ap.Name()}

	//model link
	amoutputs := ap.OutputLinks(am)

	bytes, err := json.Marshal(amoutputs[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))

}
