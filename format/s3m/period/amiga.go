package period

import (
	"fmt"
	"math"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"

	"github.com/gotracker/playback/period"
)

// Amiga defines a sampler period that follows the Amiga-style approach of note
// definition. Useful in calculating resampling.
type Amiga period.Amiga

// AddInteger truncates the current period to an integer and adds the delta integer in
// then returns the resulting period
func (p Amiga) AddInteger(delta int) Amiga {
	ret := Amiga(int(p) + delta)
	// clamp to 64 as minimum
	if ret < 64 {
		ret = 64
	}
	return ret
}

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) AddDelta(delta period.Delta) period.Period {
	ret := p
	d := period.ToPeriodDelta(delta)
	ret += Amiga(d)
	// clamp to 64 as minimum
	if ret < 64 {
		ret = 64
	}
	return ret
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Amiga) Compare(rhs period.Period) comparison.Spaceship {
	lf := p.GetFrequency()
	rf := rhs.GetFrequency()

	switch {
	case lf < rf:
		return comparison.SpaceshipRightGreater
	case lf > rf:
		return comparison.SpaceshipLeftGreater
	default:
		return comparison.SpaceshipEqual
	}
}

// Lerp linear-interpolates the current period with the `rhs` period
func (p Amiga) Lerp(t float64, rhs period.Period) period.Period {
	right := Amiga(0)
	if r, ok := rhs.(Amiga); ok {
		right = r
	}

	ret := Amiga(period.Amiga(p).Lerp(t, period.Amiga(right)))
	return ret
}

// GetSamplerAdd returns the number of samples to advance an instrument by given the period
func (p Amiga) GetSamplerAdd(samplerSpeed float64) float64 {
	return float64(period.Amiga(p).GetFrequency(period.Frequency(samplerSpeed)))
}

// GetFrequency returns the frequency defined by the period
func (p Amiga) GetFrequency() period.Frequency {
	return period.Amiga(p).GetFrequency(period.Frequency(S3MBaseClock))
}

func (p Amiga) String() string {
	return fmt.Sprintf("Amiga{ Period:%f }", float32(p))
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c2spd period.Frequency) Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	pow := math.Pow(2, float64(finetunes)/semitonesPerOctave)
	linFreq := float64(c2spd) * pow / float64(DefaultC2Spd)

	period := Amiga(float64(semitonePeriodTable[0]) / linFreq)
	return period
}
