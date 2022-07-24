package channel

import (
	"github.com/gotracker/playback"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
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

// SetTargetSemitone sets the target semitone for the channel
func (cs *State) SetTargetSemitone(st note.Semitone) {
	cs.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetTargetPeriod))
}

// SetOverrideSemitone sets the semitone override for the channel
func (cs *State) SetOverrideSemitone(st note.Semitone) {
	cs.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetPeriodOverride))
}

func (cs State) IsLinearFreqSlides() bool {
	return cs.Memory.Shared.LinearFreqSlides
}

func (cs State) CalculateSemitonePeriod(st note.Semitone) period.Period {
	inst := cs.GetTargetInst()
	if inst == nil {
		return nil
	}

	cs.Semitone = note.Semitone(int(st) + int(inst.GetSemitoneShift()))
	return itPeriod.CalcSemitonePeriod(cs.Semitone, inst.GetFinetune(), inst.GetC2Spd(), cs.IsLinearFreqSlides())
}

func (cs *State) DoPortaByDelta(delta period.Delta) {
	cur := cs.GetPeriod()
	if cur == nil {
		return
	}

	sign := 1
	if cs.IsLinearFreqSlides() {
		sign = -1
	}

	cur = cur.AddDelta(delta, sign)
	cs.SetPeriod(cur)
}
