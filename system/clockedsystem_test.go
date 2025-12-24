package system

import (
	"testing"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
)

func TestClockedSystemAccessors(t *testing.T) {
	s := ClockedSystem{
		MaxPastNotesPerChannel: 2,
		BaseClock:              1000,
		BaseFinetunes:          48,
		FinetunesPerOctave:     192,
		FinetunesPerNote:       16,
		CommonPeriod:           900,
		CommonRate:             frequency.Frequency(48000),
		SemitonePeriods:        [note.NumKeys]uint16{10, 20, 30, 40},
		OctaveShift:            1,
	}

	if s.GetMaxPastNotesPerChannel() != 2 {
		t.Fatalf("MaxPastNotesPerChannel = %d", s.GetMaxPastNotesPerChannel())
	}
	if s.GetBaseClock() != 1000 {
		t.Fatalf("BaseClock = %v", s.GetBaseClock())
	}
	if s.GetBaseFinetunes() != 48 {
		t.Fatalf("BaseFinetunes = %v", s.GetBaseFinetunes())
	}
	if s.GetFinetunesPerOctave() != 192 || s.GetFinetunesPerSemitone() != 16 {
		t.Fatalf("unexpected finetune values")
	}
	if s.GetCommonPeriod() != 900 || s.GetCommonRate() != 48000 {
		t.Fatalf("unexpected common values")
	}

	per0, ok := s.GetSemitonePeriod(note.Key(0))
	if !ok || per0 != 10>>1 {
		t.Fatalf("expected semitone period 5, got %d ok=%v", per0, ok)
	}
	per1, ok := s.GetSemitonePeriod(note.Key(1))
	if !ok || per1 != 20>>1 {
		t.Fatalf("expected semitone period 10, got %d ok=%v", per1, ok)
	}
	if _, ok := s.GetSemitonePeriod(note.Key(20)); ok {
		t.Fatalf("expected out-of-range key to return ok=false")
	}
	if s.GetOctaveShift() != 1 {
		t.Fatalf("OctaveShift = %d", s.GetOctaveShift())
	}
}
