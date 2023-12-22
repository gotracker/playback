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

	a -= int(d)
	if a < 1 {
		a = 1
	}
	p = Amiga(a)
	return p
}

func (p Amiga) PortaDown(amount int) Amiga {
	return p.Add(Delta(-amount))
}

func (p Amiga) PortaUp(amount int) Amiga {
	return p.Add(Delta(amount))
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

// GetFrequency returns the frequency defined by the period
func (p Amiga) GetFrequency() Frequency {
	panic("unimplemented") // must be implemented by format
}

func (p Amiga) String() string {
	return fmt.Sprintf("Amiga{ Period:%d }", p)
}
