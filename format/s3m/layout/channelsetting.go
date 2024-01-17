package layout

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/vol0optimization"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled           bool
	Muted             bool
	OutputChannelNum  int
	Category          s3mfile.ChannelCategory
	InitialVolume     s3mVolume.Volume
	PanEnabled        bool
	InitialPanning    s3mPanning.Panning
	Memory            channel.Memory
	DefaultFilterName string
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

func (c ChannelSetting) GetInitialVolume() s3mVolume.Volume {
	return c.InitialVolume
}

func (c ChannelSetting) GetMixingVolume() s3mVolume.FineVolume {
	return s3mVolume.FineVolume(0x7f)
}

func (c ChannelSetting) GetInitialPanning() s3mPanning.Panning {
	if c.PanEnabled {
		return c.InitialPanning
	}
	return s3mPanning.DefaultPanning
}

func (c ChannelSetting) GetMemory() song.ChannelMemory {
	return &c.Memory
}

func (c ChannelSetting) IsPanEnabled() bool {
	return c.PanEnabled
}

func (c ChannelSetting) GetDefaultFilterName() string {
	return c.DefaultFilterName
}

func (c ChannelSetting) IsDefaultFilterEnabled() bool {
	return len(c.DefaultFilterName) > 0
}

func (c ChannelSetting) GetVol0OptimizationSettings() vol0optimization.Vol0OptimizationSettings {
	return vol0optimization.Vol0OptimizationSettings{
		Enabled:    c.Memory.Shared.ZeroVolOptimization,
		MaxRowsAt0: 3,
	}
}

func (c ChannelSetting) GetOPLChannel() index.OPLChannel {
	switch c.Category {
	case s3mfile.ChannelCategoryOPL2Melody, s3mfile.ChannelCategoryOPL2Drums:
		return index.OPLChannel(c.OutputChannelNum)
	default:
		return index.InvalidOPLChannel
	}
}
