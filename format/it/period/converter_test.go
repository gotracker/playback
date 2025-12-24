package period

import (
	"math"
	"testing"

	"github.com/gotracker/playback/format/it/system"
	"github.com/gotracker/playback/period"
)

func TestConvertersConfigured(t *testing.T) {
	amiga, ok := AmigaConverter.(period.AmigaConverter)
	if !ok {
		t.Fatalf("expected AmigaConverter to be period.AmigaConverter")
	}
	if amiga.MinPeriod != 1 || amiga.MaxPeriod != math.MaxUint16 {
		t.Fatalf("unexpected AmigaConverter bounds: min %d max %d", amiga.MinPeriod, amiga.MaxPeriod)
	}
	if amiga.System != system.ITSystem {
		t.Fatalf("unexpected AmigaConverter system")
	}

	linear, ok := LinearConverter.(period.LinearConverter)
	if !ok {
		t.Fatalf("expected LinearConverter to be period.LinearConverter")
	}
	if linear.System != system.ITSystem {
		t.Fatalf("unexpected LinearConverter system")
	}
}
