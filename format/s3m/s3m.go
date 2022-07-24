package s3m

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// S3M is the exported interface to the S3M file loader
	S3M = format{}
)

func loadSong(r io.Reader, features []feature.Feature) (*Song, error) {
	l, err := load.S3M(r, features)
	if err != nil {
		return nil, err
	}

	s := Song{
		Layout: *l,
	}

	return &s, nil
}

// Load loads an IT file into a playback system
func (f format) Load(filename string, features []feature.Feature) (playback.Song, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return loadSong(r, features)
}

func (f format) LoadFromReader(r io.Reader, features []feature.Feature) (playback.Song, error) {
	return loadSong(r, features)
}
