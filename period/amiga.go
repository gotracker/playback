package period

import (
	"fmt"

	"github.com/gotracker/playback/util"
	"github.com/heucuva/comparison"
)

type Amiga uint16

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) Add(d Delta) Amiga {
	a := int(p)
	if a == 0 {
		// 0 means "not playing", so keep it that way
		return p
	}

	p = Amiga(max(a-int(d), 1))
	return p
}

func (p Amiga) PortaDown(amount Delta) Amiga {
	return p.Add(-amount)
}

func (p Amiga) PortaUp(amount Delta) Amiga {
	return p.Add(amount)
}

func (p Amiga) PortaTo(amount Delta, target Amiga) Amiga {
	switch p.Compare(target) {
	case comparison.SpaceshipLeftGreater:
		// porta down to target
		p = p.PortaDown(amount)
		if p.Compare(target) == comparison.SpaceshipRightGreater {
			return target
		}
	case comparison.SpaceshipRightGreater:
		// porta up to target
		p = p.PortaUp(amount)
		if p.Compare(target) == comparison.SpaceshipLeftGreater {
			return target
		}
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
