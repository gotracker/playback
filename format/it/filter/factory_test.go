package filter

import (
	"testing"

	pf "github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
)

func TestFactoryAmigaLPF(t *testing.T) {
	f, err := Factory("amigalpf", frequency.Frequency(8363), nil)
	if err != nil {
		t.Fatalf("expected nil error for amigalpf: %v", err)
	}
	if f == nil {
		t.Fatalf("expected non-nil filter for amigalpf")
	}
}

func TestFactoryITResonantParamsValidation(t *testing.T) {
	params := pf.ITResonantFilterParams{Cutoff: 0x40, Resonance: 0x20, ExtendedFilterRange: true, Highpass: false}
	f, err := Factory("itresonant", frequency.Frequency(44100), params)
	if err != nil {
		t.Fatalf("expected nil error for itresonant: %v", err)
	}
	if f == nil {
		t.Fatalf("expected filter instance for itresonant")
	}

	if _, err := Factory("itresonant", frequency.Frequency(44100), "bad"); err == nil {
		t.Fatalf("expected type assertion error for wrong params type")
	}
}

func TestFactoryEchoParamsValidation(t *testing.T) {
	params := pf.EchoFilterSettings{WetDryMix: 0.5, Feedback: 0.3, LeftDelay: 0.05, RightDelay: 0.06, PanDelay: 0.7}
	f, err := Factory("echo", frequency.Frequency(48000), params)
	if err != nil {
		t.Fatalf("expected nil error for echo: %v", err)
	}
	if f == nil {
		t.Fatalf("expected filter instance for echo")
	}

	if _, err := Factory("echo", frequency.Frequency(48000), 123); err == nil {
		t.Fatalf("expected type assertion error for wrong echo params")
	}
}

func TestFactoryEmptyAndUnknown(t *testing.T) {
	if f, err := Factory("", frequency.Frequency(0), nil); err != nil {
		t.Fatalf("expected empty name to succeed: %v", err)
	} else if f != nil {
		t.Fatalf("expected nil filter for empty name")
	}

	if _, err := Factory("nope", frequency.Frequency(0), nil); err == nil {
		t.Fatalf("expected error for unsupported filter")
	}
}
