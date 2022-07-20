package layout

import (
	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/format/s3m/channel"
)

// ChannelSetting is settings specific to a single channel
type ChannelSetting struct {
	Enabled          bool
	OutputChannelNum int
	Category         s3mfile.ChannelCategory
	InitialVolume    volume.Volume
	InitialPanning   panning.Position
	Memory           channel.Memory
}
