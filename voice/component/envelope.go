package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
)

// Envelope is an envelope component interface
type Envelope interface {
	//Init(env *envelope.Envelope)
	SetEnabled(enabled bool)
	IsEnabled() bool
	Reset()
	Advance()
}

// baseEnvelope is a basic modulation envelope
type baseEnvelope[TIn, TOut any] struct {
	settings EnvelopeSettings[TIn, TOut]
	updater  func() TOut
	enabled  bool
	state    envelope.State[TIn]
	value    TOut

	slimKeyModulator
}

type EnvelopeSettings[TIn, TOut any] struct {
	envelope.Envelope[TIn]
	OnFinished voice.Callback
}

func (e *baseEnvelope[TIn, TOut]) Setup(settings EnvelopeSettings[TIn, TOut], update func() TOut) {
	e.settings = settings
	e.updater = update
	e.Reset()
}

func (e baseEnvelope[TIn, TOut]) Clone(update func() TOut) baseEnvelope[TIn, TOut] {
	m := e
	m.state = e.state.Clone()
	return m
}

// Reset resets the state to defaults based on the envelope provided
func (e *baseEnvelope[TIn, TOut]) Reset() {
	e.enabled = e.settings.Enabled
	e.state.Init(&e.settings.Envelope)
	if e.enabled {
		e.value = e.updater()
	}
}

// SetEnabled sets the enabled flag for the envelope
func (e *baseEnvelope[TIn, TOut]) SetEnabled(enabled bool) {
	e.enabled = enabled
}

// IsEnabled returns the enabled flag for the envelope
func (e baseEnvelope[TIn, TOut]) IsEnabled() bool {
	return e.enabled
}

func (e baseEnvelope[TIn, TOut]) IsDone() bool {
	return e.state.Stopped()
}

// GetCurrentValue returns the current cached envelope value
func (e baseEnvelope[TIn, TOut]) GetCurrentValue() TOut {
	return e.value
}

// SetEnvelopePosition sets the current position in the envelope
func (e *baseEnvelope[TIn, TOut]) SetEnvelopePosition(pos int) voice.Callback {
	e.state.Reset()
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(); doneCB != nil {
			return doneCB
		}
	}
	return nil
}

func (e baseEnvelope[TIn, TOut]) GetEnvelopePosition() int {
	return e.state.Pos()
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *baseEnvelope[TIn, TOut]) Advance() voice.Callback {
	var doneCB voice.Callback
	if done := e.state.Advance(e.keyOn); done {
		doneCB = e.settings.OnFinished
	}
	e.value = e.updater()
	return doneCB
}

func (e baseEnvelope[TIn, TOut]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} pos{%v} value{%v}",
		e.enabled,
		e.state.Pos(),
		e.value,
	), comment)
}
