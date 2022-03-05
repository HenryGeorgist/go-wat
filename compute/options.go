package compute

type EventOptions interface{
	StartTime() string
	EndTime() string
}

type TimeWindow struct{
	StartTime string
	EndTime string
}

type StochasticEventOptions struct{
	TimeWindow TimeWindow
	RealizationNumber Int64
	LifecycleNumber Int64
	EventNumber Int64
	EventSeed Int64
	RealizationSeed Int64
}

func(s StochasticEventOptions) StartTime(){
	return s.TimeWindow.StartTime
}
func(s StochasticEventOptions) EndTime(){
	return s.TimeWindow.EndTime
}
type DeterministicEventOptions struct{
	TimeWindow TimeWindow
}
func(d DeterministicEventOptions) StartTime(){
	return s.TimeWindow.StartTime
}
func(d DeterministicEventOptions) EndTime(){
	return s.TimeWindow.EndTime
}
type ComputeOptions struct{
	EventOptions
	InputSource string
	OutputDestination string
}
