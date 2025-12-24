package layout

import (
	"testing"

	"github.com/gotracker/playback/filter"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
)

func TestChannelSettingGetters(t *testing.T) {
	cs := ChannelSetting{
		Enabled:          true,
		Muted:            false,
		OutputChannelNum: 2,
		InitialVolume:    itVolume.Volume(10),
		ChannelVolume:    itVolume.FineVolume(20),
		PanEnabled:       true,
		InitialPanning:   itPanning.Panning(0x10),
		Vol0OptEnabled:   true,
	}

	if !cs.IsEnabled() || cs.IsMuted() {
		t.Fatalf("expected enabled and unmuted")
	}
	if got := cs.GetOutputChannelNum(); got != 2 {
		t.Fatalf("unexpected output channel: %d", got)
	}
	if got := cs.GetInitialVolume(); got != 10 {
		t.Fatalf("unexpected initial volume: %d", got)
	}
	if got := cs.GetMixingVolume(); got != 20 {
		t.Fatalf("unexpected mixing volume: %d", got)
	}
	if got := cs.GetInitialPanning(); got != 0x10 {
		t.Fatalf("unexpected initial panning: %d", got)
	}
	if !cs.IsPanEnabled() {
		t.Fatalf("expected pan enabled")
	}
	vo := cs.GetVol0OptimizationSettings()
	if !vo.Enabled || vo.MaxRowsAt0 != 3 {
		t.Fatalf("unexpected vol0 optimization settings: %+v", vo)
	}
}

func TestChannelSettingDefaultPanningWhenDisabled(t *testing.T) {
	cs := ChannelSetting{PanEnabled: false, InitialPanning: itPanning.Panning(0x40)}
	if got := cs.GetInitialPanning(); got != itPanning.DefaultPanning {
		t.Fatalf("expected default panning when disabled, got %d", got)
	}
}

func TestChannelSettingFilterDefaults(t *testing.T) {
	cs := ChannelSetting{}
	if got := cs.GetDefaultFilterInfo(); got != (filter.Info{}) {
		t.Fatalf("expected empty filter info, got %+v", got)
	}
	if cs.IsDefaultFilterEnabled() {
		t.Fatalf("expected default filter disabled")
	}
}
