package compute

type EventOptions interface {
	StartTime() string
	EndTime() string
}

type TimeWindow struct {
	StartTime string
	EndTime   string
}

type StochasticEventOptions struct {
	TimeWindow        TimeWindow
	RealizationNumber int
	LifecycleNumber   int
	EventNumber       int
	EventSeed         int64
	RealizationSeed   int64
}

func (s StochasticEventOptions) StartTime() string {
	return s.TimeWindow.StartTime
}
func (s StochasticEventOptions) EndTime() string {
	return s.TimeWindow.EndTime
}

type DeterministicEventOptions struct {
	TimeWindow TimeWindow
}

func (d DeterministicEventOptions) StartTime() string {
	return d.TimeWindow.StartTime
}
func (d DeterministicEventOptions) EndTime() string {
	return d.TimeWindow.EndTime
}

type ComputeOptions struct {
	EventOptions
	InputSource       string
	OutputDestination string
}
