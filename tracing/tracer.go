package tracing

import (
	"io"
	"os"

	"github.com/gotracker/playback/index"
)

type Tracer interface {
	OutputTraces()
	SetTracingTick(order index.Order, row index.Row, tick int)
	Trace(op string)
	TraceWithComment(op, commentFmt string, commentParams ...any)
	TraceValueChange(op string, prev, new any)
	TraceValueChangeWithComment(op string, prev, new any, commentFmt string, commentParams ...any)
	TraceChannel(ch index.Channel, op string)
	TraceChannelWithComment(ch index.Channel, op, commentFmt string, commentParams ...any)
	TraceChannelValueChange(ch index.Channel, op string, prev, new any)
	TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, commentFmt string, commentParams ...any)
}

type TracerWithClose interface {
	Tracer
	io.Closer
}

func NewFromFilename(filename string) (TracerWithClose, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	tf := tracerFile{
		file: f,
	}
	return &tf, nil
}
