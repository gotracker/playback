package component

import (
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice/types"
)

// PanEnvelope is a spatial modulation envelope
type PanEnvelope[TPanning types.Panning] struct {
	baseEnvelope[TPanning, TPanning]
}

func (e *PanEnvelope[TPanning]) Setup(settings EnvelopeSettings[TPanning, TPanning]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e PanEnvelope[TPanning]) Clone() PanEnvelope[TPanning] {
	var m PanEnvelope[TPanning]
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc)
	return m
}

func (e *PanEnvelope[TPanning]) calc(y0, y1 TPanning, t float64) TPanning {
	return util.Lerp(t, y0, y1)
}
