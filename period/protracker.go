package period

import (
	"github.com/gotracker/playback/util"
)

type Amiga uint16

func (p Amiga) Lerp(t float64, rhs Amiga) Amiga {
	return Amiga(util.LerpFloat64(t, float64(p), float64(rhs)))
}

func (p Amiga) GetFrequency(baseClockRate Frequency) Frequency {
	if p == 0 {
		return 0
	}
	return baseClockRate / Frequency(p)
}
