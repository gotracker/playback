package period

import (
	"testing"

	xmSystem "github.com/gotracker/playback/format/xm/system"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

func TestCalcFinetuneC4SampleRateIdentity(t *testing.T) {
	got := CalcFinetuneC4SampleRate(xmSystem.DefaultC4SampleRate, note.Semitone(xmSystem.C4Note), 0)
	if got != xmSystem.DefaultC4SampleRate {
		t.Fatalf("expected unchanged c4 sample rate, got %v", got)
	}
}

func TestCalcFinetuneC4SampleRateAdjusts(t *testing.T) {
	got := CalcFinetuneC4SampleRate(xmSystem.DefaultC4SampleRate, note.Semitone(xmSystem.C4Note), note.Finetune(64))
	if got != 8608 {
		t.Fatalf("expected finetune to raise rate to 8608, got %v", got)
	}
}

func TestConvertersConfigured(t *testing.T) {
	linear, ok := LinearConverter.(period.LinearConverter)
	if !ok {
		t.Fatalf("expected LinearConverter to be period.LinearConverter")
	}
	amiga, ok := AmigaConverter.(period.AmigaConverter)
	if !ok {
		t.Fatalf("expected AmigaConverter to be period.AmigaConverter")
	}

	if linear.System != xmSystem.XMSystem {
		t.Fatalf("unexpected LinearConverter system")
	}
	if amiga.MaxPeriod != 31999 || amiga.MinPeriod != 1 {
		t.Fatalf("unexpected AmigaConverter bounds: min %d max %d", amiga.MinPeriod, amiga.MaxPeriod)
	}
	if amiga.System != xmSystem.XMSystem {
		t.Fatalf("unexpected AmigaConverter system")
	}
}
