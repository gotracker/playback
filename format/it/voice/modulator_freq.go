package voice

import (
	"github.com/gotracker/playback/period"
)

// == FreqModulator ==

func (v *itVoice[TPeriod]) SetPeriod(period TPeriod) error {
	if period.IsInvalid() {
		return nil
	}
	return v.freq.SetPeriod(period)
}

func (v itVoice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *itVoice[TPeriod]) SetPeriodDelta(delta period.Delta) error {
	return v.freq.SetPeriodDelta(delta)
}

func (v itVoice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetPeriodDelta()
}

func (v itVoice[TPeriod]) GetFinalPeriod() (TPeriod, error) {
	return v.finalPeriod, nil
}
