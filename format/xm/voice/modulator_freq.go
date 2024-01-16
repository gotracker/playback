package voice

import (
	"github.com/gotracker/playback/period"
)

// == FreqModulator ==

func (v *xmVoice[TPeriod]) SetPeriod(period TPeriod) error {
	return v.freq.SetPeriod(period)
}

func (v *xmVoice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *xmVoice[TPeriod]) SetPeriodDelta(delta period.Delta) error {
	return v.freq.SetPeriodDelta(delta)
}

func (v *xmVoice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetPeriodDelta()
}

func (v *xmVoice[TPeriod]) GetFinalPeriod() (TPeriod, error) {
	p, err := v.freq.GetFinalPeriod()
	if err != nil {
		return p, err
	}
	if v.IsPitchEnvelopeEnabled() {
		delta := v.GetCurrentPitchEnvelope()
		p, err = v.inst.Static.PC.AddDelta(p, delta)
	}
	return p, err
}
