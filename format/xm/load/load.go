package load

import (
	"io"

	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

// XM loads an XM file and upgrades it into an XM file internally
func XM(r io.Reader, features []feature.Feature) (song.Data, error) {
	return common.Load(r, readXM, features)
}
