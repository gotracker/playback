package layout

import (
	"testing"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
)

func TestChannelSettingDefaults(t *testing.T) {
	cs := ChannelSetting{
		Enabled:          true,
		Muted:            false,
		OutputChannelNum: 2,
		InitialVolume:    xmVolume.XmVolume(0x20),
		InitialPanning:   xmPanning.DefaultPanning,
	}

	if !cs.IsEnabled() || cs.IsMuted() {
		t.Fatalf("expected enabled and not muted")
	}
	if cs.GetOutputChannelNum() != 2 {
		t.Fatalf("expected output channel 2, got %d", cs.GetOutputChannelNum())
	}
	if cs.GetInitialVolume() != cs.InitialVolume {
		t.Fatalf("initial volume mismatch")
	}
	if cs.GetInitialPanning() != cs.InitialPanning {
		t.Fatalf("initial panning mismatch")
	}
	if !cs.IsPanEnabled() {
		t.Fatalf("expected pan enabled")
	}
	if cs.GetOPLChannel() != index.InvalidOPLChannel {
		t.Fatalf("expected invalid OPL channel")
	}
}
