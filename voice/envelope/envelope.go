package envelope

import (
	"math"

	"github.com/gotracker/playback/voice/loop"
)

// State is the state information about an envelope
type State[T any] struct {
	tick    int
	stopped bool
	env     *Envelope[T]
}

func (e *State[T]) Init(env *Envelope[T]) {
	e.env = env
	if e.env == nil || !e.env.Enabled {
		e.stopped = true
		return
	}

	e.Reset()
}

func (e State[T]) Clone() State[T] {
	return State[T]{
		tick:    e.tick,
		stopped: e.stopped,
		env:     e.env,
	}
}

// Stopped returns true if the envelope state is stopped
func (e *State[T]) Stopped() bool {
	return e.stopped
}

// Stop stops the envelope state
func (e *State[T]) Stop() {
	e.stopped = true
}

// Envelope returns the envelope that the state is based on
func (e *State[T]) Envelope() *Envelope[T] {
	return e.env
}

// Reset resets the envelope
func (e *State[T]) Reset() {
	if e.env == nil || !e.env.Enabled {
		e.stopped = true
		return
	}

	e.tick = 0
	e.stopped = false
}

func (e *State[T]) Pos() int {
	return e.tick
}

// GetCurrentValue returns the current value
func (e *State[T]) GetCurrentValue(keyOn, prevKeyOn bool) (*EnvPoint[T], *EnvPoint[T], float64) {
	if e.stopped {
		return nil, nil, 0
	}

	nPoints := len(e.env.Values)

	if nPoints == 0 {
		return nil, nil, 0
	}

	curTick, _ := loop.CalcLoopPos(e.env.Loop, e.env.Sustain, e.tick, e.env.Length, prevKeyOn)
	nextTick, _ := loop.CalcLoopPos(e.env.Loop, e.env.Sustain, e.tick+1, e.env.Length, keyOn)

	var cur EnvPoint[T]
	for _, it := range e.env.Values {
		if it.Pos > curTick {
			break
		}
		cur = it
	}

	var next EnvPoint[T]
	foundNext := false
	for _, it := range e.env.Values {
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

// Advance advances the state by 1 tick
func (e *State[T]) Advance(keyOn bool) bool {
	if e.stopped {
		return false
	}

	if e.env.Sustain.Enabled() && keyOn {
		if e.env.Sustain.Length() == 0 {
			return false
		}
	} else if e.env.Loop.Enabled() {
		if e.env.Loop.Length() == 0 {
			return false
		}
	}

	nPoints := len(e.env.Values)

	if nPoints == 0 {
		e.stopped = true
		return true
	}

	e.tick++
	curTick, _ := loop.CalcLoopPos(e.env.Loop, e.env.Sustain, e.tick, e.env.Length, keyOn)

	found := false
	for _, i := range e.env.Values {
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
