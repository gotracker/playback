package panning

import (
	"testing"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/voice/types"
)

func TestFromItPanningDisabledUsesCenterAhead(t *testing.T) {
	pos := FromItPanning(itfile.PanValue(0x80))
	if pos != panning.CenterAhead {
		t.Fatalf("expected disabled panning to map to CenterAhead, got %+v", pos)
	}
}

func TestItPanningRoundTrip(t *testing.T) {
	source := itfile.PanValue(32)
	pos := FromItPanning(source)
	if pos.Angle == 0 || pos.Distance != 1 {
		t.Fatalf("unexpected stereo position %+v", pos)
	}
	back := ToItPanning(pos)
	if back != source {
		t.Fatalf("expected round-trip to preserve panning, got %d", back)
	}
}

func TestItPanningClampsOutOfRangeRight(t *testing.T) {
	pos := FromItPanning(itfile.PanValue(0x70))
	back := ToItPanning(pos)
	if back != 64 {
		t.Fatalf("expected out-of-range pan to clamp to 64, got %d", back)
	}
}

func TestPanningClampOperations(t *testing.T) {
	if got := Panning(10).FMA(2, 5); got != 25 {
		t.Fatalf("expected FMA to yield 25, got %d", got)
	}

	if got := Panning(250).FMA(2, 0); got != MaxPanning {
		t.Fatalf("expected FMA to clamp to MaxPanning, got %d", got)
	}

	if got := Panning(0).AddDelta(types.PanDelta(-5)); got != 0 {
		t.Fatalf("expected AddDelta to clamp at 0, got %d", got)
	}

	if got := Panning(240).AddDelta(types.PanDelta(20)); got != MaxPanning {
		t.Fatalf("expected AddDelta to clamp to MaxPanning, got %d", got)
	}
}
