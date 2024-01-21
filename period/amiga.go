package period

import (
	"fmt"

	"github.com/gotracker/playback/util"
	"github.com/heucuva/comparison"
)

type Amiga uint16

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) Add(d Delta, minPeriod, maxPeriod Amiga, canSlideTo0 bool) Amiga {
	if d == 0 {
		return p
	}
	a := int(p)
	if a == 0 {
		// 0 means "not playing", so keep it that way
		return p
	}

	a -= int(d)
	if a == 0 && canSlideTo0 {
		return 0
	}
	// can't use Clamp() here because we need to clamp negatives
	c := min(Amiga(max(a, int(minPeriod))), maxPeriod)
	if c < 64 {
		_ = c
	}
	return c
}

func (p Amiga) Clamp(minPeriod, maxPeriod Amiga) Amiga {
	if p == 0 {
		return 0
	}
	return min(max(p, minPeriod), maxPeriod)
}

func (p Amiga) PortaDown(amount Delta, minPeriod, maxPeriod Amiga, canSlideTo0 bool) Amiga {
	return p.Add(-amount, minPeriod, maxPeriod, canSlideTo0)
}

func (p Amiga) PortaUp(amount Delta, minPeriod, maxPeriod Amiga, canSlideTo0 bool) Amiga {
	return p.Add(amount, minPeriod, maxPeriod, canSlideTo0)
}

func (p Amiga) PortaTo(amount Delta, target, minPeriod, maxPeriod Amiga) Amiga {
	switch p.Compare(target) {
	case comparison.SpaceshipLeftGreater:
		// porta down to target
		p = p.PortaDown(amount, minPeriod, target, false)
	case comparison.SpaceshipRightGreater:
		// porta up to target
		p = p.PortaUp(amount, target, maxPeriod, false)
	}
	return p
}

func (p Amiga) IsInvalid() bool {
	return p == 0
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Amiga) Compare(rhs Amiga) comparison.Spaceship {
	switch {
	case p < rhs:
		return comparison.SpaceshipLeftGreater
	case p > rhs:
		return comparison.SpaceshipRightGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Amiga) Lerp(t float64, rhs Amiga) Amiga {
	p = util.Lerp(t, p, rhs)
	return p
}

func (p Amiga) String() string {
	if p == 0 {
		return "Amiga{ nil }"
	}
	return fmt.Sprintf("Amiga{ Period:%d }", p)
}
