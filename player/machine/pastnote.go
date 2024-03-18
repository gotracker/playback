package machine

import (
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/voice/mixer"
)

type pastNote[TPeriod Period] struct {
	rc *render.Channel[TPeriod]
}

func (p pastNote[TPeriod]) RenderAndAdvance(pc period.PeriodConverter[TPeriod], centerAheadPan volume.Matrix, details mixer.Details) (*mixing.Data, error) {
	return p.rc.RenderAndTick(pc, centerAheadPan, details)
}
