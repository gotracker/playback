package period

import (
	"github.com/gotracker/playback/format/s3m/system"
	"github.com/gotracker/playback/period"
)

var AmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System: system.S3MSystem,
}
