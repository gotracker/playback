package period

import (
	"github.com/gotracker/playback/format/s3m/system"
	"github.com/gotracker/playback/period"
)

var S3MAmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System:    system.S3MSystem,
	MinPeriod: 64,
	MaxPeriod: 32767,
	DeltaMult: 2,
}

//MinMOD15Period = 113 << 2
//MaxMOD15Period = 856 << 2

var MODAmigaConverter period.PeriodConverter[period.Amiga] = period.AmigaConverter{
	System:    system.S3MSystem,
	MinPeriod: 14 << 2,
	MaxPeriod: 3424 << 2,
	DeltaMult: 2,
}
