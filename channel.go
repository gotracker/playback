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

// Channel is an interface for channel state
type Channel[TMemory, TChannelData any] interface {
	ResetRetriggerCount()
	SetMemory(*TMemory)
	GetMemory() *TMemory
	GetActiveVolume() volume.Volume
	SetActiveVolume(volume.Volume)
	FreezePlayback()
	UnfreezePlayback()
	GetData() *TChannelData
	GetPortaTargetPeriod() period.Period
	SetPortaTargetPeriod(period.Period)
	GetTargetPeriod() period.Period
	SetTargetPeriod(period.Period)
	SetPeriodOverride(period.Period)
	GetPeriod() period.Period
	SetPeriod(period.Period)
	SetPeriodDelta(period.PeriodDelta)
	GetPeriodDelta() period.PeriodDelta
	SetInstrument(*instrument.Instrument)
	GetInstrument() *instrument.Instrument
	GetVoice() voice.Voice
	GetTargetInst() *instrument.Instrument
	SetTargetInst(*instrument.Instrument)
	GetPrevInst() *instrument.Instrument
	GetPrevVoice() voice.Voice
	GetNoteSemitone() note.Semitone
	SetStoredSemitone(note.Semitone)
	SetTargetSemitone(note.Semitone)
	SetOverrideSemitone(note.Semitone)
	GetTargetPos() sampling.Pos
	SetTargetPos(sampling.Pos)
	GetPos() sampling.Pos
	SetPos(sampling.Pos)
	SetNotePlayTick(bool, note.Action, int)
	GetRetriggerCount() uint8
	SetRetriggerCount(uint8)
	SetPanEnabled(bool)
	GetPan() panning.Position
	SetPan(panning.Position)
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
	NoteCut()
}
