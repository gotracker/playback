package load

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/settings"
)

// IT loads an IT file
func IT(filename string, s *settings.Settings) (playback.Playback, error) {
	return load(filename, readIT, s)
}
