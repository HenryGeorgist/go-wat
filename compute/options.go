package compute

import "time"

//EventOptions defines the interface for event options, and requires a TimeWindow for the event
type EventOptions interface {
	TimeWindow() TimeWindow
	//UpdateTimeWindow(t TimeWindow) EventOptions
}

type TimeWindow struct {
	StartTime time.Time
	EndTime   time.Time
}

//StochasticEventOptions implements EventOptions and adds information about realizations, lifecycles, and events as well as random seeds
type StochasticEventOptions struct {
	timeWindow        TimeWindow
	RealizationNumber int
	LifecycleNumber   int
	EventNumber       int
	EventSeed         int64
	RealizationSeed   int64
}

func (s *StochasticEventOptions) UpdateTimeWindow(t TimeWindow) EventOptions {
	s.timeWindow = t
	return s
}
func (s StochasticEventOptions) TimeWindow() TimeWindow {
	return s.timeWindow
}

//DeterministicEventOptions implements the EventOptions interface for a deterministic compute
type DeterministicEventOptions struct {
	timeWindow TimeWindow
}

func (d *DeterministicEventOptions) UpdateTimeWindow(t TimeWindow) EventOptions {
	d.timeWindow = t
	return d
}
func (d DeterministicEventOptions) TimeWindow() TimeWindow {
	return d.timeWindow
}

//Options composes EventOptions with information about the location of inputdata and where output data should be stored
type Options struct {
	EventOptions
	InputSource       string
	OutputDestination string
}
