package period

import (
	"github.com/gotracker/playback/note"
)

type PeriodConverter[TPeriod Period] interface {
	GetPeriodGeneric(note.Note) Period
	PortaToNoteGeneric(Period, Delta, Period) (Period, error)
	PortaDownGeneric(Period, Delta) (Period, error)
	PortaUpGeneric(Period, Delta) (Period, error)

	GetPeriod(note.Note) TPeriod
	PortaToNote(TPeriod, Delta, TPeriod) (TPeriod, error)
	PortaDown(TPeriod, Delta) (TPeriod, error)
	PortaUp(TPeriod, Delta) (TPeriod, error)
	GetSamplerAdd(TPeriod, float64) float64
	GetFrequency(TPeriod) Frequency
}
