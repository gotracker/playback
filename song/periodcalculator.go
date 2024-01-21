package song

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/system"
)

type PeriodCalculatorIntf interface {
	GetSystem() system.System
}

type PeriodCalculator[TPeriod period.Period] interface {
	PeriodCalculatorIntf

	GetPeriod(note.Note) TPeriod
	PortaToNote(TPeriod, period.Delta, TPeriod) (TPeriod, error)
	PortaDown(TPeriod, period.Delta) (TPeriod, error)
	PortaUp(TPeriod, period.Delta) (TPeriod, error)
	AddDelta(TPeriod, period.Delta) (TPeriod, error)

	GetSamplerAdd(TPeriod, frequency.Frequency, frequency.Frequency) float64
	GetFrequency(TPeriod) frequency.Frequency
}
