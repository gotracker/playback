package period

import (
	"github.com/gotracker/playback/format/xm/system"
	"github.com/gotracker/playback/period"
)

var AmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System: system.XMSystem,
}
