package period

import (
	"github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/period"
)

var LinearConverter period.PeriodConverter[period.Linear] = period.LinearConverter{
	System: system.ITSystem,
}
