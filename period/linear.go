package period

import (
	"fmt"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"
)

// Linear is a linear period, based on semitone and finetune values
type Linear struct {
	Finetune   note.Finetune
	CommonRate Frequency
}

// Add adds the current period to a delta value then returns the resulting period
func (p Linear) AddDelta(delta Delta) Period {
	panic("unimplemented") // must be implemented by format
}

// Add adds the current period to a delta value then returns the resulting period
func AddLinearDelta(p Linear, delta PeriodDelta) Linear {
	// 0 means "not playing", so keep it that way
	if p.Finetune > 0 {
		p.Finetune += note.Finetune(delta)
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
func (p Linear) Compare(rhs Period) comparison.Spaceship {
	panic("unimplemented") // must be implemented by format
}

// CompareLinear compares (<=>) two linear periods
func CompareLinear(lhs, rhs Linear) comparison.Spaceship {
	switch {
	case lhs.Finetune < rhs.Finetune:
		return comparison.SpaceshipRightGreater
	case lhs.Finetune > rhs.Finetune:
		return comparison.SpaceshipLeftGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Linear) Lerp(t float64, rhs Period) Period {
	right, _ := rhs.ToLinearPeriod().(Linear)

	lnft := float64(p.Finetune)
	rnft := float64(right.Finetune)

	delta := PeriodDelta(t * (rnft - lnft))
	return p.AddDelta(delta)
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Linear) GetSamplerAdd(samplerSpeed float64) float64 {
	panic("unimplemented") // must be implemented by format
}

// GetFrequency returns the frequency defined by the period
func (p Linear) GetFrequency() Frequency {
	panic("unimplemented") // must be implemented by format
}

func (p Linear) String() string {
	return fmt.Sprintf("LinearPeriod{ Finetune:%v CommonRate:%v }", p.Finetune, p.CommonRate)
}

// ToLinearPeriod returns the linear frequency period for a given period
func (p Linear) ToLinearPeriod() Period {
	panic("unimplemented") // must be implemented by format
}

// ToAmigaPeriod returns the amiga (protracker) representation for a given period
func (p Linear) ToAmigaPeriod() Period {
	panic("unimplemented") // must be implemented by format
}
