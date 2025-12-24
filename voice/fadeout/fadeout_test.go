package fadeout

import (
	"testing"

	"github.com/gotracker/playback/mixing/volume"
)

func TestModeIsFadeoutActive(t *testing.T) {
	cases := []struct {
		mode          Mode
		force         bool
		volEnvEnabled bool
		volEnvDone    bool
		expect        bool
	}{
		{ModeDisabled, false, false, false, false},
		{ModeAlwaysActive, false, false, false, true},
		{ModeAlwaysActive, true, false, false, true},
		{ModeAlwaysActive, false, true, false, false},
		{ModeAlwaysActive, false, true, true, true},
		{ModeOnlyIfVolEnvActive, false, false, false, false},
		{ModeOnlyIfVolEnvActive, true, false, false, true},
		{ModeOnlyIfVolEnvActive, false, true, false, true},
	}

	for _, tt := range cases {
		if got := tt.mode.IsFadeoutActive(tt.force, tt.volEnvEnabled, tt.volEnvDone); got != tt.expect {
			t.Fatalf("mode %v force=%v env=%v done=%v => %v, want %v", tt.mode, tt.force, tt.volEnvEnabled, tt.volEnvDone, got, tt.expect)
		}
	}
}

func TestSettingsHoldsValues(t *testing.T) {
	s := Settings{Mode: ModeAlwaysActive, Amount: volume.Volume(0.5)}
	if s.Mode != ModeAlwaysActive {
		t.Fatalf("expected ModeAlwaysActive")
	}
	if s.Amount != 0.5 {
		t.Fatalf("expected Amount=0.5, got %v", s.Amount)
	}
}
