package period

import (
	"math"

	"github.com/gotracker/playback/note"
	"github.com/heucuva/comparison"

	"github.com/gotracker/playback/period"
)

// Amiga defines a sampler period that follows the Amiga-style approach of note
// definition. Useful in calculating resampling.
type Amiga struct {
	period.Amiga
}

var _ period.Period = (*Amiga)(nil)

// Add adds the current period to a delta value then returns the resulting period
func (p Amiga) AddDelta(delta period.Delta) period.Period {
	d := period.ToPeriodDelta(delta)
	p.Amiga += period.Amiga(d)
	return p
}

// Compare returns:
//
//	-1 if the current period is higher frequency than the `rhs` period
//	0 if the current period is equal in frequency to the `rhs` period
//	1 if the current period is lower frequency than the `rhs` period
func (p Amiga) Compare(rhs period.Period) comparison.Spaceship {
	if q, ok := rhs.ToAmigaPeriod().(Amiga); ok {
		return period.CompareAmiga(p.Amiga, q.Amiga)
	}
	return comparison.SpaceshipLeftGreater
}

// GetFrequency returns the frequency defined by the period
func (p Amiga) GetFrequency() period.Frequency {
	if p.Amiga == 0 {
		return 0
	}
	return period.Frequency(ITBaseClock) / period.Frequency(p.Amiga)
}

// ToLinearPeriod returns the linear frequency period for a given period
func (p Amiga) ToLinearPeriod() period.Period {
	return nil
}

// ToAmigaPeriod returns the amiga (protracker) representation for a given period
func (p Amiga) ToAmigaPeriod() period.Period {
	return p
}

// ToAmigaPeriod calculates an amiga period for a linear finetune period
func ToAmigaPeriod(finetunes note.Finetune, c2spd period.Frequency) Amiga {
	if finetunes < 0 {
		finetunes = 0
	}
	pow := math.Pow(2, float64(finetunes)/SlideFinesPerOctave)
	linFreq := float64(c2spd) * pow / float64(DefaultC2Spd)

	return Amiga{
		Amiga: period.Amiga(float64(semitonePeriodTable[0]) / linFreq),
	}
}
