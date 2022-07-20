// Package it does a thing.
package it

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/load"
	"github.com/gotracker/playback/settings"
)

type format struct{}

var (
	// IT is the exported interface to the IT file loader
	IT = format{}
)

// Load loads an IT file into a playback system
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	return load.IT(filename, s)
}
