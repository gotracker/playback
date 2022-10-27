package load

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/xm/layout"
	"github.com/gotracker/playback/player/feature"
)

// XM loads an XM file from a reader
func XM(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	return common.Load(r, readXM, features)
}
