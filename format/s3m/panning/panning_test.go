package panning

import "testing"

func TestPanningClampAndPosition(t *testing.T) {
	if got := Panning(0x02).FMA(2, 1); got != 0x05 {
		t.Fatalf("expected FMA result 0x05, got %d", got)
	}
	if got := Panning(0x0E).FMA(2, 1); got != MaxPanning {
		t.Fatalf("expected FMA clamp to max, got %d", got)
	}
	if got := Panning(0).AddDelta(-5); got != 0 {
		t.Fatalf("expected AddDelta clamp at 0, got %d", got)
	}
	if got := Panning(0x0F).AddDelta(5); got != MaxPanning {
		t.Fatalf("expected AddDelta clamp at max, got %d", got)
	}

	pos := PanningFromS3M(0x08)
	if pos.Distance != 1 {
		t.Fatalf("unexpected distance from panning position: %v", pos)
	}
}
