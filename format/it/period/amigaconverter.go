package period

import (
	"math"

	"github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/period"
)

var AmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System:    system.ITSystem,
	MinPeriod: 1,
	MaxPeriod: math.MaxUint16,
}
