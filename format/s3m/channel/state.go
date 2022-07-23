package channel

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/player/state"
)

type State struct {
	state.ChannelState
	state.ChannelMemory[Memory]
	state.ChannelData[Data, State]
	state.ChannelStateSemitoneSetter[State]
	state.ChannelEffect[State]
}

func (cs *State) SetData(cdata *Data) error {
	return cs.ChannelData.SetData(cdata, cs)
}

func (cs *State) AdvanceRow(effectFactory playback.EffectFactory[Data, State]) {
	cs.ChannelState.AdvanceRow()
	cs.ChannelData.AdvanceRow(&dataTxn{
		EffectFactory: effectFactory,
	})
}
