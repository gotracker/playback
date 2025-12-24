package system

import (
	"testing"

	"github.com/gotracker/playback/note"
)

func TestITBaseClockCalculation(t *testing.T) {
	if ITBaseClock != DefaultC5SampleRate*C5Period {
		t.Fatalf("expected ITBaseClock to equal DefaultC5SampleRate*C5Period, got %v", ITBaseClock)
	}
}

func TestITSystemValues(t *testing.T) {
	s := ITSystem
	if s.GetBaseClock() != ITBaseClock {
		t.Fatalf("unexpected base clock: %v", s.GetBaseClock())
	}
	if s.GetCommonRate() != DefaultC5SampleRate {
		t.Fatalf("unexpected common rate: %v", s.GetCommonRate())
	}
	if s.GetCommonPeriod() != C5Period {
		t.Fatalf("unexpected common period: %d", s.GetCommonPeriod())
	}
	if s.GetMaxPastNotesPerChannel() != 1 {
		t.Fatalf("unexpected max past notes: %d", s.GetMaxPastNotesPerChannel())
	}
	if got := s.GetOctaveShift(); got != 0 {
		t.Fatalf("expected zero octave shift, got %d", got)
	}
}

func TestITSystemSemitonePeriods(t *testing.T) {
	s := ITSystem
	period, ok := s.GetSemitonePeriod(note.KeyC)
	if !ok {
		t.Fatalf("expected semitone period for KeyC")
	}
	if period != semitonePeriodTable[0] {
		t.Fatalf("unexpected semitone period for KeyC: %d", period)
	}
	if _, ok := s.GetSemitonePeriod(note.KeyInvalid1); ok {
		t.Fatalf("expected invalid key to report missing semitone period")
	}
}
