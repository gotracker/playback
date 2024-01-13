package layout

import (
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/song"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	OutputChannelNum int
	InitialVolume    xmVolume.XmVolume
	InitialPanning   xmPanning.Panning
	Memory           channel.Memory
}

var _ song.ChannelSettings = (*ChannelSetting)(nil)

func (c ChannelSetting) GetEnabled() bool {
	return c.Enabled
}

func (c ChannelSetting) GetOutputChannelNum() int {
	return c.OutputChannelNum
}

func (c ChannelSetting) GetInitialVolume() xmVolume.XmVolume {
	return c.InitialVolume
}

func (c ChannelSetting) GetMixingVolume() xmVolume.XmVolume {
	return xmVolume.DefaultXmVolume
}

func (c ChannelSetting) GetInitialPanning() xmPanning.Panning {
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
