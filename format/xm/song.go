package xm

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/layout"
	xmPlayback "github.com/gotracker/playback/format/xm/playback"
)

type Song struct {
	*layout.Layout
}

func (s Song) ConstructPlayer() (playback.Playback, error) {
	return xmPlayback.NewManager(s.Layout)
}
