package feature

import (
	"testing"

	"github.com/gotracker/playback/voice/pcm"
)

func TestUseNativeSampleFormat(t *testing.T) {
	pc := UseNativeSampleFormat(true)
	cfg, ok := pc.(PreConvertSamples)
	if !ok {
		t.Fatalf("expected PreConvertSamples, got %T", pc)
	}
	if !cfg.Enabled {
		t.Fatalf("expected Enabled=true")
	}
	if cfg.DesiredFormat != pcm.SampleDataFormatNative {
		t.Fatalf("unexpected DesiredFormat: %v", cfg.DesiredFormat)
	}
}
