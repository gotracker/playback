package load

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/common"
	itPlayback "github.com/gotracker/playback/format/it/playback"
	"github.com/gotracker/playback/player/feature"
)

// IT loads an IT file from a reader
func IT(r io.Reader, features []feature.Feature) (playback.Playback, error) {
	return common.Load(r, readIT, itPlayback.NewManager, features)
}
