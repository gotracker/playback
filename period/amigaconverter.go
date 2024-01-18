package period

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

// AmigaConverter defines a sampler period that follows the AmigaConverter-style approach of note
// definition. Useful in calculating resampling.
type AmigaConverter struct {
	System          system.ClockableSystem
	MinPeriod       Amiga
	MaxPeriod       Amiga
	SlideTo0Allowed bool
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
	return c.System.GetBaseClock() / frequency.Frequency(p<<Amiga(c.System.GetOctaveShift()))
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (c AmigaConverter) GetSamplerAdd(p Amiga, instrumentRate, outputRate frequency.Frequency) float64 {
	if p.IsInvalid() {
		return 0
	}
	//return float64(c.System.GetCommonPeriod()) / float64(p<<Amiga(c.System.GetOctaveShift())) * float64(instrumentRate/outputRate)
	// this produces a slightly more correct value than above
	return float64(c.GetFrequency(p) * (instrumentRate / (c.System.GetCommonRate() * outputRate)))
}

func (c AmigaConverter) GetPeriod(n note.Note) Amiga {
	switch n.Type() {
	case note.SpecialTypeEmpty, note.SpecialTypeRelease, note.SpecialTypeStop, note.SpecialTypeStopOrRelease:
		return Amiga(0)
	case note.SpecialTypeNormal:
		semi := note.Semitone(n.(note.Normal))
		key := semi.Key()
		octave := int(semi.Octave())

		keyPeriod, valid := c.System.GetSemitonePeriod(key)
		if !valid {
			return Amiga(0)
		}

		p := Amiga(keyPeriod) >> octave

		return p.Clamp(c.MinPeriod, c.MaxPeriod)
	case note.SpecialTypeInvalid:
		fallthrough
	default:
		panic("unsupported note type")
	}
}

func (c AmigaConverter) PortaToNote(p Amiga, delta Delta, target Amiga) (Amiga, error) {
	return p.PortaTo(delta, target, c.MinPeriod, c.MaxPeriod), nil
}

func (c AmigaConverter) PortaDown(p Amiga, delta Delta) (Amiga, error) {
	return p.PortaDown(delta, c.MinPeriod, c.MaxPeriod, c.SlideTo0Allowed), nil
}

func (c AmigaConverter) PortaUp(p Amiga, delta Delta) (Amiga, error) {
	return p.PortaUp(delta, c.MinPeriod, c.MaxPeriod, c.SlideTo0Allowed), nil
}

func (c AmigaConverter) AddDelta(p Amiga, delta Delta) (Amiga, error) {
	return p.Add(delta, c.MinPeriod, c.MaxPeriod, c.SlideTo0Allowed), nil
}

func (c AmigaConverter) clamp(p Amiga) Amiga {
	return min(max(p, c.MinPeriod), c.MaxPeriod)
}
