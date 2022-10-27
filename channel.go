package playback

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

type Channel[TMemory, TChannelData any] interface {
	ChannelState
	ChannelMemory[TMemory]
	ChannelData[TChannelData]

	FreezePlayback()
	UnfreezePlayback()
	NoteCut()
	ResetRetriggerCount()
	SetActiveVolume(av volume.Volume)
	SetVolumeActive(enabled bool)
	SetChannelVolume(cv volume.Volume)
	SetVolumeEnvelopeEnable(enabled bool)
	SetPanningEnvelopeEnable(enabled bool)
	SetPitchEnvelopeEnable(enabled bool)
	SetEnvelopePosition(pos int)
	SetGlobalVolume(gv volume.Volume)
	SetInstrument(inst *instrument.Instrument)
	SetNewNoteAction(action note.Action)
	SetNotePlayTick(enabled bool, action note.Action, tick int)
	SetOverrideSemitone(st note.Semitone)
	SetPan(pos panning.Position)
	SetPanEnabled(enabled bool)
	SetPeriod(p period.Period)
	SetPeriodDelta(delta period.PeriodDelta)
	SetTargetPeriod(p period.Period)
	SetPeriodOverride(p period.Period)
	SetPortaTargetPeriod(p period.Period)
	SetPos(pos sampling.Pos)
	SetRenderChannel(rc *render.Channel)
	SetTargetInst(inst *instrument.Instrument)
	SetStoredSemitone(st note.Semitone)
	SetTargetSemitone(st note.Semitone)
	SetTargetPos(pos sampling.Pos)
	SetRetriggerCount(c uint8)
	TransitionActiveToPastState()
}

type ChannelMemory[TMemory any] interface {
	SetMemory(*TMemory)
	GetMemory() *TMemory
}

type ChannelData[TChannelData any] interface {
	GetData() *TChannelData
}

type ChannelEffect[TChannelState ChannelState] interface {
	GetActiveEffect() Effecter[TChannelState]
}

// ChannelState is an interface for channel state
type ChannelState interface {
	GetActiveVolume() volume.Volume
	GetPortaTargetPeriod() period.Period
	GetTargetPeriod() period.Period
	GetPeriod() period.Period
	GetPeriodDelta() period.PeriodDelta
	GetInstrument() *instrument.Instrument
	GetVoice() voice.Voice
	GetTargetInst() *instrument.Instrument
	GetPrevInst() *instrument.Instrument
	GetPrevVoice() voice.Voice
	GetNoteSemitone() note.Semitone
	GetTargetPos() sampling.Pos
	GetPos() sampling.Pos
	GetRetriggerCount() uint8
	GetPan() panning.Position
	GetRenderChannel() *render.Channel
	GetChannelVolume() volume.Volume
	GetNewNoteAction() note.Action
}
