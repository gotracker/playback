package component

import (
	"fmt"
	"sort"

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
	e.ensureDefaults()
	e.updater = update
	e.Reset()
}

func (e baseEnvelope[TIn, TOut]) Clone(update func(TIn, TIn, float64) TOut, onFinished voice.Callback) baseEnvelope[TIn, TOut] {
	m := e
	m.settings.OnFinished = onFinished
	m.ensureDefaults()
	m.updater = update
	return m
}

// Reset resets the state to defaults based on the envelope provided
func (e *baseEnvelope[TIn, TOut]) Reset() error {
	e.keyed.active = e.settings.Enabled
	return e.stateReset()
}

func (e *baseEnvelope[TIn, TOut]) ensureDefaults() {
	if e.settings.Envelope.Loop == nil {
		e.settings.Envelope.Loop = &loop.Disabled{}
	}
	if e.settings.Envelope.Sustain == nil {
		e.settings.Envelope.Sustain = &loop.Disabled{}
	}
	if e.settings.Envelope.Length < 0 {
		e.settings.Envelope.Length = 0
	}
}

func (e baseEnvelope[TIn, TOut]) loopPos(pos int, keyOn bool) (int, bool) {
	return loop.CalcLoopPos(e.settings.Envelope.Loop, e.settings.Envelope.Sustain, pos, e.settings.Envelope.Length, keyOn)
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
	// XXX: this is gross, but currently the most optimal way to find the correct position
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
	e.ensureDefaults()
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

	curTick, _ := e.loopPos(e.keyed.pos, e.prevKeyOn)
	vals := e.settings.Envelope.Values
	curIdx := sort.Search(len(vals), func(i int) bool { return vals[i].Pos > curTick }) - 1
	if curIdx < 0 {
		curIdx = 0
	} else if curIdx >= nPoints {
		curIdx = nPoints - 1
	}

	nextIdx := curIdx + 1
	if nextIdx >= nPoints {
		e.value = e.updater(vals[curIdx].Y, vals[curIdx].Y, 0)
		return nil
	}

	cur := vals[curIdx]
	next := vals[nextIdx]

	t := float64(0)
	if cur.Length > 0 && curTick >= cur.Pos {
		if tl := curTick - cur.Pos; tl > 0 {
			t = max(min(float64(tl)/float64(cur.Length), 1), 0)
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
	curTick, looped := e.loopPos(e.keyed.pos, keyOn)

	vals := e.settings.Envelope.Values
	idx := sort.Search(len(vals), func(i int) bool { return vals[i].Pos >= curTick })
	if idx == len(vals) && !looped && !keyOn && curTick >= e.settings.Length {
		e.keyed.done = true
		return true
	}

	e.updateValue()
	return false
}
