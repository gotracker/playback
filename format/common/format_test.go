package common

import (
	"testing"

	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine/settings"
	optional "github.com/heucuva/optional"
)

func TestConvertFeaturesToSettings(t *testing.T) {
	us := settings.UserSettings{}

	features := []feature.Feature{
		feature.SongLoop{Count: 2},
		feature.StartOrderAndRow{},
		feature.PlayUntilOrderAndRow{Order: 3, Row: 5},
		feature.SetDefaultTempo{Tempo: 6},
		feature.SetDefaultBPM{BPM: 7},
		feature.IgnoreUnknownEffect{Enabled: true},
	}

	if err := (Format{}).ConvertFeaturesToSettings(&us, features); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if us.SongLoopCount != 2 {
		t.Fatalf("expected SongLoopCount=2, got %d", us.SongLoopCount)
	}
	if o, set := us.Start.Order.Get(); set {
		t.Fatalf("expected Start.Order to remain unset, got %v", o)
	}
	if r, set := us.Start.Row.Get(); set {
		t.Fatalf("expected Start.Row to remain unset, got %v", r)
	}
	if o, set := us.PlayUntil.Order.Get(); !set || o != 3 {
		t.Fatalf("expected PlayUntil.Order=3, got %v set=%v", o, set)
	}
	if r, set := us.PlayUntil.Row.Get(); !set || r != 5 {
		t.Fatalf("expected PlayUntil.Row=5, got %v set=%v", r, set)
	}
	if us.Start.Tempo != 6 {
		t.Fatalf("expected Start.Tempo=6, got %d", us.Start.Tempo)
	}
	if us.Start.BPM != 7 {
		t.Fatalf("expected Start.BPM=7, got %d", us.Start.BPM)
	}
	if !us.IgnoreUnknownEffect {
		t.Fatalf("expected IgnoreUnknownEffect=true")
	}
}

func TestConvertFeaturesSetsStartOrderAndRow(t *testing.T) {
	us := settings.UserSettings{}

	features := []feature.Feature{
		feature.StartOrderAndRow{
			Order: optional.NewValue(4),
			Row:   optional.NewValue(2),
		},
	}

	if err := (Format{}).ConvertFeaturesToSettings(&us, features); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o, set := us.Start.Order.Get(); !set || o != 4 {
		t.Fatalf("expected Start.Order=4 set, got %v set=%v", o, set)
	}
	if r, set := us.Start.Row.Get(); !set || r != 2 {
		t.Fatalf("expected Start.Row=2 set, got %v set=%v", r, set)
	}
}

type unknownFeature struct{}

func TestConvertFeaturesIgnoresUnknown(t *testing.T) {
	us := settings.UserSettings{}
	us.SongLoopCount = 5

	features := []feature.Feature{unknownFeature{}}

	if err := (Format{}).ConvertFeaturesToSettings(&us, features); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if us.SongLoopCount != 5 {
		t.Fatalf("expected SongLoopCount to remain 5, got %d", us.SongLoopCount)
	}
}

func TestConvertFeaturesSetsQuirksMode(t *testing.T) {
	us := settings.UserSettings{}

	features := []feature.Feature{
		feature.QuirksMode{
			Profile:      optional.NewValue("ft2.09"),
			LinearSlides: optional.NewValue(true),
		},
	}

	if err := (Format{}).ConvertFeaturesToSettings(&us, features); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prof, ok := us.Quirks.Profile.Get(); !ok || prof != "ft2.09" {
		t.Fatalf("expected quirks profile ft2.09 set, got %v set=%v", prof, ok)
	}
	if linear, ok := us.Quirks.LinearSlidesOverride.Get(); !ok || !linear {
		t.Fatalf("expected linear slides override set true, got %v set=%v", linear, ok)
	}
}

func TestConvertFeaturesLeavesDefaultsOnEmpty(t *testing.T) {
	us := settings.UserSettings{}
	us.SongLoopCount = 4
	us.Start.Tempo = 3
	us.Start.BPM = 2
	us.IgnoreUnknownEffect = true

	if err := (Format{}).ConvertFeaturesToSettings(&us, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if us.SongLoopCount != 4 {
		t.Fatalf("expected SongLoopCount to stay 4, got %d", us.SongLoopCount)
	}
	if us.Start.Tempo != 3 || us.Start.BPM != 2 {
		t.Fatalf("expected Start tempo/BPM unchanged, got tempo=%d bpm=%d", us.Start.Tempo, us.Start.BPM)
	}
	if o, set := us.Start.Order.Get(); set {
		t.Fatalf("expected Start.Order to remain unset, got %v", o)
	}
	if r, set := us.Start.Row.Get(); set {
		t.Fatalf("expected Start.Row to remain unset, got %v", r)
	}
	if !us.IgnoreUnknownEffect {
		t.Fatalf("expected IgnoreUnknownEffect to remain true")
	}
}
