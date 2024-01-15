package period

import (
	"errors"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	System system.ClockableSystem
}

var _ PeriodConverter[Amiga] = (*AmigaConverter)(nil)

// GetFrequency returns the frequency defined by the period
func (c AmigaConverter) GetFrequency(p Amiga) frequency.Frequency {
	if p.IsInvalid() {
		return 0
	}
	return c.System.GetBaseClock() / frequency.Frequency(p)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c AmigaConverter) GetSamplerAdd(p Amiga, samplerSpeed float64) float64 {
	if p.IsInvalid() {
		return 0
	}
	return samplerSpeed / (float64(p) * float64(c.System.GetCommonRate()))
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

		return Amiga(float64(keyPeriod) / float64(uint32(1)<<octave))
	case note.SpecialTypeInvalid:
		fallthrough
	default:
		panic("unsupported note type")
	}
}

func (c AmigaConverter) GetPeriodGeneric(n note.Note) Period {
	return c.GetPeriod(n)
}

func (c AmigaConverter) PortaToNote(lhs Amiga, delta Delta, rhs Amiga) (Amiga, error) {
	return PortaTo(lhs, int(delta), rhs), nil
}

func (c AmigaConverter) PortaToNoteGeneric(p Period, delta Delta, target Period) (Period, error) {
	lhs, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	rhs, ok := target.(Amiga)
	if !ok {
		return p, errors.New("invalid target period type conversion")
	}

	return PortaTo(lhs, int(delta), rhs), nil
}

func (c AmigaConverter) PortaDown(p Amiga, delta Delta) (Amiga, error) {
	return p.PortaDown(int(delta)), nil
}

func (c AmigaConverter) PortaDownGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return cur.PortaDown(int(delta)), nil
}

func (c AmigaConverter) PortaUp(p Amiga, delta Delta) (Amiga, error) {
	return p.PortaUp(int(delta)), nil
}

func (c AmigaConverter) PortaUpGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return cur.PortaUp(int(delta)), nil
}
