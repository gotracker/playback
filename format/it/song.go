package it

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/layout"
	itPlayback "github.com/gotracker/playback/format/it/playback"
)

type Song struct {
	*layout.Layout
}

func (s Song) ConstructPlayer() (playback.Playback, error) {
	return itPlayback.NewManager(s.Layout)
}
