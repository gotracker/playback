package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gotracker/playback/index"
)

type fakeTracer struct {
	calls    []string
	closed   bool
	lastTick struct {
		order index.Order
		row   index.Row
		tick  int
	}
}

func (f *fakeTracer) Close() error {
	f.closed = true
	return nil
}

func (f *fakeTracer) OutputTraces() {
	f.calls = append(f.calls, "output")
}

func (f *fakeTracer) SetTracingTick(order index.Order, row index.Row, tick int) {
	f.calls = append(f.calls, "tick")
	f.lastTick = struct {
		order index.Order
		row   index.Row
		tick  int
	}{order, row, tick}
}

func (f *fakeTracer) Trace(op string) {
	f.calls = append(f.calls, "trace:"+op)
}

func (f *fakeTracer) TraceWithComment(op, commentFmt string, commentParams ...any) {
	f.calls = append(f.calls, "tracecomment:"+op)
}

func (f *fakeTracer) TraceValueChange(op string, prev, new any) {
	f.calls = append(f.calls, "traceval:"+op)
}

func (f *fakeTracer) TraceValueChangeWithComment(op string, prev, new any, commentFmt string, commentParams ...any) {
	f.calls = append(f.calls, "tracevalcomment:"+op)
}

func (f *fakeTracer) TraceChannel(ch index.Channel, op string) {
	f.calls = append(f.calls, "channel:"+op)
}

func (f *fakeTracer) TraceChannelWithComment(ch index.Channel, op, commentFmt string, commentParams ...any) {
	f.calls = append(f.calls, "channelcomment:"+op)
}

func (f *fakeTracer) TraceChannelValueChange(ch index.Channel, op string, prev, new any) {
	f.calls = append(f.calls, "channelval:"+op)
}

func (f *fakeTracer) TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, commentFmt string, commentParams ...any) {
	f.calls = append(f.calls, "channelvalcomment:"+op)
}

func TestUserSettingsReset(t *testing.T) {
	ft := &fakeTracer{}
	s := UserSettings{Tracer: ft, SongLoopCount: 2}
	s.Start.Order.Set(index.Order(5))
	s.Start.Row.Set(index.Row(6))
	s.Start.Tempo = 7
	s.Start.BPM = 8
	s.PlayUntil.Order.Set(index.Order(9))
	s.PlayUntil.Row.Set(index.Row(10))
	s.LongChannelOutput = false
	s.IgnoreUnknownEffect = true
	s.EnableNewNoteActions = false

	s.Reset()

	if s.Tracer != ft {
		t.Fatalf("expected tracer preserved after reset")
	}
	if s.SongLoopCount != 0 {
		t.Fatalf("expected SongLoopCount reset to 0, got %d", s.SongLoopCount)
	}
	if _, ok := s.Start.Order.Get(); ok {
		t.Fatalf("expected Start.Order cleared")
	}
	if _, ok := s.Start.Row.Get(); ok {
		t.Fatalf("expected Start.Row cleared")
	}
	if s.Start.Tempo != 0 || s.Start.BPM != 0 {
		t.Fatalf("expected Start tempo/BPM reset")
	}
	if _, ok := s.PlayUntil.Order.Get(); ok {
		t.Fatalf("expected PlayUntil.Order cleared")
	}
	if _, ok := s.PlayUntil.Row.Get(); ok {
		t.Fatalf("expected PlayUntil.Row cleared")
	}
	if !s.LongChannelOutput {
		t.Fatalf("expected LongChannelOutput default true")
	}
	if s.IgnoreUnknownEffect {
		t.Fatalf("expected IgnoreUnknownEffect default false")
	}
	if !s.EnableNewNoteActions {
		t.Fatalf("expected EnableNewNoteActions default true")
	}
}

func TestUserSettingsTraceDelegation(t *testing.T) {
	ft := &fakeTracer{}
	s := UserSettings{Tracer: ft}

	s.OutputTraces()
	s.SetTracingTick(1, 2, 3)
	s.Trace("a")
	s.TraceWithComment("b", "fmt")
	s.TraceValueChange("c", 1, 2)
	s.TraceValueChangeWithComment("d", 1, 2, "fmt")
	s.TraceChannel(4, "e")
	s.TraceChannelWithComment(5, "f", "fmt")
	s.TraceChannelValueChange(6, "g", 1, 2)
	s.TraceChannelValueChangeWithComment(7, "h", 1, 2, "fmt")

	if len(ft.calls) != 10 {
		t.Fatalf("expected 10 tracer calls, got %d", len(ft.calls))
	}
	if ft.lastTick.order != 1 || ft.lastTick.row != 2 || ft.lastTick.tick != 3 {
		t.Fatalf("unexpected SetTracingTick args: %+v", ft.lastTick)
	}
	if ft.closed {
		t.Fatalf("tracer should not be closed yet")
	}
	if err := s.CloseTracing(); err != nil {
		t.Fatalf("CloseTracing returned error: %v", err)
	}
	if !ft.closed {
		t.Fatalf("expected tracer closed after CloseTracing")
	}
}

func TestUserSettingsTracingNilGuard(t *testing.T) {
	var s UserSettings
	s.OutputTraces()
	s.SetTracingTick(0, 0, 0)
	s.Trace("noop")
	s.TraceWithComment("noop", "fmt")
	s.TraceValueChange("noop", 1, 2)
	s.TraceValueChangeWithComment("noop", 1, 2, "fmt")
	s.TraceChannel(0, "noop")
	s.TraceChannelWithComment(0, "noop", "fmt")
	s.TraceChannelValueChange(0, "noop", 1, 2)
	s.TraceChannelValueChangeWithComment(0, "noop", 1, 2, "fmt")
	if err := s.CloseTracing(); err != nil {
		t.Fatalf("CloseTracing returned error on nil tracer: %v", err)
	}
}

func TestSetupTracingWithFilename(t *testing.T) {
	var s UserSettings
	path := filepath.Join(t.TempDir(), "trace.log")
	if err := s.SetupTracingWithFilename(path); err != nil {
		t.Fatalf("SetupTracingWithFilename error: %v", err)
	}
	if s.Tracer == nil {
		t.Fatalf("expected tracer to be created")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected trace file to exist: %v", err)
	}
	if err := s.CloseTracing(); err != nil {
		t.Fatalf("CloseTracing error: %v", err)
	}
}
