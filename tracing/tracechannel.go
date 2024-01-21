package tracing

import "github.com/gotracker/playback/index"

type TraceChannel interface {
	Trace(op string)
	TraceWithComment(op, comment string)
	TraceValueChange(op string, prev, new any)
	TraceValueChangeWithComment(op string, prev, new any, comment string)
}

type channelTracer struct {
	t       Tracer
	channel index.Channel
}

func (c channelTracer) Trace(op string) {
	c.t.TraceChannel(c.channel, op)
}

func (c channelTracer) TraceWithComment(op, comment string) {
	c.t.TraceChannelWithComment(c.channel, op, comment)
}

func (c channelTracer) TraceValueChange(op string, prev, new any) {
	c.t.TraceChannelValueChange(c.channel, op, prev, new)
}

func (c channelTracer) TraceValueChangeWithComment(op string, prev, new any, comment string) {
	c.t.TraceChannelValueChangeWithComment(c.channel, op, prev, new, comment)
}
