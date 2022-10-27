package state

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/render"
	voiceImpl "github.com/gotracker/playback/player/voice"
	"github.com/gotracker/playback/song"
	"github.com/heucuva/optional"
)

type NoteTrigger struct {
	NoteAction note.Action
	Tick       int
}

type VolOp[TChannelState any] interface {
	Process(p playback.Playback, cs *TChannelState) error
}

type NoteOp[TChannelState any] interface {
	Process(p playback.Playback, cs *TChannelState) error
}

// ChannelState is the state of a single channel
type ChannelState struct {
	ActiveState Active
	TargetState Playback
	PrevState   Active

	s song.Data

	StoredSemitone    note.Semitone // from pattern, unmodified, current note
	PortaTargetPeriod optional.Value[period.Period]
	Trigger           optional.Value[NoteTrigger]
	RetriggerCount    uint8
	freezePlayback    bool
	Semitone          note.Semitone // from TargetSemitone, modified further, used in period calculations
	UseTargetPeriod   bool
	periodOverride    period.Period
	UsePeriodOverride bool
	volumeActive      bool
	PanEnabled        bool
	NewNoteAction     note.Action

	RenderChannel *render.Channel
}

// WillTriggerOn returns true if a note will trigger on the tick specified
func (cs *ChannelState) WillTriggerOn(tick int) (bool, note.Action) {
	if trigger, ok := cs.Trigger.Get(); ok {
		return trigger.Tick == tick, trigger.NoteAction
	}

	return false, note.ActionContinue
}

// RenderRowTick renders a channel's row data for a single tick
func (cs *ChannelState) RenderRowTick(details RenderDetails, pastNotes []*Active) ([]mixing.Data, error) {
	if cs.PlaybackFrozen() {
		return nil, nil
	}

	mixData := RenderStatesTogether(&cs.ActiveState, pastNotes, details)

	return mixData, nil
}

// ResetStates resets the channel's internal states
func (cs *ChannelState) ResetStates() {
	cs.ActiveState.Reset()
	cs.TargetState.Reset()
	cs.PrevState.Reset()
}

// FreezePlayback suspends mixer progression on the channel
func (cs *ChannelState) FreezePlayback() {
	cs.freezePlayback = true
}

// UnfreezePlayback resumes mixer progression on the channel
func (cs *ChannelState) UnfreezePlayback() {
	cs.freezePlayback = false
}

// PlaybackFrozen returns true if the mixer progression for the channel is suspended
func (cs ChannelState) PlaybackFrozen() bool {
	return cs.freezePlayback
}

// ResetRetriggerCount sets the retrigger count to 0
func (cs *ChannelState) ResetRetriggerCount() {
	cs.RetriggerCount = 0
}

// GetActiveVolume returns the current active volume on the channel
func (cs ChannelState) GetActiveVolume() volume.Volume {
	return cs.ActiveState.Volume
}

// SetActiveVolume sets the active volume on the channel
func (cs *ChannelState) SetActiveVolume(vol volume.Volume) {
	if vol != volume.VolumeUseInstVol {
		cs.ActiveState.Volume = vol
	}
}

func (cs *ChannelState) SetSongDataInterface(s song.Data) {
	cs.s = s
}

func (cs ChannelState) GetSongDataInterface() song.Data {
	return cs.s
}

// GetPortaTargetPeriod returns the current target portamento (to note) sampler period
func (cs ChannelState) GetPortaTargetPeriod() period.Period {
	if p, ok := cs.PortaTargetPeriod.Get(); ok {
		return p
	}
	return nil
}

// SetPortaTargetPeriod sets the current target portamento (to note) sampler period
func (cs *ChannelState) SetPortaTargetPeriod(period period.Period) {
	if period != nil {
		cs.PortaTargetPeriod.Set(period)
	} else {
		cs.PortaTargetPeriod.Reset()
	}
}

// GetTargetPeriod returns the soon-to-be-committed sampler period (when the note retriggers)
func (cs ChannelState) GetTargetPeriod() period.Period {
	return cs.TargetState.Period
}

// SetTargetPeriod sets the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState) SetTargetPeriod(period period.Period) {
	cs.TargetState.Period = period
}

// GetTargetPeriod returns the soon-to-be-committed sampler period (when the note retriggers)
func (cs ChannelState) GetPeriodOverride() period.Period {
	return cs.periodOverride
}

