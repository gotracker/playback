package state

import "github.com/gotracker/playback"

type ChannelEffect[TChannelState playback.ChannelState] struct {
	ActiveEffect playback.Effecter[TChannelState]
}

func (cs ChannelEffect[TChannelState]) GetActiveEffect() playback.Effecter[TChannelState] {
	return cs.ActiveEffect
}

func (cs *ChannelEffect[TChannelState]) SetActiveEffect(e playback.Effecter[TChannelState]) {
	cs.ActiveEffect = e
}
