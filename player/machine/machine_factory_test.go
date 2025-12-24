package machine

import (
	"testing"

	"github.com/gotracker/playback/player/machine/settings"
)

func TestNewMachineErrorsForUnregisteredTypes(t *testing.T) {
	if _, err := NewMachine(stubSongData{}, settings.UserSettings{}); err == nil {
		t.Fatalf("expected error for unregistered machine types")
	}
}

func TestRegisterMachinePanicsOnDuplicate(t *testing.T) {
	ms := &settings.MachineSettings[stubPeriod, stubGV, stubGV, stubGV, stubPan]{
		PeriodConverter: stubPeriodCalc{},
	}
	RegisterMachine(ms)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected duplicate registration to panic")
		}
	}()

	RegisterMachine(ms)
}
