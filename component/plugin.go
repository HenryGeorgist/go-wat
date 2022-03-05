package component

import "github.com/HenryGeorgist/go-wat/compute"

//Computable defines the interface to perform a compute and facilitate links
type Computable interface {
	Compute(model Model, options compute.ComputeOptions) error
	InputLinks() []InputDataLocation
	OutputLinks() []OutputDataLocation
}

//OutputReporter provides the interface to compute output varibles
type OutputReporter interface {
	OutputVariables() []string
	ComputeOutputVariables(selectedVariables []string) []string //TODO fix this
}

//EventGenerator is an additional interface to generate timewindows for a lifecycle
type EventGenerator interface {
	GenerateTimeWindows(timewindow compute.TimeWindow) []compute.TimeWindow
}
