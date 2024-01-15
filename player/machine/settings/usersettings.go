package settings

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/heucuva/optional"
)

type UserSettings struct {
	tracing.Tracer
	SongLoopCount int
	Start         struct {
		Order optional.Value[index.Order] // default: based on song
		Row   optional.Value[index.Row]   // default: 0
		Tempo int                         // 0 = based on song
		BPM   int                         // 0 = based on song
	}
	PlayUntil struct {
		Order optional.Value[index.Order] // default: based on song
		Row   optional.Value[index.Row]   // default: based on song
	}
	LongChannelOutput    bool
	IgnoreUnknownEffect  bool
	EnableNewNoteActions bool
}

func (s UserSettings) SetTracingTick(order index.Order, row index.Row, tick int) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.SetTracingTick(order, row, tick)
}

func (s UserSettings) Trace(op string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.Trace(op)
}

func (s UserSettings) TraceWithComment(op, comment string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceWithComment(op, comment)
}

func (s UserSettings) TraceValueChange(op string, prev, new any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceValueChange(op, prev, new)
}

func (s UserSettings) TraceValueChangeWithComment(op string, prev, new any, comment string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceValueChangeWithComment(op, prev, new, comment)
}

func (s UserSettings) TraceChannel(ch index.Channel, op string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannel(ch, op)
}

func (s UserSettings) TraceChannelWithComment(ch index.Channel, op, comment string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelWithComment(ch, op, comment)
}

func (s UserSettings) TraceChannelValueChange(ch index.Channel, op string, prev, new any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelValueChange(ch, op, prev, new)
}

func (s UserSettings) TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, comment string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelValueChangeWithComment(ch, op, prev, new, comment)
}
