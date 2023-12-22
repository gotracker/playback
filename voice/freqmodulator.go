package voice

import (
	"github.com/gotracker/playback/period"
)

// FreqModulator is the instrument frequency control interface
type FreqModulator[TPeriod period.Period] interface {
	SetPeriod(period TPeriod)
	GetPeriod() TPeriod
	SetPeriodDelta(delta period.Delta)
	GetPeriodDelta() period.Delta
	GetFinalPeriod() TPeriod
}
