package component

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
)

// PitchEnvelope is an frequency modulation envelope
type PitchEnvelope[TPeriod period.Period] struct {
	enabled   bool
	state     envelope.State[filter.PitchFiltValue]
	delta     period.Delta
	keyOn     bool
	prevKeyOn bool
}

func (e *PitchEnvelope[TPeriod]) Init(env *envelope.Envelope[filter.PitchFiltValue]) {
	e.state.Init(env)
	e.Reset()
}

func (e PitchEnvelope[TPeriod]) Clone() PitchEnvelope[TPeriod] {
	return PitchEnvelope[TPeriod]{
		enabled:   e.enabled,
		state:     e.state.Clone(),
		delta:     e.delta,
		keyOn:     false,
		prevKeyOn: false,
	}
}

// Reset resets the state to defaults based on the envelope provided
func (e *PitchEnvelope[TPeriod]) Reset() {
	e.state.Reset()
	e.keyOn = false
	e.prevKeyOn = false
	e.update()
}

// SetEnabled sets the enabled flag for the envelope
func (e *PitchEnvelope[TPeriod]) SetEnabled(enabled bool) {
	e.enabled = enabled
}

// IsEnabled returns the enabled flag for the envelope
func (e *PitchEnvelope[TPeriod]) IsEnabled() bool {
	return e.enabled
}

// GetCurrentValue returns the current cached envelope value
func (e *PitchEnvelope[TPeriod]) GetCurrentValue() period.Delta {
	return e.delta
}

// SetEnvelopePosition sets the current position in the envelope
func (e *PitchEnvelope[TPeriod]) SetEnvelopePosition(pos int) voice.Callback {
	keyOn := e.keyOn
	prevKeyOn := e.prevKeyOn
	e.state.Reset()
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(keyOn, prevKeyOn); doneCB != nil {
			return doneCB
		}
	}
	return nil
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *PitchEnvelope[TPeriod]) Advance(keyOn bool, prevKeyOn bool) voice.Callback {
	e.keyOn = keyOn
	e.prevKeyOn = prevKeyOn
	var doneCB voice.Callback
	if done := e.state.Advance(e.keyOn, e.prevKeyOn); done {
		doneCB = e.state.Envelope().OnFinished
	}
	e.update()
	return doneCB
}

func (e *PitchEnvelope[TPeriod]) update() {
	cur, next, t := e.state.GetCurrentValue(e.keyOn)

	var y0 filter.PitchFiltValue
	if cur != nil {
		y0 = cur.Value()
	}

	var y1 filter.PitchFiltValue
	if next != nil {
		y1 = next.Value()
	}

	e.delta = -period.Delta(util.Lerp(float64(t), y0.AsPitch(), y1.AsPitch()))
}
