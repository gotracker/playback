package period

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

type PeriodConverter[TPeriod Period] interface {
	GetSystem() system.System

	GetPeriod(note.Note) TPeriod
	PortaToNote(TPeriod, Delta, TPeriod) (TPeriod, error)
	PortaDown(TPeriod, Delta) (TPeriod, error)
	PortaUp(TPeriod, Delta) (TPeriod, error)
	AddDelta(TPeriod, Delta) (TPeriod, error)

	GetSamplerAdd(TPeriod, frequency.Frequency, frequency.Frequency) float64
	GetFrequency(TPeriod) frequency.Frequency
}
