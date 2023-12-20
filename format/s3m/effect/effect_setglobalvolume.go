package effect

import (
	"fmt"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume ChannelCommand // 'V'

// PreStart triggers when the effect enters onto the channel state
func (e SetGlobalVolume) PreStart(cs S3MChannel, p playback.Playback) error {
	p.SetGlobalVolume(s3mVolume.VolumeFromS3M(s3mfile.Volume(channel.DataEffect(e))))
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetGlobalVolume) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetGlobalVolume) String() string {
	return fmt.Sprintf("V%0.2x", channel.DataEffect(e))
}
