package panning

import (
	"testing"

	"github.com/gotracker/playback/mixing/panning"
)

func TestPanningRoundTripAndClamp(t *testing.T) {
	pos := PanningFromXm(DefaultPanning)
	if got := PanningToXm(pos); got != uint8(DefaultPanning) {
		t.Fatalf("expected round trip panning %d, got %d", DefaultPanning, got)
	}

	if got := Panning(0xFF).AddDelta(5); got != MaxPanning {
		t.Fatalf("expected AddDelta clamp to MaxPanning, got %d", got)
	}
	if got := Panning(1).AddDelta(-5); got != 0 {
		t.Fatalf("expected AddDelta clamp to 0, got %d", got)
	}

	posLeft := PanningFromXm(DefaultPanningLeft)
	if panning.FromStereoPosition(posLeft, 0, 0xFF) >= 0x80 {
		t.Fatalf("expected left default to be left-biased")
	}
}
