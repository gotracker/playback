package period

import (
	"fmt"
	"math"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"

	"github.com/gotracker/playback/period"
)

// Linear is a linear period, based on semitone and finetune values
type Linear struct {
	Finetune note.Finetune
	C2Spd    period.Frequency
}

// Add adds the current period to a delta value then returns the resulting period
func (p Linear) AddDelta(delta period.Delta) period.Period {
	// 0 means "not playing", so keep it that way
	if p.Finetune > 0 {
		d := period.ToPeriodDelta(delta)
		p.Finetune += note.Finetune(d)
		if p.Finetune < 1 {
			p.Finetune = 1
		}
	}
	return p
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Linear) Compare(rhs period.Period) comparison.Spaceship {
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
func (p Linear) Lerp(t float64, rhs period.Period) period.Period {
	right := ToLinearPeriod(rhs)

	lnft := float64(p.Finetune)
	rnft := float64(right.Finetune)

	delta := period.PeriodDelta(t * (rnft - lnft))
	return p.AddDelta(delta)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Linear) GetSamplerAdd(samplerSpeed float64) float64 {
	return float64(p.GetFrequency()) * samplerSpeed / float64(ITBaseClock)
}

// GetFrequency returns the frequency defined by the period
func (p Linear) GetFrequency() period.Frequency {
	pft := float64(p.Finetune-C5SlideFines) / float64(SlideFinesPerOctave)
	f := p.C2Spd * period.Frequency(math.Pow(2.0, pft))
	return f
}

func (p Linear) String() string {
	return fmt.Sprintf("Linear{ Finetune:%v C2Spd:%v }", p.Finetune, p.C2Spd)
}

// ToLinearPeriod returns the linear frequency period for a given period
func ToLinearPeriod(p period.Period) Linear {
	switch pp := p.(type) {
	case Linear:
		return pp
	case Amiga:
		linFreq := float64(semitonePeriodTable[0]) / float64(pp)

		fts := note.Finetune(SlideFinesPerOctave * math.Log2(linFreq))

		lp := Linear{
			Finetune: fts,
			C2Spd:    DefaultC2Spd,
		}
		return lp
	}
	return Linear{}
}
