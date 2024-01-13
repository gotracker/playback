package layout

import (
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/song"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	OutputChannelNum int
	InitialVolume    itVolume.Volume
	ChannelVolume    itVolume.FineVolume
	InitialPanning   itPanning.Panning
	Memory           channel.Memory
}

var _ song.ChannelSettings = (*ChannelSetting)(nil)

func (c ChannelSetting) GetEnabled() bool {
	return c.Enabled
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
	return c.InitialPanning
}

func (c ChannelSetting) GetMemory() song.ChannelMemory {
	return &c.Memory
}

func (c ChannelSetting) GetPanEnabled() bool {
	return true
}

func (c ChannelSetting) GetDefaultFilterName() string {
	return ""
}

func (c ChannelSetting) IsDefaultFilterEnabled() bool {
	return false
}
