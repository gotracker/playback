package component

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/util"
)

// FilterEnvelope is a filter frequency cutoff modulation envelope
type FilterEnvelope struct {
	baseEnvelope[filter.PitchFiltValue, uint8]
}

func (e *FilterEnvelope) Setup(settings EnvelopeSettings[filter.PitchFiltValue, uint8]) {
	e.baseEnvelope.Setup(settings, e.calc)
}

func (e FilterEnvelope) Clone() FilterEnvelope {
	var m FilterEnvelope
	m.baseEnvelope = e.baseEnvelope.Clone(m.calc)
	return m
}

func (e *FilterEnvelope) calc() uint8 {
	cur, next, t := e.state.GetCurrentValue(e.keyOn, e.prevKeyOn)

	var y0 filter.PitchFiltValue
	if cur != nil {
		y0 = cur.Y
	}

	var y1 filter.PitchFiltValue
	if next != nil {
		y1 = next.Y
	}

	v := util.Lerp(float64(t), y0, y1)
	return uint8(32 + v)
}
