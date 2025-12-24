package pitchpan

import "testing"

func TestPitchPanDefaults(t *testing.T) {
	var p PitchPan
	if p.Enabled {
		t.Fatalf("expected default Enabled=false")
	}
	if p.Center != 0 {
		t.Fatalf("expected default Center=0, got %d", p.Center)
	}
	if p.Separation != 0 {
		t.Fatalf("expected default Separation=0, got %v", p.Separation)
	}
}

func TestPitchPanValues(t *testing.T) {
	p := PitchPan{Enabled: true, Center: 5, Separation: 0.25}
	if !p.Enabled {
		t.Fatalf("expected Enabled=true")
	}
	if p.Center != 5 {
		t.Fatalf("expected Center=5, got %d", p.Center)
	}
	if p.Separation != 0.25 {
		t.Fatalf("expected Separation=0.25, got %v", p.Separation)
	}
}
