package component

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/util"
)

// PitchEnvelope is an frequency modulation envelope
type PitchEnvelope struct {
	baseEnvelope[filter.PitchFiltValue, period.Delta]
}

func (e *PitchEnvelope) Setup(settings EnvelopeSettings[filter.PitchFiltValue, period.Delta]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e PitchEnvelope) Clone() PitchEnvelope {
	var m PitchEnvelope
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc)
	return m
}

func (e *PitchEnvelope) calc() period.Delta {
	cur, next, t := e.state.GetCurrentValue(e.keyOn, e.prevKeyOn)

	var y0 filter.PitchFiltValue
	if cur != nil {
		y0 = cur.Y
	}

	var y1 filter.PitchFiltValue
	if next != nil {
		y1 = next.Y
	}

	return -period.Delta(util.Lerp(float64(t), y0, y1))
}
