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

func (e *PitchEnvelope) calc(y0, y1 filter.PitchFiltValue, t float64) period.Delta {
	return -period.Delta(util.Lerp(t, y0, y1))
}
