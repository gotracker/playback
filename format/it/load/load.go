package load

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/it/layout"
	"github.com/gotracker/playback/player/feature"
)

// IT loads an IT file from a reader
func IT(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	return common.Load(r, readIT, features)
}
