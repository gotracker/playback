package load

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/settings"
)

// XM loads an XM file and upgrades it into an XM file internally
func XM(filename string, s *settings.Settings) (playback.Playback, error) {
	return load(filename, readXM, s)
}
