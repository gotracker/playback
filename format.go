package playback

import "github.com/gotracker/playback/settings"

// Format is an interface to a music file format loader
type Format[TChannelData any] interface {
	Load(filename string, s *settings.Settings) (Playback, error)
}
