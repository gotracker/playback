package playback

import (
	"io"

	"github.com/gotracker/playback/player/feature"
)

// Format is an interface to a music file format loader
type Format[TChannelData any] interface {
	Load(filename string, features []feature.Feature) (Playback, error)
	LoadFromReader(r io.Reader, features []feature.Feature) (Playback, error)
}
