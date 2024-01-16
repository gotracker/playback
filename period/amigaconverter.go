package period

import (
	"errors"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
	"github.com/heucuva/comparison"
)

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	System    system.ClockableSystem
	MinPeriod Amiga
	MaxPeriod Amiga
	DeltaMult Delta
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

		return min(max(Amiga(float64(keyPeriod)/float64(uint32(1)<<octave)), c.MinPeriod), c.MaxPeriod)
	case note.SpecialTypeInvalid:
		fallthrough
	default:
		panic("unsupported note type")
	}
}

func (c AmigaConverter) GetPeriodGeneric(n note.Note) Period {
	return c.GetPeriod(n)
}

func (c AmigaConverter) PortaToNote(p Amiga, delta Delta, target Amiga) (Amiga, error) {
	var err error
	switch p.Compare(target) {
	case comparison.SpaceshipRightGreater:
		p, err = c.PortaUp(p, delta)
		if ComparePeriods(p, target) == comparison.SpaceshipLeftGreater {
			p = target
		}
	case comparison.SpaceshipLeftGreater:
		p, err = c.PortaDown(p, delta)
		if ComparePeriods(p, target) == comparison.SpaceshipRightGreater {
			p = target
		}
	}
	return min(max(p, c.MinPeriod), c.MaxPeriod), err
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

	return c.PortaToNote(lhs, delta, rhs)
}

func (c AmigaConverter) PortaDown(p Amiga, delta Delta) (Amiga, error) {
	return min(max(p.PortaDown(int(delta*c.DeltaMult)), c.MinPeriod), c.MaxPeriod), nil
}

func (c AmigaConverter) PortaDownGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return c.PortaDown(cur, delta)
}

func (c AmigaConverter) PortaUp(p Amiga, delta Delta) (Amiga, error) {
	return min(max(p.PortaUp(int(delta*c.DeltaMult)), c.MinPeriod), c.MaxPeriod), nil
}

func (c AmigaConverter) PortaUpGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return c.PortaUp(cur, delta)
}

func (c AmigaConverter) AddDelta(p Amiga, delta Delta) (Amiga, error) {
	return min(max(p.Add(delta*c.DeltaMult), c.MinPeriod), c.MaxPeriod), nil
}

func (c AmigaConverter) AddDeltaGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Amiga)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return c.AddDelta(cur, delta)
}
