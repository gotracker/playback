package machine

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/heucuva/optional"
)

type ActionTick struct {
	Action note.Action
	Tick   int
}

type NewNoteInfo[TPeriod Period, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	Period        optional.Value[TPeriod]
	Inst          optional.Value[*instrument.Instrument[TMixingVolume, TVolume, TPanning]]
	Pos           optional.Value[sampling.Pos]
	MixVol        optional.Value[TMixingVolume]
	Vol           optional.Value[TVolume]
	Pan           optional.Value[TPanning]
	ActionTick    optional.Value[ActionTick]
	NewNoteAction optional.Value[note.Action]
}

func (n NewNoteInfo[TPeriod, TMixingVolume, TVolume, TPanning]) IsSet() bool {
	return n.Period.IsSet() ||
		n.Inst.IsSet() ||
		n.Pos.IsSet() ||
		n.MixVol.IsSet() ||
		n.Vol.IsSet() ||
		n.Pan.IsSet() ||
		n.ActionTick.IsSet() ||
		n.NewNoteAction.IsSet()
}