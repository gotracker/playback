package load

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/s3m/layout"
	"github.com/gotracker/playback/format/s3m/load/modconv"
	s3mPlayback "github.com/gotracker/playback/format/s3m/playback"
	"github.com/gotracker/playback/settings"
)

func readMOD(r io.Reader, s *settings.Settings) (*layout.Song, error) {
	f, err := modconv.Read(r)
	if err != nil {
		return nil, err
	}

	return convertS3MFileToSong(f, func(patNum int) uint8 {
		return 64
	}, s, true)
}

// MOD loads a MOD file and upgrades it into an S3M file internally
func MOD(r io.Reader, s *settings.Settings) (playback.Playback, error) {
	return common.Load(r, readMOD, s3mPlayback.NewManager, s)
}

// S3M loads an S3M file into a new Playback object
func S3M(r io.Reader, s *settings.Settings) (playback.Playback, error) {
	return common.Load(r, readS3M, s3mPlayback.NewManager, s)
}
