package component

import "github.com/HenryGeorgist/go-wat/compute"

type Computable interface {
	Compute(model Model, options compute.ComputeOptions) error
	InputLinks() []InputDataLocation
	OutputLinks() []OutputDataLocation
}

type OutputReporter interface {
	OutputVariables() []string
	ComputeOutputVariables(selectedVariables []string) []string //TODO fix this
}

type EventGenerator interface {
	GenerateTimeWindows() []compute.TimeWindow
}
