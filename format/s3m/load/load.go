package load

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/s3m/layout"
	"github.com/gotracker/playback/format/s3m/load/modconv"
	"github.com/gotracker/playback/player/feature"
)

func readMOD(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	f, err := modconv.Read(r)
	if err != nil {
		return nil, err
	}

	return convertS3MFileToSong(f, func(patNum int) uint8 {
		return 64
	}, features, true)
}

// MOD loads a MOD file from a reader and upgrades it into an S3M file internally
func MOD(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	return common.Load(r, readMOD, features)
}

// S3M loads an S3M file from a reader
func S3M(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	return common.Load(r, readS3M, features)
}
