package component

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/util"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
)

// FilterEnvelope is a filter frequency cutoff modulation envelope
type FilterEnvelope struct {
	enabled   bool
	state     envelope.State[filter.PitchFiltValue]
	value     uint8
	keyOn     bool
	prevKeyOn bool
}

func (e *FilterEnvelope) Init(env *envelope.Envelope[filter.PitchFiltValue]) {
	e.state.Init(env)
	e.Reset()
}

func (e FilterEnvelope) Clone() FilterEnvelope {
	return FilterEnvelope{
		enabled:   e.enabled,
		state:     e.state.Clone(),
		value:     e.value,
		keyOn:     false,
		prevKeyOn: false,
	}
}

// Reset resets the state to defaults based on the envelope provided
func (e *FilterEnvelope) Reset() {
	e.state.Reset()
	e.keyOn = false
	e.prevKeyOn = false
	e.update()
}

// SetEnabled sets the enabled flag for the envelope
func (e *FilterEnvelope) SetEnabled(enabled bool) {
	e.enabled = enabled
}

// IsEnabled returns the enabled flag for the envelope
func (e *FilterEnvelope) IsEnabled() bool {
	return e.enabled
}

// GetCurrentValue returns the current cached envelope value
func (e *FilterEnvelope) GetCurrentValue() uint8 {
	return e.value
}

// SetEnvelopePosition sets the current position in the envelope
func (e *FilterEnvelope) SetEnvelopePosition(pos int) voice.Callback {
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
func (e *FilterEnvelope) Advance(keyOn bool, prevKeyOn bool) voice.Callback {
	e.keyOn = keyOn
	e.prevKeyOn = prevKeyOn
	var doneCB voice.Callback
	if done := e.state.Advance(e.keyOn, e.prevKeyOn); done {
		doneCB = e.state.Envelope().OnFinished
	}
	e.update()
	return doneCB
}

func (e *FilterEnvelope) update() {
	cur, next, t := e.state.GetCurrentValue(e.keyOn)

	var y0 filter.PitchFiltValue
	if cur != nil {
		y0 = cur.Value()
	}

	var y1 filter.PitchFiltValue
	if next != nil {
		y1 = next.Value()
	}

	e.value = util.Lerp(float64(t), y0.AsFilter(), y1.AsFilter())
}
