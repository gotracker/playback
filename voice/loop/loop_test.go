package loop

import (
	"fmt"
	"testing"
)

type stubLoop struct {
	enabled bool
	pos     int
	looped  bool
}

func (s stubLoop) Enabled() bool                           { return s.enabled }
func (s stubLoop) Length() int                             { return 0 }
func (s stubLoop) CalcPos(pos int, length int) (int, bool) { return s.pos, s.looped }

func TestNewLoopReturnsExpectedType(t *testing.T) {
	cases := []struct {
		mode   Mode
		expect string
	}{
		{ModeDisabled, "*loop.Disabled"},
		{ModeLegacy, "*loop.Legacy"},
		{ModeNormal, "*loop.Normal"},
		{ModePingPong, "*loop.PingPong"},
	}

	for _, tt := range cases {
		l := NewLoop(tt.mode, Settings{})
		if got := fmt.Sprintf("%T", l); got != tt.expect {
			t.Fatalf("mode %v returned %s, want %s", tt.mode, got, tt.expect)
		}
	}
}

func TestDisabledLoopClamps(t *testing.T) {
	l := &Disabled{}
	if pos, looped := l.CalcPos(-1, 5); pos != 0 || looped {
		t.Fatalf("neg pos -> %d looped=%v, want 0 false", pos, looped)
	}
	if pos, looped := l.CalcPos(10, 5); pos != 5 || looped {
		t.Fatalf("past end -> %d looped=%v, want 5 false", pos, looped)
	}
}

func TestNormalLoopWraps(t *testing.T) {
	l := &Normal{Settings: Settings{Begin: 2, End: 5}}
	if pos, looped := l.CalcPos(3, 8); pos != 3 || looped {
		t.Fatalf("inside -> %d looped=%v", pos, looped)
	}
	if pos, looped := l.CalcPos(5, 8); pos != 2 || !looped {
		t.Fatalf("at end -> %d looped=%v, want 2 true", pos, looped)
	}
	if pos, looped := l.CalcPos(7, 8); pos != 4 || !looped {
		t.Fatalf("after end -> %d looped=%v, want 4 true", pos, looped)
	}
}

func TestPingPongLoopBounces(t *testing.T) {
	l := &PingPong{Settings: Settings{Begin: 2, End: 5}}
	if pos, looped := l.CalcPos(5, 8); pos != 4 || !looped {
		t.Fatalf("bounce first -> %d looped=%v, want 4 true", pos, looped)
	}
	if pos, looped := l.CalcPos(6, 8); pos != 3 || !looped {
		t.Fatalf("bounce second -> %d looped=%v, want 3 true", pos, looped)
	}
	if pos, looped := l.CalcPos(8, 8); pos != 2 || !looped {
		t.Fatalf("reverse direction -> %d looped=%v, want 2 true", pos, looped)
	}
}

func TestLegacyLoopAfterEnd(t *testing.T) {
	l := &Legacy{Settings: Settings{Begin: 2, End: 5}}
	if pos, looped := l.CalcPos(9, 8); pos != 3 || !looped {
		t.Fatalf("after end -> %d looped=%v, want 3 true", pos, looped)
	}
}

func TestCalcLoopPosPrefersSustainWhenKeyOn(t *testing.T) {
	sustain := stubLoop{enabled: true, pos: 9, looped: true}
	main := stubLoop{enabled: true, pos: 1, looped: false}
	if pos, looped := CalcLoopPos(main, sustain, 0, 10, true); pos != 9 || !looped {
		t.Fatalf("sustain expected 9/true, got %d/%v", pos, looped)
	}
	if pos, looped := CalcLoopPos(main, sustain, 0, 10, false); pos != 1 || looped {
		t.Fatalf("non-sustain expected 1/false, got %d/%v", pos, looped)
	}
}

func TestCalcLoopPosHandlesNilLoops(t *testing.T) {
	pos, looped := CalcLoopPos(nil, nil, -2, 5, false)
	if pos != 0 || looped {
		t.Fatalf("nil loops should clamp and not loop, got %d/%v", pos, looped)
	}
	pos, looped = CalcLoopPos(nil, nil, 9, 5, false)
	if pos != 5 || looped {
		t.Fatalf("nil loops should clamp at length, got %d/%v", pos, looped)
	}
}
