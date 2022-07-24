// Package xm does a thing.
package xm

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/load"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// XM is the exported interface to the XM file loader
	XM = format{}
)

func loadSong(r io.Reader, features []feature.Feature) (*Song, error) {
	l, err := load.XM(r, features)
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
