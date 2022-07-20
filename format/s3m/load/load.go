package load

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/layout"
	"github.com/gotracker/playback/format/s3m/load/modconv"
	"github.com/gotracker/playback/settings"
	formatutil "github.com/gotracker/playback/util"
)

func readMOD(filename string, s *settings.Settings) (*layout.Song, error) {
	buffer, err := formatutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	f, err := modconv.Read(buffer)
	if err != nil {
		return nil, err
	}

	return convertS3MFileToSong(f, func(patNum int) uint8 {
		return 64
	}, s, true)
}

// MOD loads a MOD file and upgrades it into an S3M file internally
func MOD(filename string, s *settings.Settings) (playback.Playback, error) {
	return load(filename, readMOD, s)
}

// S3M loads an S3M file into a new Playback object
func S3M(filename string, s *settings.Settings) (playback.Playback, error) {
	return load(filename, readS3M, s)
}
