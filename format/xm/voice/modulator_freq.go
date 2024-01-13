package voice

import (
	"github.com/gotracker/playback/period"
)

// == FreqModulator ==

func (v *xmVoice[TPeriod]) SetPeriod(period TPeriod) {
	v.freq.SetPeriod(period)
}

func (v *xmVoice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *xmVoice[TPeriod]) SetPeriodDelta(delta period.Delta) {
	v.freq.SetPeriodDelta(delta)
}

func (v *xmVoice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetPeriodDelta()
}

func (v *xmVoice[TPeriod]) GetFinalPeriod() TPeriod {
	p := v.freq.GetFinalPeriod()
	if v.IsPitchEnvelopeEnabled() {
		delta := v.GetCurrentPitchEnvelope()
		p = period.AddDelta(p, delta)
	}
	return p
}
