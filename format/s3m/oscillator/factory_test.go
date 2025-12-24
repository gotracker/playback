package oscillator

import "testing"

func TestOscillatorFactoryKnown(t *testing.T) {
	for _, name := range []string{"vibrato", "tremolo", "panbrello"} {
		osc, err := OscillatorFactory(name)
		if err != nil {
			t.Fatalf("expected nil error for %s: %v", name, err)
		}
		if osc == nil {
			t.Fatalf("expected oscillator for %s", name)
		}
	}
}

func TestOscillatorFactoryEmpty(t *testing.T) {
	osc, err := OscillatorFactory("")
	if err != nil {
		t.Fatalf("expected nil error for empty name: %v", err)
	}
	if osc != nil {
		t.Fatalf("expected nil oscillator for empty name")
	}
}

func TestOscillatorFactoryUnknown(t *testing.T) {
	if _, err := OscillatorFactory("nope"); err == nil {
		t.Fatalf("expected error for unsupported oscillator")
	}
}
