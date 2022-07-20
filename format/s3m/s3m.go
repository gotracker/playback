// Package s3m does a thing.
package s3m

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/load"
	"github.com/gotracker/playback/format/settings"
)

type format struct{}

var (
	// S3M is the exported interface to the S3M file loader
	S3M = format{}
)

// LoadMOD loads a MOD file and upgrades it into an S3M file internally
func LoadMOD(filename string, s *settings.Settings) (playback.Playback, error) {
	return load.MOD(filename, s)
}

// Load loads an S3M file into a playback system
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	return load.S3M(filename, s)
}
