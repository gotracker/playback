package period

import (
	"github.com/gotracker/playback/period"
)

var AmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	BaseClock: S3MBaseClock,
}
