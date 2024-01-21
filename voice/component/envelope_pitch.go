package component

import (
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/types"
)

// PitchEnvelope is an frequency modulation envelope
type PitchEnvelope struct {
	baseEnvelope[types.PitchFiltValue, period.Delta]
}

func (e *PitchEnvelope) Setup(settings EnvelopeSettings[types.PitchFiltValue, period.Delta]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e PitchEnvelope) Clone(onFinished voice.Callback) PitchEnvelope {
	var m PitchEnvelope
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc, onFinished)
	return m
}

func (e *PitchEnvelope) calc(y0, y1 types.PitchFiltValue, t float64) period.Delta {
	return -period.Delta(util.Lerp(t, y0, y1))
}
