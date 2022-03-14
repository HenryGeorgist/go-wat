package compute

import "time"

//EventOptions defines the interface for event options, and requires a TimeWindow for the event
type EventOptions interface {
	TimeWindow() TimeWindow
	//UpdateTimeWindow(t TimeWindow) EventOptions
	CurrentModelIndex() int
	IncrementModelIndex() EventOptions
	ResetModelIndex() EventOptions
}

type TimeWindow struct {
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}

//StochasticEventOptions implements EventOptions and adds information about realizations, lifecycles, and events as well as random seeds
type StochasticEventOptions struct {
	EventTimeWindow   TimeWindow `json:"timewindow"`
	RealizationNumber int        `json:"realizationnumber"`
	LifecycleNumber   int        `json:"lifecyclenumber"`
	EventNumber       int        `json:"eventnumber"`
	EventSeeds        []int64    `json:"eventseeds"`
	RealizationSeeds  []int64    `json:"realizationseeds"`
	CurrentModel      int        `json:"CurrentModel"`
}

func (s StochasticEventOptions) TimeWindow() TimeWindow {
	return s.EventTimeWindow
}

//DeterministicEventOptions implements the EventOptions interface for a deterministic compute
type DeterministicEventOptions struct {
	EventTimeWindow TimeWindow `json:"timewindow"`
}

func (d DeterministicEventOptions) TimeWindow() TimeWindow {
	return d.EventTimeWindow
}

//Options composes EventOptions with information about the location of inputdata and where output data should be stored
type Options struct {
	EventOptions      `json:"eventoptions"`
	InputSource       string              `json:"inputsource"`
	OutputDestination string              `json:"outputdestination"`
	OutputVariables   map[string][]string `json:"output_variables"`
}

func (d DeterministicEventOptions) CurrentModelIndex() int {
	return 0
}
func (d DeterministicEventOptions) IncrementModelIndex() EventOptions {
	return d
}
func (d DeterministicEventOptions) ResetModelIndex() EventOptions {
	return d
}
func (s StochasticEventOptions) CurrentModelIndex() int {
	return s.CurrentModel
}
func (s StochasticEventOptions) IncrementModelIndex() EventOptions {
	s.CurrentModel += 1
	return s
}
func (s StochasticEventOptions) ResetModelIndex() EventOptions {
	s.CurrentModel = 0
	return s
}