// SetTargetPeriod sets the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState) SetPeriodOverride(period period.Period) {
	cs.periodOverride = period
	cs.UsePeriodOverride = true
}

// SetPeriodDelta sets the vibrato (ephemeral) delta sampler period
func (cs *ChannelState) SetPeriodDelta(delta period.PeriodDelta) {
	cs.ActiveState.PeriodDelta = delta
}

// GetPeriodDelta gets the vibrato (ephemeral) delta sampler period
func (cs ChannelState) GetPeriodDelta() period.PeriodDelta {
	return cs.ActiveState.PeriodDelta
}

// SetVolumeActive enables or disables the sample of the instrument
func (cs *ChannelState) SetVolumeActive(on bool) {
	cs.volumeActive = on
}

// GetInstrument returns the interface to the active instrument
func (cs ChannelState) GetInstrument() *instrument.Instrument {
	return cs.ActiveState.Instrument
}

// SetInstrument sets the interface to the active instrument
func (cs *ChannelState) SetInstrument(inst *instrument.Instrument) {
	cs.ActiveState.Instrument = inst
	if inst != nil {
		if inst == cs.PrevState.Instrument {
			cs.ActiveState.Voice = cs.PrevState.Voice
		} else {
			cs.ActiveState.Voice = voiceImpl.New(inst, cs.RenderChannel)
		}
	}
}

// GetVoice returns the active voice interface
func (cs ChannelState) GetVoice() voice.Voice {
	return cs.ActiveState.Voice
}

// GetTargetInst returns the interface to the soon-to-be-committed active instrument (when the note retriggers)
func (cs ChannelState) GetTargetInst() *instrument.Instrument {
	return cs.TargetState.Instrument
}

// SetTargetInst sets the soon-to-be-committed active instrument (when the note retriggers)
func (cs *ChannelState) SetTargetInst(inst *instrument.Instrument) {
	cs.TargetState.Instrument = inst
}

// GetPrevInst returns the interface to the last row's active instrument
func (cs ChannelState) GetPrevInst() *instrument.Instrument {
	return cs.PrevState.Instrument
}

// GetPrevVoice returns the interface to the last row's active voice
func (cs ChannelState) GetPrevVoice() voice.Voice {
	return cs.PrevState.Voice
}

// GetNoteSemitone returns the note semitone for the channel
func (cs ChannelState) GetNoteSemitone() note.Semitone {
	return cs.StoredSemitone
}

// GetTargetPos returns the soon-to-be-committed sample position of the instrument
func (cs ChannelState) GetTargetPos() sampling.Pos {
	return cs.TargetState.Pos
}

// SetTargetPos sets the soon-to-be-committed sample position of the instrument
func (cs *ChannelState) SetTargetPos(pos sampling.Pos) {
	cs.TargetState.Pos = pos
}

// GetPeriod returns the current sampler period of the active instrument
func (cs ChannelState) GetPeriod() period.Period {
	return cs.ActiveState.Period
}

// SetPeriod sets the current sampler period of the active instrument
func (cs *ChannelState) SetPeriod(period period.Period) {
	cs.ActiveState.Period = period
}

// GetPos returns the sample position of the active instrument
func (cs ChannelState) GetPos() sampling.Pos {
	return cs.ActiveState.Pos
}

// SetPos sets the sample position of the active instrument
func (cs *ChannelState) SetPos(pos sampling.Pos) {
	cs.ActiveState.Pos = pos
}

// SetNotePlayTick sets the tick on which the note will retrigger
func (cs *ChannelState) SetNotePlayTick(enabled bool, action note.Action, tick int) {
	if enabled {
		cs.Trigger.Set(NoteTrigger{
			NoteAction: action,
			Tick:       tick,
		})
	} else {
		cs.Trigger.Reset()
	}
}

// GetRetriggerCount returns the current count of the retrigger counter
func (cs ChannelState) GetRetriggerCount() uint8 {
	return cs.RetriggerCount
}

// SetRetriggerCount sets the current count of the retrigger counter
func (cs *ChannelState) SetRetriggerCount(cnt uint8) {
	cs.RetriggerCount = cnt
}

// SetPanEnabled activates or deactivates the panning. If enabled, then pan updates work (see SetPan)
func (cs *ChannelState) SetPanEnabled(on bool) {
	cs.PanEnabled = on
}

