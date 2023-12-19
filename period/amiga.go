package period

import (
	"fmt"

	"github.com/gotracker/playback/util"
	"github.com/heucuva/comparison"
)

type Amiga uint16

// AddInteger truncates the current period to an integer and adds the delta integer in
// then returns the resulting period
func (p Amiga) AddInteger(delta int) Amiga {
	period := Amiga(int(p) + delta)
	return period
}

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) AddDelta(delta Delta) Period {
	panic("unimplemented") // must be implemented by format
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Amiga) Compare(rhs Period) comparison.Spaceship {
	panic("unimplemented") // must be implemented by format
}

func CompareAmiga(lhs, rhs Amiga) comparison.Spaceship {
	switch {
	case lhs < rhs:
		return comparison.SpaceshipLeftGreater
	case lhs > rhs:
		return comparison.SpaceshipRightGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Amiga) Lerp(t float64, rhs Period) Period {
	right := Amiga(0)
	if r, ok := rhs.(Amiga); ok {
		right = r
	}

	return Amiga(util.LerpFloat64(t, float64(p), float64(right)))
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Amiga) GetSamplerAdd(samplerSpeed float64) float64 {
	if p == 0 {
		return 0
	}
	return float64(samplerSpeed) / float64(p)
}

// GetFrequency returns the frequency defined by the period
func (p Amiga) GetFrequency() Frequency {
	panic("unimplemented") // must be implemented by format
}

func (p Amiga) String() string {
	return fmt.Sprintf("Amiga{ Period:%d }", p)
}

// ToLinearPeriod returns the linear frequency period for a given period
func (p Amiga) ToLinearPeriod() Period {
	panic("unimplemented") // must be implemented by format
}

// ToAmigaPeriod returns the amiga (protracker) representation for a given period
func (p Amiga) ToAmigaPeriod() Period {
	panic("unimplemented") // must be implemented by format
}
