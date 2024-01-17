package period

import (
	"github.com/gotracker/playback/format/s3m/system"
	"github.com/gotracker/playback/period"
)

var S3MAmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System:    system.S3MSystem,
	MinPeriod: 64,
	MaxPeriod: 32767,
}

//MinMOD15Period = 452
//MaxMOD15Period = 3424

var MODAmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System:    system.S3MSystem,
	MinPeriod: 56,
	MaxPeriod: 13696,
}
