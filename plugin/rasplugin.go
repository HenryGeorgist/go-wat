package plugin

// import (
// 	"fmt"
// 	"time"

// 	"go-wat/component"
// 	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
// )

// type RasPlugin struct {
// }

// // Name of the model (ras - bc-1)
// // ParentPluginName (hyrdoscalaar plugin)
// // Flows shape set
// // Flow frequency LPIII (or other BootstrappableDistribution)
// // Links are output links (what can I produce)
// type RasModel struct {
// 	Name             string                                `json:"name"`
// 	ParentPluginName string                                `json:"parent_plugin_name"`
// 	Flows            []float64                             `json:"flows"`
// 	TimeStep         time.Duration                         `json:"timestep"`
// 	FlowFrequency    statistics.BootstrappableDistribution `json:"flow_frequency"`
// 	Links            component.ModelLinks                  `json:"-"`
// }

// //model implementation
// func (hsm RasModel) ModelName() string {
// 	return hsm.Name
// }

// func (hsm RasModel) PluginName() string {
// 	return hsm.ParentPluginName
// }

// // Todo: Write this first
// func (hsm RasModel) ModelLinkages() component.ModelLinks {
// 	return hsm.Links
// }

// // Get input list of BC's from a plan from HDF...
// func (hsp RasPlugin) InputLinks(model component.Model) []component.InputDataLocation {
// 	// no links needed here, this serves as a generator in this context
// 	ret := make([]component.InputDataLocation, 0)
// 	return ret
// }

// // List output file (p*.hdf)
// func (hsp RasPlugin) OutputLinks(model component.Model) []component.OutputDataLocation {
// 	ret := make([]component.OutputDataLocation, 0)
// 	output := component.OutputDataLocation{
// 		Name:                 "Hydrograph",
// 		Parameter:            "Flow",
// 		Format:               "Array",
// 		LinkInfo:             component.LocalCSVLink{Path: fmt.Sprintf("/%v.csv", model.ModelName())},
// 		GeneratingModelName:  model.ModelName(),
// 		GeneratingPluginName: hsp.Name(),
// 	}
// 	ret = append(ret, output)
// 	return ret
// }
