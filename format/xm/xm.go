// Package xm does a thing.
package xm

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/settings"
	"github.com/gotracker/playback/format/xm/load"
)

type format struct{}

var (
	// XM is the exported interface to the XM file loader
	XM = format{}
)

// Load loads an XM file into a playback system
func (f format) Load(filename string, s *settings.Settings) (playback.Playback, error) {
	return load.XM(filename, s)
}
