// Package it does a thing.
package it

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/load"
	"github.com/gotracker/playback/settings"
	"github.com/gotracker/playback/util"
)

type format struct{}

var (
	// IT is the exported interface to the IT file loader
	IT = format{}
)

// Load loads an IT file into a playback system
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	r, err := util.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f.LoadFromReader(r, s)
}

// LoadFromReader loads an IT file on a reader into a playback system
func (f format) LoadFromReader(r io.Reader, s *settings.Settings) (playback.Playback, error) {
	return load.IT(r, s)
}
