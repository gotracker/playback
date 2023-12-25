package tracing

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

func ChannelStateHeaders(preamble string) []string {
	return []string{
		fmt.Sprintf("%s Instrument", preamble),
		fmt.Sprintf("%s Period", preamble),
		fmt.Sprintf("%s Volume", preamble),
		fmt.Sprintf("%s Position", preamble),
		fmt.Sprintf("%s Pan", preamble),
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
