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

func (e *PanEnvelope[TPanning]) calc() TPanning {
	cur, next, t := e.state.GetCurrentValue(e.keyOn, e.prevKeyOn)

	var y0 TPanning
	if cur != nil {
		y0 = cur.Y
	}

	var y1 TPanning
	if next != nil {
		y1 = next.Y
	}

	return util.Lerp(float64(t), y0, y1)
}
