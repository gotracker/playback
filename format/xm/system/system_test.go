package system

import (
	"testing"

	"github.com/gotracker/playback/note"
)

func TestXMSystemConstants(t *testing.T) {
	if XMBaseClock != DefaultC4SampleRate*C4Period {
		t.Fatalf("expected base clock product, got %v", XMBaseClock)
	}
	if len(semitonePeriodTable) != 12 {
		t.Fatalf("expected 12 semitone periods, got %d", len(semitonePeriodTable))
	}
	if semitonePeriodTable[0] != 27392 {
		t.Fatalf("unexpected first semitone period: %d", semitonePeriodTable[0])
	}
	if XMSystem.GetBaseClock() != XMBaseClock {
		t.Fatalf("expected XMSystem base clock to match")
	}
	if XMSystem.GetCommonPeriod() != C4Period {
		t.Fatalf("expected XMSystem common period to match C4 period")
	}
	if XMSystem.GetCommonRate() != DefaultC4SampleRate {
		t.Fatalf("expected XMSystem common rate to match default C4 sample rate")
	}
	if p, ok := XMSystem.GetSemitonePeriod(note.Key(0)); !ok || p != semitonePeriodTable[0] {
		t.Fatalf("expected first semitone period from XMSystem")
	}
}
