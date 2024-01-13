package op

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice/types"
	"github.com/heucuva/optional"
)

type ChannelTargets[TPeriod types.Period, TVolume types.Volume, TPanning types.Panning] struct {
	NoteAction optional.Value[note.Action]
	NoteCalcST optional.Value[note.Semitone]

	TargetPos            optional.Value[sampling.Pos]
	TargetInst           optional.Value[instrument.InstrumentIntf]
	TargetPeriod         optional.Value[TPeriod]
	TargetStoredSemitone optional.Value[note.Semitone]
	TargetNewNoteAction  optional.Value[note.Action]
	TargetVolume         optional.Value[TVolume]
}
