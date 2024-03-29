package layout

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/vol0optimization"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	Muted            bool
	OutputChannelNum int
	InitialVolume    itVolume.Volume
	ChannelVolume    itVolume.FineVolume
	PanEnabled       bool
	InitialPanning   itPanning.Panning
	Memory           channel.Memory
	Vol0OptEnabled   bool
}

var _ song.ChannelSettings = (*ChannelSetting)(nil)

func (c ChannelSetting) IsEnabled() bool {
	return c.Enabled
}

func (c ChannelSetting) IsMuted() bool {
	return c.Muted
}

func (c ChannelSetting) GetOutputChannelNum() int {
	return c.OutputChannelNum
}

func (c ChannelSetting) GetInitialVolume() itVolume.Volume {
	return c.InitialVolume
}

func (c ChannelSetting) GetMixingVolume() itVolume.FineVolume {
	return c.ChannelVolume
}

func (c ChannelSetting) GetInitialPanning() itPanning.Panning {
	if c.PanEnabled {
		return c.InitialPanning
	}
	return itPanning.DefaultPanning
}

func (c ChannelSetting) GetMemory() song.ChannelMemory {
	return &c.Memory
}

func (c ChannelSetting) IsPanEnabled() bool {
	return c.PanEnabled
}

func (c ChannelSetting) GetDefaultFilterInfo() filter.Info {
	return filter.Info{}
}

func (c ChannelSetting) IsDefaultFilterEnabled() bool {
	return false
}

func (c ChannelSetting) GetVol0OptimizationSettings() vol0optimization.Vol0OptimizationSettings {
	return vol0optimization.Vol0OptimizationSettings{
		Enabled:    c.Vol0OptEnabled,
		MaxRowsAt0: 3,
	}
}

func (ChannelSetting) GetOPLChannel() index.OPLChannel {
	return index.InvalidOPLChannel
}
