package s3m

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/layout"
	s3mPlayback "github.com/gotracker/playback/format/s3m/playback"
)

type Song struct {
	*layout.Layout
}

func (s Song) ConstructPlayer() (playback.Playback, error) {
	return s3mPlayback.NewManager(s.Layout)
}
