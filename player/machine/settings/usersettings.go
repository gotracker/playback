package settings

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/heucuva/optional"
)

type UserSettings struct {
	Tracer        tracing.TracerWithClose
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

// Reset applies the defaults
//
//	NOTE: does not reset the Tracer value
func (s *UserSettings) Reset() {
	// don't touch the Tracer here
	s.SongLoopCount = 0
	s.Start.Order.Reset()
	s.Start.Row.Reset()
	s.Start.Tempo = 0
	s.Start.BPM = 0
	s.PlayUntil.Order.Reset()
	s.PlayUntil.Row.Reset()
	s.LongChannelOutput = true
	s.IgnoreUnknownEffect = false
	s.EnableNewNoteActions = true
}

func (s *UserSettings) SetupTracingWithFilename(filename string) error {
	var err error
	s.Tracer, err = tracing.NewFromFilename(filename)
	return err
}

func (s *UserSettings) OutputTraces() {
	if s.Tracer != nil {
		s.Tracer.OutputTraces()
	}
}

func (s *UserSettings) CloseTracing() error {
	if s.Tracer != nil {
		return s.Tracer.Close()
	}
	return nil
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

func (s UserSettings) TraceWithComment(op, commentFmt string, commentParams ...any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceWithComment(op, commentFmt, commentParams...)
}

func (s UserSettings) TraceValueChange(op string, prev, new any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceValueChange(op, prev, new)
}

func (s UserSettings) TraceValueChangeWithComment(op string, prev, new any, commentFmt string, commentParams ...any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceValueChangeWithComment(op, prev, new, commentFmt, commentParams...)
}

func (s UserSettings) TraceChannel(ch index.Channel, op string) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannel(ch, op)
}

func (s UserSettings) TraceChannelWithComment(ch index.Channel, op, commentFmt string, commentParams ...any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelWithComment(ch, op, commentFmt, commentParams...)
}

func (s UserSettings) TraceChannelValueChange(ch index.Channel, op string, prev, new any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelValueChange(ch, op, prev, new)
}

func (s UserSettings) TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, commentFmt string, commentParams ...any) {
	if s.Tracer == nil {
		return
	}

	s.Tracer.TraceChannelValueChangeWithComment(ch, op, prev, new, commentFmt, commentParams...)
}
