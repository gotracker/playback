// Package s3m does a thing.
package s3m

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/settings"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// S3M is the exported interface to the S3M file loader
	S3M = format{}
)

// Load loads an S3M file into a playback system
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, s)
}

// Load loads an S3M file on a reader into a playback system
func (f format) LoadFromReader(r io.Reader, s *settings.Settings) (playback.Playback, error) {
	return load.S3M(r, s)
}
