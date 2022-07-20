package period

import (
	"github.com/gotracker/playback/util"
)

type AmigaPeriod float64

func (p AmigaPeriod) Lerp(t float64, rhs AmigaPeriod) AmigaPeriod {
	return AmigaPeriod(util.LerpFloat64(t, float64(p), float64(rhs)))
}

func (p AmigaPeriod) GetFrequency(baseClockRate Frequency) Frequency {
	if p == 0 {
		return 0
	}
	return baseClockRate / Frequency(p)
}
