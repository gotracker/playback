package channel

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/state"
)

type State struct {
	state.ChannelState
	state.ChannelMemory[Memory]
	state.ChannelData[Data, State]
	state.ChannelStateSemitoneSetter[State]
	state.ChannelEffect[State]

	PastNotes *state.PastNotesProcessor
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

// DoPastNoteEffect performs an action on all past-note playbacks associated with the channel
func (cs *State) DoPastNoteEffect(action note.Action) {
	cs.PastNotes.Do(cs.RenderChannel.ChannelNum, action)
}

// TransitionActiveToPastState will transition the current active state to the 'past' state
// and will activate the specified New-Note Action on it
func (cs *State) TransitionActiveToPastState() {
	if cs.PastNotes != nil {
		switch cs.NewNoteAction {
		case note.ActionCut:
			// reset at end

		case note.ActionContinue:
			// nothing
			pn := cs.ActiveState.Clone()
			if nc := pn.Voice; nc != nil {
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}

		case note.ActionRelease:
			pn := cs.ActiveState.Clone()
			if nc := pn.Voice; nc != nil {
				nc.Release()
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}

		case note.ActionFadeout:
			pn := cs.ActiveState.Clone()
			if nc := pn.Voice; nc != nil {
				nc.Release()
				nc.Fadeout()
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}
		}
	}
	cs.ChannelState.TransitionActiveToPastState()
}
