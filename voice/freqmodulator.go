package voice

import (
	"github.com/gotracker/playback/period"
)

// FreqModulator is the instrument frequency control interface
type FreqModulator interface {
	SetPeriod(period period.Period)
	GetPeriod() period.Period
	SetPeriodDelta(delta period.Delta)
	GetPeriodDelta() period.Delta
	GetFinalPeriod() period.Period
}
