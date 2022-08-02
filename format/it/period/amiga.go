package period

import (
	"fmt"
	"math"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"

	"github.com/gotracker/playback/period"
)

// Amiga defines a sampler period that follows the Amiga-style approach of note
// definition. Useful in calculating resampling.
type Amiga struct {
	period.AmigaPeriod
	Coeff period.Frequency
}

// AddInteger truncates the current period to an integer and adds the delta integer in
// then returns the resulting period
func (p Amiga) AddInteger(delta int) Amiga {
	p.AmigaPeriod = period.AmigaPeriod(int(p.AmigaPeriod) + delta)
	return p
}

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) AddDelta(delta period.Delta, sign int) period.Period {
	d := period.ToPeriodDelta(delta) * period.PeriodDelta(sign)
	p.AmigaPeriod += period.AmigaPeriod(d)
	return p
}

// Compare returns:
//  -1 if the current period is higher frequency than the `rhs` period
//  0 if the current period is equal in frequency to the `rhs` period
//  1 if the current period is lower frequency than the `rhs` period
func (p Amiga) Compare(rhs period.Period) comparison.Spaceship {
	lf := p.GetFrequency()
	rf := rhs.GetFrequency()

	switch {
	case lf < rf:
		return comparison.SpaceshipRightGreater
	case lf > rf:
		return comparison.SpaceshipLeftGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Amiga) Lerp(t float64, rhs period.Period) period.Period {
	var right Amiga
	if r, ok := rhs.(Amiga); ok {
		right = r
	}

	p.AmigaPeriod = p.AmigaPeriod.Lerp(t, right.AmigaPeriod)
	return p
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Amiga) GetSamplerAdd(samplerSpeed period.Frequency) float64 {
	return float64(p.AmigaPeriod.GetFrequency(samplerSpeed) * p.Coeff)
}

// GetFrequency returns the frequency defined by the period
func (p Amiga) GetFrequency() period.Frequency {
	return p.AmigaPeriod.GetFrequency(BaseClock) * p.Coeff
}

func (p Amiga) String() string {
	return fmt.Sprintf("Amiga{ Period:%f, Coeff:%f }", float32(p.AmigaPeriod), p.Coeff)
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c2spd period.Frequency) Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	st := note.Semitone(finetunes / semitonesPerNote)
	ft := finetunes % semitonesPerNote

	k := st.Key()
	o := st.Octave()
	linFreq := math.Pow(2, float64(o)+(float64(ft)/semitonesPerOctave))

	p := Amiga{
		AmigaPeriod: period.AmigaPeriod(float64(semitonePeriodTable[k]) / linFreq),
		Coeff:       c2spd / MiddleCFrequency,
	}
	return p
}
