package mod

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// MOD is the exported interface to the MOD file loader
	MOD = format{}
)

func loadSong(r io.Reader, features []feature.Feature) (*s3m.Song, error) {
	l, err := load.MOD(r, features)
	if err != nil {
		return nil, err
	}

	s := s3m.Song{
		Layout: l,
	}

	return &s, nil
}

// Load loads a MOD file into a playback system
func (f format) Load(filename string, features []feature.Feature) (playback.Song, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return loadSong(r, features)
}

// LoadFromReader loads a MOD file on a reader into a playback system
func (f format) LoadFromReader(r io.Reader, features []feature.Feature) (playback.Song, error) {
	return loadSong(r, features)
}
