package tracing

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

func ChannelStateHeaders() []string {
	return []string{
		"Instrument",
		"Period",
		"Volume",
		"Position",
		"Pan",
	}
}

func ChannelState[TPeriod period.Period](cs *playback.ChannelState[TPeriod]) []any {
	var instStr string
	if cs.Instrument != nil {
		instStr = cs.Instrument.GetID().String()
	}
	return []any{
		instStr,
		cs.Period,
		cs.GetVolume(),
		cs.Pos,
		cs.Pan,
	}
}
