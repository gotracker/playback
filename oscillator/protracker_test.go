package oscillator

import "testing"

func TestGetProtrackerSineWrapsAndSigns(t *testing.T) {
	if v := GetProtrackerSine(0); v != 0 {
		t.Fatalf("expected sine at 0 to be 0, got %v", v)
	}

	posPos := GetProtrackerSine(8) // table value 180/255
	posNeg := GetProtrackerSine(40)
	if posPos <= 0 {
		t.Fatalf("expected positive sine at pos 8, got %v", posPos)
	}
	if posNeg >= 0 {
		t.Fatalf("expected negative sine after wrap, got %v", posNeg)
	}
	if diff := posPos - (-posNeg); diff > 1e-4 {
		t.Fatalf("expected magnitudes to mirror within tolerance, diff=%v", diff)
	}
}

func TestProtrackerOscillatorAdvanceAndReset(t *testing.T) {
	osc := NewProtrackerOscillator().(*protrackerOscillator)
	osc.SetWaveform(WaveTableSelectSineRetrigger)

	osc.Pos = 60
	osc.Advance(5) // wraps past 63
	if osc.Pos != 1 {
		t.Fatalf("expected position to wrap to 1, got %d", osc.Pos)
	}

	osc.HardReset()
	if osc.Pos != 0 {
		t.Fatalf("expected hard reset to zero position, got %d", osc.Pos)
	}

	osc.Advance(16)
	wave := osc.GetWave(0.5)
	expected := GetProtrackerSine(16) * 0.5
	if wave != expected {
		t.Fatalf("unexpected wave output: got %v want %v", wave, expected)
	}
}
