package period

import (
	"math"

	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// Linear is a linear period, based on semitone and finetune values
type Linear struct {
	period.Linear
}

var _ period.Period = (*Linear)(nil)

// Add adds the current period to a delta value then returns the resulting period
func (p Linear) Add(delta period.PeriodDelta) *Linear {
	p.Linear = period.AddLinearDelta(p.Linear, delta)
	return &p
}

// AddDelta adds the current period to a delta value then returns the resulting period
func (p Linear) AddDelta(delta period.Delta) period.Period {
	return p.Add(period.ToPeriodDelta(delta))
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Linear) Compare(rhs period.Period) comparison.Spaceship {
	if q, ok := rhs.ToLinearPeriod().(Linear); ok {
		return period.CompareLinear(p.Linear, q.Linear)
	}
	return comparison.SpaceshipLeftGreater
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Linear) GetSamplerAdd(samplerSpeed float64) float64 {
	return float64(p.GetFrequency()) * samplerSpeed / float64(ITBaseClock)
}

// GetFrequency returns the frequency defined by the period
func (p Linear) GetFrequency() period.Frequency {
	pft := float64(p.Finetune-C5SlideFines) / float64(SlideFinesPerOctave)
	f := p.CommonRate * period.Frequency(math.Pow(2.0, pft))
	return f
}

// ToLinearPeriod returns the linear frequency period for a given period
func (p Linear) ToLinearPeriod() period.Period {
	return p
}

// ToAmigaPeriod returns the amiga (protracker) representation for a given period
func (p Linear) ToAmigaPeriod() period.Period {
	return nil
}
