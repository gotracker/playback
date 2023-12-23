package period

import (
	"github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/period"
)

var AmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System: system.ITSystem,
}
