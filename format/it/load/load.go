package load

import (
	"github.com/gotracker/playback/format/settings"
	"github.com/gotracker/playback/player/intf"
)

// IT loads an IT file
func IT(filename string, s *settings.Settings) (intf.Playback, error) {
	return load(filename, readIT, s)
}
