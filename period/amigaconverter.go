package period

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	System    system.ClockableSystem
	MinPeriod Amiga
	MaxPeriod Amiga
}

var _ PeriodConverter[Amiga] = (*AmigaConverter)(nil)

func (c AmigaConverter) GetSystem() system.System {
	return c.System
}

// GetFrequency returns the frequency defined by the period
func (c AmigaConverter) GetFrequency(p Amiga) frequency.Frequency {
	if p.IsInvalid() {
		return 0
	}
	return c.System.GetBaseClock() / frequency.Frequency(p)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c AmigaConverter) GetSamplerAdd(p Amiga, instrumentRate, outputRate frequency.Frequency) float64 {
	if p.IsInvalid() {
		return 0
	}
	return c.System.GetCommonPeriod() * float64(instrumentRate/outputRate) / float64(p)
}

func (c AmigaConverter) GetPeriod(n note.Note) Amiga {
	switch n.Type() {
	case note.SpecialTypeEmpty, note.SpecialTypeRelease, note.SpecialTypeStop, note.SpecialTypeStopOrRelease:
		return Amiga(0)
	case note.SpecialTypeNormal:
		semi := note.Semitone(n.(note.Normal))
		octave := uint32(semi.Octave())

		keyPeriod, valid := c.System.GetSemitonePeriod(semi.Key())
		if !valid {
			return Amiga(0)
		}

		return c.clamp(Amiga(float64(keyPeriod) / float64(uint32(1)<<octave)))
	case note.SpecialTypeInvalid:
		fallthrough
	default:
		panic("unsupported note type")
	}
}

func (c AmigaConverter) PortaToNote(p Amiga, delta Delta, target Amiga) (Amiga, error) {
	return c.clamp(p.PortaTo(delta, target)), nil
}

func (c AmigaConverter) PortaDown(p Amiga, delta Delta) (Amiga, error) {
	return c.clamp(p.PortaDown(delta)), nil
}

func (c AmigaConverter) PortaUp(p Amiga, delta Delta) (Amiga, error) {
	return c.clamp(p.PortaUp(delta)), nil
}

func (c AmigaConverter) AddDelta(p Amiga, delta Delta) (Amiga, error) {
	return c.clamp(p.Add(delta)), nil
}

func (c AmigaConverter) clamp(p Amiga) Amiga {
	return min(max(p, c.MinPeriod), c.MaxPeriod)
}
