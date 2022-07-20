package load

import (
	"io"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/common"
	xmPlayback "github.com/gotracker/playback/format/xm/playback"
	"github.com/gotracker/playback/settings"
)

// XM loads an XM file and upgrades it into an XM file internally
func XM(r io.Reader, s *settings.Settings) (playback.Playback, error) {
	return common.Load(r, readXM, xmPlayback.NewManager, s)
}
