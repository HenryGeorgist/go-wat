package plugin

import (
	"time"

	"github.com/HenryGeorgist/go-wat/compute"
)

type AnnualEventGenerator struct {
}

func (a AnnualEventGenerator) GenerateTimeWindows(t compute.TimeWindow) []compute.TimeWindow {
	timewindows := make([]compute.TimeWindow, 0)
	eventsToGenerate := (t.EndTime.Year() - t.StartTime.Year())
	year := t.StartTime.Year()
	for i := 0; i < eventsToGenerate; i++ {
		eventtwstart := time.Date(year, time.January, 1, 1, 1, 1, 1, time.Local)
		if i == 0 {
			eventtwstart = t.StartTime
		}
		eventtwend := time.Date(year, time.December, 31, 23, 59, 59, 1, time.Local)
		if i == eventsToGenerate-1 {
			eventtwend = t.EndTime
		}
		tw := compute.TimeWindow{
			StartTime: eventtwstart,
			EndTime:   eventtwend,
		}
		year++
		timewindows = append(timewindows, tw) //TODO: split on october first or january 1st
	}

	return timewindows
}
