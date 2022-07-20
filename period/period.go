package period

import (
	"github.com/heucuva/comparison"
)

// Period is an interface that defines a sampler period
type Period interface {
	AddDelta(Delta) Period
	Compare(Period) comparison.Spaceship // <=>
	Lerp(float64, Period) Period
	GetSamplerAdd(float64) float64
	GetFrequency() Frequency
}

// PeriodDelta is an amount of delta specific to the period type it modifies
// it's intended to be non-specific unit type, so it's up to the implementer
// to keep track of the expected unit type.
type PeriodDelta float64

// ToPeriodDelta works as a conversion system for different types of 'delta' values to a single common one
func ToPeriodDelta(delta Delta) PeriodDelta {
	switch d := delta.(type) {
	case PeriodDelta:
		return d
	case float32:
		return PeriodDelta(d)
	default:
		panic("unknown type conversion for Delta")
	}
}

// ComparePeriods compares two periods, taking nil into account
func ComparePeriods(lhs Period, rhs Period) comparison.Spaceship {
	if lhs == nil {
		if rhs == nil {
			return comparison.SpaceshipEqual
		}
		return comparison.SpaceshipRightGreater
	} else if rhs == nil {
		return comparison.SpaceshipLeftGreater
	}

	return lhs.Compare(rhs)
}
