package component

import (
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/types"
)

// FilterEnvelope is a filter frequency cutoff modulation envelope
type FilterEnvelope struct {
	baseEnvelope[types.PitchFiltValue, uint8]
}

func (e *FilterEnvelope) Setup(settings EnvelopeSettings[types.PitchFiltValue, uint8]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e FilterEnvelope) Clone(onFinished voice.Callback) FilterEnvelope {
	var m FilterEnvelope
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc, onFinished)
	return m
}

func (e *FilterEnvelope) calc(y0, y1 types.PitchFiltValue, t float64) uint8 {
	v := util.Lerp(t, y0, y1)
	return uint8(32 + v)
}
