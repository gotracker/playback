package load

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/settings"
)

// IT loads an IT file
func IT(filename string, s *settings.Settings) (playback.Playback, error) {
	return load(filename, readIT, s)
}
