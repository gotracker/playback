package song

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

type PeriodCalculatorIntf interface {
	GetPeriodGeneric(note.Note) period.Period
	PortaToNoteGeneric(period.Period, period.Delta, period.Period) (period.Period, error)
	PortaDownGeneric(period.Period, period.Delta) (period.Period, error)
	PortaUpGeneric(period.Period, period.Delta) (period.Period, error)
	AddDeltaGeneric(period.Period, period.Delta) (period.Period, error)
}

type PeriodCalculator[TPeriod period.Period] interface {
	PeriodCalculatorIntf

	GetPeriod(note.Note) TPeriod
	PortaToNote(TPeriod, period.Delta, TPeriod) (TPeriod, error)
	PortaDown(TPeriod, period.Delta) (TPeriod, error)
	PortaUp(TPeriod, period.Delta) (TPeriod, error)
	AddDelta(TPeriod, period.Delta) (TPeriod, error)

	GetSamplerAdd(TPeriod, float64) float64
	GetFrequency(TPeriod) frequency.Frequency
}
