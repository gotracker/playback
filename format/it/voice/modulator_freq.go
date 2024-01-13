package voice

import (
	"github.com/gotracker/playback/period"
)

// == FreqModulator ==

func (v *itVoice[TPeriod]) SetPeriod(period TPeriod) {
	v.freq.SetPeriod(period)
}

func (v *itVoice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *itVoice[TPeriod]) SetPeriodDelta(delta period.Delta) {
	v.freq.SetPeriodDelta(delta)
}

func (v *itVoice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetPeriodDelta()
}

func (v *itVoice[TPeriod]) GetFinalPeriod() TPeriod {
	p := v.freq.GetFinalPeriod()
	if v.IsPitchEnvelopeEnabled() {
		delta := v.GetCurrentPitchEnvelope()
		p = period.AddDelta(p, delta)
	}
	return p
}
