package effect

import (
	"fmt"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"

	"github.com/gotracker/playback"
	s3mVolume "github.com/gotracker/playback/format/s3m/conversion/volume"
	"github.com/gotracker/playback/format/s3m/layout/channel"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume ChannelCommand // 'V'

// PreStart triggers when the effect enters onto the channel state
func (e SetGlobalVolume) PreStart(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	p.SetGlobalVolume(s3mVolume.VolumeFromS3M(s3mfile.Volume(channel.DataEffect(e))))
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetGlobalVolume) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetGlobalVolume) String() string {
	return fmt.Sprintf("V%0.2x", channel.DataEffect(e))
}
