package layout

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/format/it/channel"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	OutputChannelNum int
	InitialVolume    volume.Volume
	ChannelVolume    volume.Volume
	InitialPanning   panning.Position
	Memory           channel.Memory
}
