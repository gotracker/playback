package it

import (
	"bytes"
	"testing"

	itFeature "github.com/gotracker/playback/format/it/feature"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/feature"
)

func TestLoadFromReaderRejectsInvalidData(t *testing.T) {
	_, err := IT.LoadFromReader(bytes.NewReader([]byte("bad")), nil)
	if err == nil {
		t.Fatalf("expected error for invalid IT data")
	}
}

func TestConvertFeaturesToSettings(t *testing.T) {
	var us settings.UserSettings
	features := []feature.Feature{
		itFeature.LongChannelOutput{Enabled: true},
		itFeature.NewNoteActions{Enabled: false},
		feature.SongLoop{Count: 3},
	}

	if err := IT.ConvertFeaturesToSettings(&us, features); err != nil {
		t.Fatalf("ConvertFeaturesToSettings error: %v", err)
	}
	if !us.LongChannelOutput {
		t.Fatalf("expected LongChannelOutput set")
	}
	if us.EnableNewNoteActions {
		t.Fatalf("expected EnableNewNoteActions cleared")
	}
	if us.SongLoopCount != 3 {
		t.Fatalf("expected SongLoopCount=3, got %d", us.SongLoopCount)
	}
}
