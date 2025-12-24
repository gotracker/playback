package filter

import (
	"testing"

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
