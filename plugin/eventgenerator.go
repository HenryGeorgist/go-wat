package plugin

import (
	"github.com/HenryGeorgist/go-wat/compute"
)

type AnnualEventGenerator struct {
}

func (a AnnualEventGenerator) GenerateTimeWindows(t compute.TimeWindow) []compute.TimeWindow {
	timewindows := make([]compute.TimeWindow, 0)
	timewindows = append(timewindows, t) //TODO: split on october first or january 1st
	return timewindows
}