// SetPan sets the active panning value of the channel
func (cs *ChannelState) SetPan(pan panning.Position) {
	if cs.PanEnabled {
		cs.ActiveState.Pan = pan
	}
}

// GetPan gets the active panning value of the channel
func (cs ChannelState) GetPan() panning.Position {
	return cs.ActiveState.Pan
}

func (cs *ChannelState) AdvanceRow() {
	cs.PrevState = cs.ActiveState
	cs.TargetState = cs.ActiveState.Playback
	cs.Trigger.Reset()
	cs.RetriggerCount = 0
	cs.ActiveState.PeriodDelta = 0
	cs.UseTargetPeriod = false
}

// SetTargetSemitone sets the target semitone for the channel
func (cs *ChannelState) SetTargetSemitone(st note.Semitone) {
	// TODO: this should be overridden with a setter that knows how to convert the semitone
	// ChannelData[channel.Data].AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetTargetPeriod))
}

// SetOverrideSemitone sets the semitone override for the channel
func (cs *ChannelState) SetOverrideSemitone(st note.Semitone) {
	// TODO: this should be overridden with a setter that knows how to convert the semitone
	//ChannelData[channel.Data].AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetPeriodOverride))
}

// SetStoredSemitone sets the stored semitone for the channel
func (cs *ChannelState) SetStoredSemitone(st note.Semitone) {
	cs.StoredSemitone = st
}

// SetRenderChannel sets the output channel for the channel
func (cs *ChannelState) SetRenderChannel(outputCh *render.Channel) {
	cs.RenderChannel = outputCh
}

// GetRenderChannel returns the output channel for the channel
func (cs ChannelState) GetRenderChannel() *render.Channel {
	return cs.RenderChannel
}

// SetGlobalVolume sets the last-known global volume on the channel
func (cs *ChannelState) SetGlobalVolume(gv volume.Volume) {
	cs.RenderChannel.LastGlobalVolume = gv
	cs.RenderChannel.SetGlobalVolume(gv)
}

// SetChannelVolume sets the channel volume on the channel
func (cs *ChannelState) SetChannelVolume(cv volume.Volume) {
	cs.RenderChannel.ChannelVolume = cv
}

// GetChannelVolume gets the channel volume on the channel
func (cs ChannelState) GetChannelVolume() volume.Volume {
	return cs.RenderChannel.ChannelVolume
}

// SetEnvelopePosition sets the envelope position for the active instrument
func (cs *ChannelState) SetEnvelopePosition(ticks int) {
	if nc := cs.GetVoice(); nc != nil {
		voice.SetVolumeEnvelopePosition(nc, ticks)
		voice.SetPanEnvelopePosition(nc, ticks)
		voice.SetPitchEnvelopePosition(nc, ticks)
		voice.SetFilterEnvelopePosition(nc, ticks)
	}
}

// TransitionActiveToPastState will transition the current active state to the 'past' state
// and will activate the specified New-Note Action on it
func (cs *ChannelState) TransitionActiveToPastState() {
	cs.ActiveState.Reset()
}

// SetNewNoteAction sets the New-Note Action on the channel
func (cs *ChannelState) SetNewNoteAction(nna note.Action) {
	cs.NewNoteAction = nna
}

// GetNewNoteAction gets the New-Note Action on the channel
func (cs ChannelState) GetNewNoteAction() note.Action {
	return cs.NewNoteAction
}

// SetVolumeEnvelopeEnable sets the enable flag on the active volume envelope
func (cs *ChannelState) SetVolumeEnvelopeEnable(enabled bool) {
	voice.EnableVolumeEnvelope(cs.ActiveState.Voice, enabled)
}

// SetPanningEnvelopeEnable sets the enable flag on the active panning envelope
func (cs *ChannelState) SetPanningEnvelopeEnable(enabled bool) {
	voice.EnablePanEnvelope(cs.ActiveState.Voice, enabled)
}

// SetPitchEnvelopeEnable sets the enable flag on the active pitch/filter envelope
func (cs *ChannelState) SetPitchEnvelopeEnable(enabled bool) {
	voice.EnablePitchEnvelope(cs.ActiveState.Voice, enabled)
}

func (cs *ChannelState) NoteCut() {
	cs.ActiveState.Period = nil
}
