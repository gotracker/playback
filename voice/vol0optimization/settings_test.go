package vol0optimization

import "testing"

func TestVol0OptimizationSettingsDefaults(t *testing.T) {
	var s Vol0OptimizationSettings
	if s.Enabled {
		t.Fatalf("expected default Enabled=false")
	}
	if s.MaxRowsAt0 != 0 {
		t.Fatalf("expected default MaxRowsAt0=0, got %d", s.MaxRowsAt0)
	}
}

func TestVol0OptimizationSettingsValues(t *testing.T) {
	s := Vol0OptimizationSettings{Enabled: true, MaxRowsAt0: 4}
	if !s.Enabled {
		t.Fatalf("expected Enabled=true")
	}
	if s.MaxRowsAt0 != 4 {
		t.Fatalf("expected MaxRowsAt0=4, got %d", s.MaxRowsAt0)
	}
}
