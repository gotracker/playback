package playback

import (
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/op"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

// Channel is an interface for channel state
type Channel[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] interface {
	ResetRetriggerCount()
	SetMemory(TMemory)
	GetMemory() TMemory
	FreezePlayback()
	UnfreezePlayback()
	GetChannelData() TChannelData
	GetPortaTargetPeriod() TPeriod
	SetPortaTargetPeriod(TPeriod)
	SetPeriodOverride(TPeriod)
	SetPeriodDelta(period.Delta)
	GetPeriodDelta() period.Delta
	SetInstrument(*instrument.Instrument[TMixingVolume, TVolume, TPanning])
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
	SetRenderChannel(*render.Channel[TGlobalVolume, TMixingVolume, TPanning])
	GetRenderChannel() *render.Channel[TGlobalVolume, TMixingVolume, TPanning]
	SetVolumeActive(bool)
	SetGlobalVolume(TGlobalVolume)
	SetChannelVolume(TMixingVolume)
	GetChannelVolume() TMixingVolume
	SetEnvelopePosition(int)
	TransitionActiveToPastState()
	SetNewNoteAction(note.Action)
	GetNewNoteAction() note.Action
	DoPastNoteEffect(action note.Action)
	SetVolumeEnvelopeEnable(bool)
	SetPanningEnvelopeEnable(bool)
	SetPitchEnvelopeEnable(bool)
	GetActiveEffects() []Effect
	GetUseTargetPeriod() bool

	GetPreviousState() ChannelState[TPeriod, TVolume, TPanning]
	GetActiveState() *ChannelState[TPeriod, TVolume, TPanning]
	GetTargetState() *ChannelState[TPeriod, TVolume, TPanning]
}

type ChannelTargeter[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] func(out *op.ChannelTargets[TPeriod, TVolume, TPanning], d TChannelData, s song.Data, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error
