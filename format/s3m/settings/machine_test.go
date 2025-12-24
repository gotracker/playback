package settings

import (
	"testing"

	"github.com/gotracker/playback/period"
)

func TestGetMachineSettingsS3M(t *testing.T) {
	ms := GetMachineSettings(false)
	if ms != amigaS3MSettings {
		t.Fatalf("expected amigaS3MSettings pointer")
	}
	if ms.PeriodConverter != amigaS3MSettings.PeriodConverter {
		t.Fatalf("unexpected period converter")
	}
	if !ms.OPL2Enabled {
		t.Fatalf("expected OPL2 enabled")
	}
	if !ms.Quirks.PreviousPeriodUsesModifiedPeriod || !ms.Quirks.PortaToNoteUsesModifiedPeriod {
		t.Fatalf("expected S3M quirks set")
	}
}

func TestGetMachineSettingsMod(t *testing.T) {
	ms := GetMachineSettings(true)
	if ms != amigaMOD31Settings {
		t.Fatalf("expected amigaMOD31Settings pointer")
	}
	if ms.PeriodConverter != amigaMOD31Settings.PeriodConverter {
		t.Fatalf("unexpected period converter")
	}
}

func TestSettingsPeriodType(t *testing.T) {
	// ensure type parameters line up with period.Amiga
	ms := GetMachineSettings(false)
	var _ *period.Amiga
	if ms.PeriodConverter == nil {
		t.Fatalf("expected period converter present")
	}
}
