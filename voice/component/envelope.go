package component

import (
	"fmt"

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
	unkeyed  struct{}
	keyed    struct {
		active bool
		pos    int
		done   bool
	}
	value TOut

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

func (e baseEnvelope[TIn, TOut]) Clone(update func(TIn, TIn, float64) TOut, onFinished voice.Callback) baseEnvelope[TIn, TOut] {
	m := e
	m.settings.OnFinished = onFinished
	m.updater = update
	return m
}

// Reset resets the state to defaults based on the envelope provided
func (e *baseEnvelope[TIn, TOut]) Reset() error {
	e.keyed.active = e.settings.Enabled
	return e.stateReset()
}

func (e baseEnvelope[TIn, TOut]) CanLoop() bool {
	return e.settings.Loop != nil && e.settings.Loop.Enabled()
}

// SetEnabled sets the enabled flag for the envelope
func (e *baseEnvelope[TIn, TOut]) SetEnabled(enabled bool) error {
	e.keyed.active = enabled
	return nil
}

// IsEnabled returns the enabled flag for the envelope
func (e baseEnvelope[TIn, TOut]) IsEnabled() bool {
	return e.keyed.active
}

func (e baseEnvelope[TIn, TOut]) IsDone() bool {
	return e.keyed.done
}

// GetCurrentValue returns the current cached envelope value
func (e baseEnvelope[TIn, TOut]) GetCurrentValue() TOut {
	return e.value
}

// SetEnvelopePosition sets the current position in the envelope
func (e *baseEnvelope[TIn, TOut]) SetEnvelopePosition(pos int) (voice.Callback, error) {
	prev := e.keyed.active
	e.keyed.active = true
	e.keyed.done = false
	e.stateReset()
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(); doneCB != nil {
			return doneCB, nil
		}
	}
	e.keyed.active = prev
	return nil, nil
}

func (e baseEnvelope[TIn, TOut]) GetEnvelopePosition() int {
	return e.keyed.pos
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *baseEnvelope[TIn, TOut]) Advance() voice.Callback {
	var doneCB voice.Callback
	if done := e.stateAdvance(e.keyOn); done {
		doneCB = e.settings.OnFinished
	}
	return doneCB
}

func (e baseEnvelope[TIn, TOut]) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("active{%v} pos{%v} stopped{%v} value{%v}",
		e.keyed.active,
		e.keyed.pos,
		e.keyed.done,
		e.value,
	), comment)
}

func (e *baseEnvelope[TIn, TOut]) stateReset() error {
	if !e.settings.Envelope.Enabled {
		e.keyed.done = true
		return nil
	}

	e.keyed.pos = 0
	e.keyed.done = false
	return e.updateValue()
}

func (e *baseEnvelope[TIn, TOut]) updateValue() error {
	if !e.keyed.active || e.keyed.done {
		return nil
	}

	nPoints := len(e.settings.Envelope.Values)

	if nPoints == 0 {
		return nil
	}

	curTick, _ := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, e.keyed.pos, e.settings.Envelope.Length, e.prevKeyOn)
	nextTick, _ := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, curTick+1, e.settings.Envelope.Length, e.keyOn)

	curPoint := -1
	for i, it := range e.settings.Envelope.Values {
		if it.Pos > curTick {
			curPoint = i - 1
			break
		}
	}
	var cur envelope.Point[TIn]
	if curPoint >= 0 && curPoint < nPoints {
		cur = e.settings.Values[curPoint]
	} else {
		cur = e.settings.Values[nPoints-1]
	}

	nextPoint := -1
	for i, it := range e.settings.Envelope.Values {
		if it.Pos > nextTick {
			nextPoint = i
			break
		}
	}

	if nextPoint < 0 || nextPoint >= nPoints {
		e.value = e.updater(cur.Y, cur.Y, 0)
		return nil
	}

	next := e.settings.Values[nextPoint]

	t := float64(0)
	if cur.Length > 0 {
		if tl := curTick - cur.Pos; tl > 0 {
			t = max(min((float64(tl)/float64(cur.Length)), 1), 0)
		}
	}

	e.value = e.updater(cur.Y, next.Y, t)
	return nil
}

func (e *baseEnvelope[TIn, TOut]) stateAdvance(keyOn bool) bool {
	if e.keyed.done {
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
		e.keyed.done = true
		return true
	}

	e.keyed.pos++
	curTick, looped := loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, e.keyed.pos, e.settings.Envelope.Length, keyOn)

	found := false
	for _, i := range e.settings.Envelope.Values {
		if i.Pos >= curTick {
			found = true
			break
		}
	}

	if !found {
		e.keyed.done = true
		return true
	}

	if !keyOn && !looped && curTick >= e.settings.Length {
		e.keyed.done = false
		return true
	}

	e.updateValue()
	return false
}
