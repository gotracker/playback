package period

import (
	"math"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

// LinearConverter defines a sampler period that follows the Linear-style approach of note
// definition. Useful in calculating resampling.
type LinearConverter struct {
	System system.ClockableSystem
}

var _ PeriodConverter[Linear] = (*LinearConverter)(nil)

func (c LinearConverter) GetSystem() system.System {
	return c.System
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c LinearConverter) GetSamplerAdd(p Linear, instrumentRate, outputRate frequency.Frequency) float64 {
	return float64(c.GetFrequency(p) * instrumentRate / outputRate)
}

// GetFrequency returns the frequency defined by the period
func (c LinearConverter) GetFrequency(p Linear) frequency.Frequency {
	if p.Finetune == 0 {
		return 0
	}
	pft := float64(p.Finetune-c.System.GetBaseFinetunes()) / float64(c.System.GetFinetunesPerOctave())
	f := frequency.Frequency(math.Pow(2.0, pft))
	return f
}

func (c LinearConverter) GetPeriod(n note.Note) Linear {
	switch n.Type() {
	case note.SpecialTypeEmpty, note.SpecialTypeRelease, note.SpecialTypeStop, note.SpecialTypeStopOrRelease:
		return Linear{}
	case note.SpecialTypeNormal:
		st := note.Semitone(n.(note.Normal))
		return Linear{
			Finetune: note.Finetune(st) * c.System.GetFinetunesPerSemitone(),
		}
	case note.SpecialTypeInvalid:
		fallthrough
	default:
		panic("unsupported note type")
	}
}

func (c LinearConverter) PortaToNote(p Linear, delta Delta, target Linear) (Linear, error) {
	return p.PortaTo(int(delta), target), nil
}

func (c LinearConverter) PortaDown(p Linear, delta Delta) (Linear, error) {
	return p.PortaDown(int(delta)), nil
}

func (c LinearConverter) PortaUp(p Linear, delta Delta) (Linear, error) {
	return p.PortaUp(int(delta)), nil
}

func (c LinearConverter) AddDelta(p Linear, delta Delta) (Linear, error) {
	return p.Add(delta), nil
}
