package memory

import "testing"

func TestValueCoalesceAndXY(t *testing.T) {
	var v Value[uint8]

	if got := v.Coalesce(0); got != 0 {
		t.Fatalf("expected zero fallback when empty, got %d", got)
	}

	if got := v.Coalesce(0x3c); got != 0x3c {
		t.Fatalf("expected stored value, got %d", got)
	}
	if got := v.Coalesce(0); got != 0x3c {
		t.Fatalf("expected coalesce to reuse memory value, got %d", got)
	}

	v.Reset()
	hi, lo := v.CoalesceXY(0xAB)
	if hi != 0x0a || lo != 0x0b {
		t.Fatalf("unexpected XY split: hi=%X lo=%X", hi, lo)
	}

	hi, lo = v.CoalesceXY(0)
	if hi != 0x0a || lo != 0x0b {
		t.Fatalf("expected XY split to persist stored value: hi=%X lo=%X", hi, lo)
	}
}
