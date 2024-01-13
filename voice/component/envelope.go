package component

import (
	"fmt"
	"math"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/loop"
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
	updater  func(TIn, TIn, float64) TOut
	enabled  bool
	pos      int
	stopped  bool
	value    TOut

	slimKeyModulator
}

type EnvelopeSettings[TIn, TOut any] struct {
	envelope.Envelope[TIn]
	OnFinished voice.Callback
}

func (e *baseEnvelope[TIn, TOut]) Setup(settings EnvelopeSettings[TIn, TOut], update func(TIn, TIn, float64) TOut) {
	e.settings = settings
	e.updater = update
	e.Reset()
}

func (e baseEnvelope[TIn, TOut]) Clone(update func(TIn, TIn, float64) TOut) baseEnvelope[TIn, TOut] {
	m := e
	return m
}

// Reset resets the state to defaults based on the envelope provided
func (e *baseEnvelope[TIn, TOut]) Reset() {
	e.enabled = e.settings.Enabled
	e.stateReset()
	e.updateValue()
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
	return e.stopped
}

// GetCurrentValue returns the current cached envelope value
func (e baseEnvelope[TIn, TOut]) GetCurrentValue() TOut {
	return e.value
}

// SetEnvelopePosition sets the current position in the envelope
func (e *baseEnvelope[TIn, TOut]) SetEnvelopePosition(pos int) voice.Callback {
	e.stateReset()
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(); doneCB != nil {
			return doneCB
		}
	}
	return nil
}

func (e baseEnvelope[TIn, TOut]) GetEnvelopePosition() int {
	return e.pos
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *baseEnvelope[TIn, TOut]) Advance() voice.Callback {
	var doneCB voice.Callback
	if done := e.stateAdvance(e.keyOn); done {
		doneCB = e.settings.OnFinished
	}
	e.updateValue()
	return doneCB
}

func (e baseEnvelope[TIn, TOut]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} pos{%v} stopped{%v} value{%v}",
		e.enabled,
		e.pos,
		e.stopped,
		e.value,
	), comment)
}

func (e *baseEnvelope[TIn, TOut]) updateValue() {
	if !e.enabled {
		return
	}

	curVal, nextVal, t := e.getCurrentPoints()

	var y0 TIn
	if curVal != nil {
		y0 = curVal.Y
	}

	var y1 TIn
	if nextVal != nil {
		y1 = nextVal.Y
	}

	e.value = e.updater(y0, y1, t)
}

func (e *baseEnvelope[TIn, TOut]) stateReset() {
	if !e.settings.Envelope.Enabled {
		e.stopped = true
		return
	}

	e.pos = 0
	e.stopped = false
}

func (e *baseEnvelope[TIn, TOut]) getCurrentPoints() (*envelope.Point[TIn], *envelope.Point[TIn], float64) {
	if e.stopped {
		return nil, nil, 0
	}

	nPoints := len(e.settings.Envelope.Values)

	if nPoints == 0 {
		return nil, nil, 0
	}

	curTick, _ := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, e.pos, e.settings.Envelope.Length, e.prevKeyOn)
	nextTick, _ := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, e.pos+1, e.settings.Envelope.Length, e.keyOn)

	var cur envelope.Point[TIn]
	for _, it := range e.settings.Envelope.Values {
		if it.Pos > curTick {
			break
		}
		cur = it
	}

	var next envelope.Point[TIn]
	foundNext := false
	for _, it := range e.settings.Envelope.Values {
		if it.Pos > nextTick {
			next = it
			foundNext = true
			break
		}
	}

	if !foundNext {
		return &cur, &cur, 0
	}

	t := float64(0)
	if cur.Length > 0 {
		if tl := curTick - cur.Pos; tl > 0 {
			t = max(min((float64(tl)/float64(cur.Length)), 1), 0)
		}
	}
	return &cur, &next, t
}

func (e *baseEnvelope[TIn, TOut]) stateAdvance(keyOn bool) bool {
	if e.stopped {
		return false
	}

	if e.settings.Envelope.Sustain.Enabled() && keyOn {
		if e.settings.Envelope.Sustain.Length() == 0 {
			return false
		}
	} else if e.settings.Envelope.Loop.Enabled() {
		if e.settings.Envelope.Loop.Length() == 0 {
			return false
		}
	}

	nPoints := len(e.settings.Envelope.Values)

	if nPoints == 0 {
		e.stopped = true
		return true
	}

	e.pos++
	curTick, _ := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, e.pos, e.settings.Envelope.Length, keyOn)

	found := false
	for _, i := range e.settings.Envelope.Values {
		if i.Pos >= curTick && i.Length != math.MaxInt {
			found = true
			break
		}
	}

	if !found {
		e.stopped = true
		return true
	}

	return false
}
