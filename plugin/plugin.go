package plugin

type Plugin interface{
	Compute(options compute.ComputeOptions) error
	InputLinks() []string
	OutputLinks() []string
}

type OutputReporter interface{
	OutputVariables() []string
	ComputeOutputVariables(selectedVariables []string)[]interface //TODO fix this
}

type EventGenerator interface{
	GenerateTimeWindows() []compute.TimeWindow
}