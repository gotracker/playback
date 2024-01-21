package load

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

// IT loads an IT file from a reader
func IT(r io.Reader, features []feature.Feature) (song.Data, error) {
	return common.Load(r, readIT, features)
}
