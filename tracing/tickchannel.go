package tracing

import (
	"fmt"

	"github.com/gotracker/playback/index"
)

type tickChannel struct {
	tick Tick
	ch   index.Channel
}

func (e tickChannel) String() string {
	return fmt.Sprintf("%v %03d", e.tick, e.ch+1)
}

func (e tickChannel) GetTick() Tick {
	return e.tick
}

///////////////////////////////////////////////////////////

func (t *tracerFile) traceChannel(tick Tick, ch index.Channel, op string) {
	t.traceChannelWithComment(tick, ch, op, "")
}

func (t *tracerFile) traceChannelWithComment(tick Tick, ch index.Channel, op string, comment string) {
	if t.file == nil {
		return
	}
	tc := tickChannel{
		tick: tick,
		ch:   ch,
	}
	traceWithPayload(t, tc, op, comment, empty)
}
