package period

import (
	"fmt"

	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/util"
	"github.com/heucuva/comparison"
)

// Linear is a linear period, based on semitone and finetune values
type Linear struct {
	Finetune   note.Finetune
	CommonRate Frequency
}

// Add adds the current period to a delta value then returns the resulting period
func (p Linear) Add(d Delta) Linear {
	a := int(p.Finetune)
	if a == 0 {
		// 0 means "not playing", so keep it that way
		return p
	}

	a += int(d)
	if a < 1 {
		a = 1
	}
	p.Finetune = note.Finetune(a)
	return p
}

func (p Linear) PortaDown(amount int) Linear {
	return p.Add(Delta(-amount))
}

func (p Linear) PortaUp(amount int) Linear {
	return p.Add(Delta(amount))
}

func (p Linear) IsInvalid() bool {
	return p.Finetune == 0
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Linear) Compare(rhs Linear) comparison.Spaceship {
	switch {
	case p.Finetune < rhs.Finetune:
		return comparison.SpaceshipRightGreater
	case p.Finetune > rhs.Finetune:
		return comparison.SpaceshipLeftGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Linear) Lerp(t float64, rhs Linear) Period {
	p.Finetune = util.Lerp(t, p.Finetune, rhs.Finetune)
	return p
}

func (p Linear) String() string {
	return fmt.Sprintf("LinearPeriod{ Finetune:%v CommonRate:%v }", p.Finetune, p.CommonRate)
}
