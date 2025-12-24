package layout

import (
	"testing"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
)

func TestChannelSettingGetters(t *testing.T) {
	cs := ChannelSetting{
		Enabled:          true,
		Muted:            false,
		OutputChannelNum: 3,
		Category:         s3mfile.ChannelCategoryOPL2Melody,
		InitialVolume:    s3mVolume.Volume(12),
		PanEnabled:       true,
		InitialPanning:   s3mPanning.Panning(0x0a),
		Memory:           channel.Memory{Shared: &channel.SharedMemory{ZeroVolOptimization: true}},
		DefaultFilter:    filter.Info{Name: "amigalpf"},
	}

	if !cs.IsEnabled() || cs.IsMuted() {
		t.Fatalf("expected enabled and unmuted")
	}
	if got := cs.GetOutputChannelNum(); got != 3 {
		t.Fatalf("unexpected output channel: %d", got)
	}
	if got := cs.GetInitialVolume(); got != 12 {
		t.Fatalf("unexpected initial volume: %d", got)
	}
	if got := cs.GetMixingVolume(); got != s3mVolume.FineVolume(0x7f) {
		t.Fatalf("unexpected mixing volume: %d", got)
	}
	if got := cs.GetInitialPanning(); got != 0x0a {
		t.Fatalf("unexpected initial panning: %d", got)
	}
	if !cs.IsPanEnabled() {
		t.Fatalf("expected pan enabled")
	}
	if !cs.IsDefaultFilterEnabled() {
		t.Fatalf("expected default filter enabled")
	}
	vo := cs.GetVol0OptimizationSettings()
	if !vo.Enabled || vo.MaxRowsAt0 != 3 {
		t.Fatalf("unexpected vol0 optimization settings: %+v", vo)
	}
	if ch := cs.GetOPLChannel(); ch != 3 {
		t.Fatalf("expected OPL channel passthrough, got %d", ch)
	}
}

func TestChannelSettingDefaults(t *testing.T) {
	cs := ChannelSetting{PanEnabled: false, InitialPanning: s3mPanning.Panning(0x02)}
	if got := cs.GetInitialPanning(); got != s3mPanning.DefaultPanning {
		t.Fatalf("expected default panning when disabled, got %d", got)
	}
	if got := cs.GetDefaultFilterInfo(); got != (filter.Info{}) {
		t.Fatalf("expected empty filter info, got %+v", got)
	}
	if cs.IsDefaultFilterEnabled() {
		t.Fatalf("expected default filter disabled")
	}
	if ch := cs.GetOPLChannel(); ch != index.InvalidOPLChannel {
		t.Fatalf("expected invalid OPL channel for non-OPL category, got %d", ch)
	}
}
