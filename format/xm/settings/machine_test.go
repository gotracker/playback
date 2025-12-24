package settings

import (
    "testing"

    xmPeriod "github.com/gotracker/playback/format/xm/period"
    "github.com/gotracker/playback/period"
)

func TestGetMachineSettingsAmiga(t *testing.T) {
    ms := GetMachineSettings[period.Amiga]()
    if ms.PeriodConverter != xmPeriod.AmigaConverter {
        t.Fatalf("expected Amiga converter")
    }
    if ms.OPL2Enabled {
        t.Fatalf("expected OPL2 disabled")
    }
    if _, err := ms.GetFilterFactory("", 0, nil); err != nil {
        t.Fatalf("expected empty filter ok: %v", err)
    }
}

func TestGetMachineSettingsLinear(t *testing.T) {
    ms := GetMachineSettings[period.Linear]()
    if ms.PeriodConverter != xmPeriod.LinearConverter {
        t.Fatalf("expected Linear converter")
    }
}
