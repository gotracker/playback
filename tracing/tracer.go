package tracing

import "github.com/gotracker/playback/index"

type Tracer interface {
	SetTracingTick(order index.Order, row index.Row, tick int)
	Trace(op string)
	TraceWithComment(op, comment string)
	TraceValueChange(op string, prev, new any)
	TraceValueChangeWithComment(op string, prev, new any, comment string)
	TraceChannel(ch index.Channel, op string)
	TraceChannelWithComment(ch index.Channel, op, comment string)
	TraceChannelValueChange(ch index.Channel, op string, prev, new any)
	TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, comment string)
}
