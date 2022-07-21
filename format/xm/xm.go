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

// Load loads an XM file into a playback system
func (f format) Load(filename string, features []feature.Feature) (playback.Playback, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, features)
}

// LoadFromReader loads an XM file on a reader into a playback system
func (f format) LoadFromReader(r io.Reader, features []feature.Feature) (playback.Playback, error) {
	return load.XM(r, features)
}
