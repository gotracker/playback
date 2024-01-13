// Package s3m does a thing.
package s3m

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/format/s3m/settings"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// S3M is the exported interface to the S3M file loader
	S3M = format{}
)

// Load loads an S3M file into a playback system
func (f format) Load(filename string, features []feature.Feature) (playback.Playback, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, features)
}

// Load loads an S3M file on a reader into a playback system
func (f format) LoadFromReader(r io.Reader, features []feature.Feature) (playback.Playback, error) {
	return load.S3M(r, features)
}

func init() {
	machine.RegisterMachine(settings.GetMachineSettings())
}
