package channel

import (
	"github.com/gotracker/playback"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
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

// SetTargetSemitone sets the target semitone for the channel
func (cs *State) SetTargetSemitone(st note.Semitone) {
	cs.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetTargetPeriod))
}

// SetOverrideSemitone sets the semitone override for the channel
func (cs *State) SetOverrideSemitone(st note.Semitone) {
	cs.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetPeriodOverride))
}

func (cs State) CalculateSemitonePeriod(st note.Semitone) period.Period {
	inst := cs.GetTargetInst()
	if inst == nil {
		return nil
	}

	cs.Semitone = note.Semitone(int(st) + int(inst.GetSemitoneShift()))
	return s3mPeriod.CalcSemitonePeriod(cs.Semitone, inst.GetFinetune(), inst.GetC2Spd())
}

func (cs *State) DoPortaByDelta(delta period.Delta) {
	cur := cs.GetPeriod()
	if cur == nil {
		return
	}

	sign := 1
	p := cur.AddDelta(delta, sign)
	cs.SetPeriod(p)
}
