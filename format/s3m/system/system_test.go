package system

import (
	"testing"

	"github.com/gotracker/playback/note"
)

func TestS3MBaseClockCalculation(t *testing.T) {
	if S3MBaseClock != DefaultC4SampleRate*C4Period {
		t.Fatalf("expected S3MBaseClock to equal DefaultC4SampleRate*C4Period, got %v", S3MBaseClock)
	}
}

func TestS3MSystemValues(t *testing.T) {
	s := S3MSystem
	if s.GetBaseClock() != S3MBaseClock {
		t.Fatalf("unexpected base clock: %v", s.GetBaseClock())
	}
	if s.GetCommonRate() != DefaultC4SampleRate {
		t.Fatalf("unexpected common rate: %v", s.GetCommonRate())
	}
	if s.GetCommonPeriod() != C4Period {
		t.Fatalf("unexpected common period: %d", s.GetCommonPeriod())
	}
	if s.GetMaxPastNotesPerChannel() != 0 {
		t.Fatalf("unexpected max past notes: %d", s.GetMaxPastNotesPerChannel())
	}
	if got := s.GetOctaveShift(); got != 1 {
		t.Fatalf("expected octave shift 1, got %d", got)
	}
}

func TestS3MSemitonePeriods(t *testing.T) {
	s := S3MSystem
	period, ok := s.GetSemitonePeriod(note.KeyC)
	if !ok {
		t.Fatalf("expected semitone period for KeyC")
	}
	if period != semitonePeriodTable[0]>>s.GetOctaveShift() {
		t.Fatalf("unexpected semitone period for KeyC: %d", period)
	}
	if _, ok := s.GetSemitonePeriod(note.KeyInvalid1); ok {
		t.Fatalf("expected invalid key to report missing semitone period")
	}
}
