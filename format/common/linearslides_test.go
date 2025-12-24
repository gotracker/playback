package common

import (
	"testing"

	"github.com/gotracker/playback/player/feature"
	"github.com/heucuva/optional"
)

func TestResolveLinearSlidesUsesOverride(t *testing.T) {
	features := []feature.Feature{
		feature.QuirksMode{LinearSlides: optional.NewValue(false)},
	}
	if linear := ResolveLinearSlides(true, features); linear {
		t.Fatalf("expected override to disable linear slides")
	}
}

func TestResolveLinearSlidesFallsBackToDefault(t *testing.T) {
	if linear := ResolveLinearSlides(false, nil); linear {
		t.Fatalf("expected default value to be used when no override present")
	}
}
