package channel

import (
	"fmt"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"

	"github.com/gotracker/playback"
)

// SetVolume defines a SetVolume effect
type SetVolume s3mfile.Volume

// Start triggers on the first tick, but before the Tick() function is called
func (e SetVolume) Start(cs S3MChannel, p playback.Playback) error {
	cs.SetActiveVolume(s3mVolume.VolumeFromS3M(s3mfile.Volume(e)))
	return nil
}

func (e SetVolume) String() string {
	return fmt.Sprintf("%02d", s3mfile.Volume(e))
}
