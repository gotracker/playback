package period

import (
	"github.com/gotracker/playback/period"
)

var LinearConverter period.PeriodConverter[period.Linear] = period.LinearConverter{
	BaseClock:      XMBaseClock,
	BaseFinetune:   C4SlideFines,
	FinesPerOctave: SlideFinesPerOctave,
}
