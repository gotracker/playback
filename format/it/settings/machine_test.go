package settings

import (
	"testing"

	"github.com/gotracker/playback/period"
)

func TestGetMachineSettingsAmiga(t *testing.T) {
	ms := GetMachineSettings[period.Amiga]()
	if ms != amigaMachine {
		t.Fatalf("expected amiga machine settings pointer")
	}
	if ms.PeriodConverter != amigaMachine.PeriodConverter {
		t.Fatalf("expected amiga period converter")
	}
	if ms.VoiceFactory != amigaMachine.VoiceFactory {
		t.Fatalf("expected amiga voice factory")
	}
	if ms.OPL2Enabled {
		t.Fatalf("unexpected OPL2 enabled for amiga machine")
	}
}

func TestGetMachineSettingsLinear(t *testing.T) {
	ms := GetMachineSettings[period.Linear]()
	if ms != linearMachine {
		t.Fatalf("expected linear machine settings pointer")
	}
	if ms.PeriodConverter != linearMachine.PeriodConverter {
		t.Fatalf("expected linear period converter")
	}
	if ms.VoiceFactory != linearMachine.VoiceFactory {
		t.Fatalf("expected linear voice factory")
	}
	if ms.OPL2Enabled {
		t.Fatalf("unexpected OPL2 enabled for linear machine")
	}
}

type unsupportedPeriod struct{}

func (unsupportedPeriod) IsInvalid() bool { return false }

func TestGetMachineSettingsUnsupportedPanics(t *testing.T) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unsupported period type")
		}
	}()

	_ = GetMachineSettings[unsupportedPeriod]()
}
