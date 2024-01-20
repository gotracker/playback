package layout

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/vol0optimization"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	Muted            bool
	OutputChannelNum int
	InitialVolume    xmVolume.XmVolume
	InitialPanning   xmPanning.Panning
	Memory           channel.Memory
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

func (c ChannelSetting) IsPanEnabled() bool {
	return true
}

func (c ChannelSetting) GetDefaultFilterInfo() filter.Info {
	return filter.Info{}
}

func (c ChannelSetting) IsDefaultFilterEnabled() bool {
	return false
}

func (c ChannelSetting) GetVol0OptimizationSettings() vol0optimization.Vol0OptimizationSettings {
	return vol0optimization.Vol0OptimizationSettings{
		Enabled:    true,
		MaxRowsAt0: 3,
	}
}

func (ChannelSetting) GetOPLChannel() index.OPLChannel {
	return index.InvalidOPLChannel
}
