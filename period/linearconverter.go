package period

import (
	"errors"
	"math"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
	"github.com/heucuva/comparison"
)

// LinearConverter defines a sampler period that follows the Linear-style approach of note
// definition. Useful in calculating resampling.
type LinearConverter struct {
	System system.ClockableSystem
}

var _ PeriodConverter[Linear] = (*LinearConverter)(nil)

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c LinearConverter) GetSamplerAdd(p Linear, samplerSpeed float64) float64 {
	return float64(c.GetFrequency(p)) * samplerSpeed / float64(c.System.GetBaseClock())
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

func (c LinearConverter) GetPeriodGeneric(n note.Note) Period {
	return c.GetPeriod(n)
}

func (c LinearConverter) PortaToNote(p Linear, delta Delta, target Linear) (Linear, error) {
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
	return p, err
}

func (c LinearConverter) PortaToNoteGeneric(p Period, delta Delta, target Period) (Period, error) {
	lhs, ok := p.(Linear)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	rhs, ok := target.(Linear)
	if !ok {
		return p, errors.New("invalid period target type conversion")
	}

	return c.PortaToNote(lhs, delta, rhs)
}

func (c LinearConverter) PortaDown(p Linear, delta Delta) (Linear, error) {
	return p.PortaDown(int(delta)), nil
}

func (c LinearConverter) PortaDownGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Linear)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return cur.PortaDown(int(delta)), nil
}

func (c LinearConverter) PortaUp(p Linear, delta Delta) (Linear, error) {
	return p.PortaUp(int(delta)), nil
}

func (c LinearConverter) PortaUpGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Linear)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return cur.PortaUp(int(delta)), nil
}

func (c LinearConverter) AddDelta(p Linear, delta Delta) (Linear, error) {
	return p.Add(delta), nil
}

func (c LinearConverter) AddDeltaGeneric(p Period, delta Delta) (Period, error) {
	cur, ok := p.(Linear)
	if !ok {
		return p, errors.New("invalid period type conversion")
	}

	return c.AddDelta(cur, delta)
}
