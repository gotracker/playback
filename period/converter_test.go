package period

import (
	"testing"

	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
)

func TestLinearConverterSamplerAddAndPeriod(t *testing.T) {
	sys := system.ClockedSystem{
		BaseFinetunes:      0,
		FinetunesPerOctave: 12,
		FinetunesPerNote:   1,
	}
	conv := LinearConverter{System: sys}

	p := Linear{Finetune: 12}

	freq := conv.GetFrequency(p)
	if freq != 2 {
		t.Fatalf("expected frequency 2, got %v", freq)
	}

	add := conv.GetSamplerAdd(p, 10, 20)
	if add != 1 {
		t.Fatalf("expected sampler add 1, got %v", add)
	}

	n := note.Normal(note.NewSemitone(note.KeyC, 1))
	got := conv.GetPeriod(n)
	if got.Finetune != 12 {
		t.Fatalf("expected finetune 12 from note C1, got %d", got.Finetune)
	}

	empty := conv.GetPeriod(note.EmptyNote{})
	if empty.Finetune != 0 {
		t.Fatalf("expected empty note to yield zero period")
	}
}

func TestAmigaConverterSamplerAddAndPeriod(t *testing.T) {
	var semis [note.NumKeys]uint16
	semis[note.KeyC] = 1000

	sys := system.ClockedSystem{
		BaseClock:       1000,
		CommonRate:      10,
		SemitonePeriods: semis,
		OctaveShift:     0,
	}

	conv := AmigaConverter{System: sys, MinPeriod: 100, MaxPeriod: 900, SlideTo0Allowed: false}

	n := note.Normal(note.NewSemitone(note.KeyC, 1))
	p := conv.GetPeriod(n)
	if p != 500 {
		t.Fatalf("expected period 500, got %d", p)
	}

	freq := conv.GetFrequency(p)
	if freq != 2 {
		t.Fatalf("expected frequency 2, got %v", freq)
	}

	add := conv.GetSamplerAdd(p, 20, 40)
	if add != 0.1 {
		t.Fatalf("expected sampler add 0.1, got %v", add)
	}

	// invalid period should yield zero frequency/add
	zero := conv.GetFrequency(Amiga(0))
	if zero != 0 {
		t.Fatalf("expected zero frequency for invalid period")
	}
	if conv.GetSamplerAdd(Amiga(0), 20, 40) != 0 {
		t.Fatalf("expected zero sampler add for invalid period")
	}
}
