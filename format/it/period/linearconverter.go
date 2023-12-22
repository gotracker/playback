package period

import (
	"github.com/gotracker/playback/period"
)

var LinearConverter period.PeriodConverter[period.Linear] = period.LinearConverter{
	BaseClock:      ITBaseClock,
	BaseFinetune:   C5SlideFines,
	FinesPerOctave: SlideFinesPerOctave,
}
