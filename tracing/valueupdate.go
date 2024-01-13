package tracing

import (
	"fmt"
	"reflect"

	"github.com/gotracker/playback/index"
)

type valueUpdate struct {
	old any
	new any
}

func (e valueUpdate) String() string {
	return fmt.Sprintf("%v <- %v", e.new, e.old)
}

///////////////////////////////////////////////////////////

func (t *Tracing) traceValueChange(tick Tick, op string, prev, new any) {
	t.traceValueChangeWithComment(tick, op, prev, new, "")
}

func (t *Tracing) traceValueChangeWithComment(tick Tick, op string, prev, new any, comment string) {
	if t.tracingFile == nil {
		return
	}
	if reflect.DeepEqual(prev, new) {
		return
	}
	traceWithPayload(t, tick, op, comment, valueUpdate{
		old: prev,
		new: new,
	})
}

func (t *Tracing) traceChannelValueChange(tick Tick, ch index.Channel, op string, prev, new any) {
	t.traceChannelValueChangeWithComment(tick, ch, op, prev, new, "")
}

func (t *Tracing) traceChannelValueChangeWithComment(tick Tick, ch index.Channel, op string, prev, new any, comment string) {
	if t.tracingFile == nil {
		return
	}
	if reflect.DeepEqual(prev, new) {
		return
	}

	tc := tickChannel{
		tick: tick,
		ch:   ch,
	}
	traceWithPayload(t, tc, op, comment, valueUpdate{
		old: prev,
		new: new,
	})
}
