package playback

import (
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/op"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

// Channel is an interface for channel state
type Channel[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	ResetRetriggerCount()
	SetMemory(*TMemory)
	GetMemory() *TMemory
	FreezePlayback()
	UnfreezePlayback()
	GetChannelData() TChannelData
	GetPortaTargetPeriod() TPeriod
	SetPortaTargetPeriod(TPeriod)
	SetPeriodOverride(TPeriod)
	SetPeriodDelta(period.Delta)
	GetPeriodDelta() period.Delta
	SetInstrument(*instrument.Instrument)
	GetVoice() voice.Voice
	GetPrevVoice() voice.Voice
	GetNoteSemitone() note.Semitone
	SetStoredSemitone(note.Semitone)
	SetTargetSemitone(note.Semitone)
	SetOverrideSemitone(note.Semitone)
	SetNotePlayTick(bool, note.Action, int)
	GetRetriggerCount() uint8
	SetRetriggerCount(uint8)
	SetPanEnabled(bool)
	SetRenderChannel(*render.Channel)
	GetRenderChannel() *render.Channel
	SetVolumeActive(bool)
	SetGlobalVolume(volume.Volume)
	SetChannelVolume(volume.Volume)
	GetChannelVolume() volume.Volume
	SetEnvelopePosition(int)
	TransitionActiveToPastState()
	SetNewNoteAction(note.Action)
	GetNewNoteAction() note.Action
	DoPastNoteEffect(action note.Action)
	SetVolumeEnvelopeEnable(bool)
	SetPanningEnvelopeEnable(bool)
	SetPitchEnvelopeEnable(bool)

	GetPreviousState() ChannelState[TPeriod]
	GetActiveState() *ChannelState[TPeriod]
	GetTargetState() *ChannelState[TPeriod]
}

type ChannelTargeter[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] func(out *op.ChannelTargets[TPeriod], d TChannelData, s song.Data, cs Channel[TPeriod, TMemory, TChannelData]) error
