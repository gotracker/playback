package period

import (
	"github.com/heucuva/comparison"
)

// Period is an interface that defines a sampler period
type Period interface {
	IsInvalid() bool
}

type PeriodPorta[TPeriod Period] interface {
	Add(d Delta) TPeriod
	PortaDown(amount int) TPeriod
	PortaUp(amount int) TPeriod
	PortaTo(amount int, target TPeriod) TPeriod
}

// ComparePeriods compares two periods, taking nil into account
func ComparePeriods[TPeriod Period](lhs TPeriod, rhs TPeriod) comparison.Spaceship {
	if lhs.IsInvalid() {
		if rhs.IsInvalid() {
			return comparison.SpaceshipEqual
		}
		return comparison.SpaceshipRightGreater
	} else if rhs.IsInvalid() {
		return comparison.SpaceshipLeftGreater
	}

	switch p := any(lhs).(type) {
	case Linear:
		return p.Compare(any(rhs).(Linear))
	case Amiga:
		return p.Compare(any(rhs).(Amiga))
	default:
		panic("unhandled period type")
	}
}
