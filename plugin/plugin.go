package plugin

import "github.com/HenryGeorgist/go-wat/compute"

type Plugin interface {
	Compute(options compute.ComputeOptions) error
	InputLinks() []string
	OutputLinks() []string
}

type OutputReporter interface {
	OutputVariables() []string
	ComputeOutputVariables(selectedVariables []string) []string //TODO fix this
}

type EventGenerator interface {
	GenerateTimeWindows() []compute.TimeWindow
}
