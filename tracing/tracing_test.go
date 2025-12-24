package tracing

import (
	"os"
	"testing"
)

func TestTickEqualsAndString(t *testing.T) {
	a := Tick{Order: 1, Row: 2, Tick: 3}
	b := Tick{Order: 1, Row: 2, Tick: 3}
	c := Tick{Order: 1, Row: 2, Tick: 4}

	if !a.Equals(b) {
		t.Fatalf("expected ticks to be equal")
	}
	if a.Equals(c) {
		t.Fatalf("expected ticks to differ")
	}
	if got := a.String(); got != "001:002  3" {
		t.Fatalf("unexpected tick string: %q", got)
	}
}

func TestTracerSkipsWhenFileNil(t *testing.T) {
	tf := tracerFile{}
	tf.traceWithComment(Tick{}, "op", "comment")
	tf.traceValueChange(Tick{}, "op", 1, 2)
	if len(tf.updates) != 0 {
		t.Fatalf("expected no updates when file is nil, got %d", len(tf.updates))
	}
}

func TestTracerValueChangeIgnoresEqual(t *testing.T) {
	f, err := os.CreateTemp("", "trace-test-*.log")
	if err != nil {
		t.Fatalf("temp file error: %v", err)
	}
	defer os.Remove(f.Name())
	tf := tracerFile{file: f}

	tf.traceValueChange(Tick{}, "op", 1, 1)
	tf.traceValueChange(Tick{}, "op", 1, 2)

	if len(tf.updates) != 1 {
		t.Fatalf("expected only unequal change to be recorded, got %d", len(tf.updates))
	}
}
