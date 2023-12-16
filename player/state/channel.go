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

type VolOp[TMemory any] interface {
	Process(p playback.Playback, cs *ChannelState[TMemory]) error
}

type NoteOp[TMemory any] interface {
	Process(p playback.Playback, cs *ChannelState[TMemory]) error
}

type PeriodUpdateFunc func(period.Period)

type SemitoneSetterFactory[TMemory any] func(note.Semitone, PeriodUpdateFunc) NoteOp[TMemory]

// ChannelState is the state of a single channel
type ChannelState[TMemory any] struct {
	activeState Active
	targetState Playback
	prevState   Active

	ActiveEffect playback.Effect

	s       song.Data
	txn     ChannelDataTransaction[TMemory]
	prevTxn ChannelDataTransaction[TMemory]

	SemitoneSetterFactory SemitoneSetterFactory[TMemory]

	StoredSemitone    note.Semitone // from pattern, unmodified, current note
	PortaTargetPeriod optional.Value[period.Period]
	Trigger           optional.Value[NoteTrigger]
	RetriggerCount    uint8
	Memory            *TMemory
	freezePlayback    bool
	Semitone          note.Semitone // from TargetSemitone, modified further, used in period calculations
	UseTargetPeriod   bool
	periodOverride    period.Period
	UsePeriodOverride bool
	volumeActive      bool
	PanEnabled        bool
	NewNoteAction     note.Action

	PastNotes     *PastNotesProcessor
	RenderChannel *render.Channel
}

// WillTriggerOn returns true if a note will trigger on the tick specified
func (cs *ChannelState[TMemory]) WillTriggerOn(tick int) (bool, note.Action) {
	if trigger, ok := cs.Trigger.Get(); ok {
		return trigger.Tick == tick, trigger.NoteAction
	}

	return false, note.ActionContinue
}

// AdvanceRow will update the current state to make room for the next row's state data
func (cs *ChannelState[TMemory]) AdvanceRow(txn ChannelDataTransaction[TMemory]) {
	cs.prevState = cs.activeState
	cs.targetState = cs.activeState.Playback
	cs.Trigger.Reset()
	cs.RetriggerCount = 0
	cs.activeState.PeriodDelta = 0

	cs.UseTargetPeriod = false
	cs.prevTxn = cs.txn
	cs.txn = txn
}

// RenderRowTick renders a channel's row data for a single tick
func (cs *ChannelState[TMemory]) RenderRowTick(details RenderDetails, pastNotes []*Active) ([]mixing.Data, error) {
	if cs.PlaybackFrozen() {
		return nil, nil
	}

	mixData := RenderStatesTogether(&cs.activeState, pastNotes, details)

	return mixData, nil
}

// ResetStates resets the channel's internal states
func (cs *ChannelState[TMemory]) ResetStates() {
	cs.activeState.Reset()
	cs.targetState.Reset()
	cs.prevState.Reset()
}

func (cs *ChannelState[TMemory]) GetActiveEffect() playback.Effect {
	return cs.ActiveEffect
}

func (cs *ChannelState[TMemory]) SetActiveEffect(e playback.Effect) {
	cs.ActiveEffect = e
}

// FreezePlayback suspends mixer progression on the channel
func (cs *ChannelState[TMemory]) FreezePlayback() {
	cs.freezePlayback = true
}

// UnfreezePlayback resumes mixer progression on the channel
func (cs *ChannelState[TMemory]) UnfreezePlayback() {
	cs.freezePlayback = false
}

// PlaybackFrozen returns true if the mixer progression for the channel is suspended
func (cs ChannelState[TMemory]) PlaybackFrozen() bool {
	return cs.freezePlayback
}

// ResetRetriggerCount sets the retrigger count to 0
func (cs *ChannelState[TMemory]) ResetRetriggerCount() {
	cs.RetriggerCount = 0
}

// GetMemory returns the interface to the custom effect memory module
func (cs *ChannelState[TMemory]) GetMemory() *TMemory {
	return cs.Memory
}

// SetMemory sets the custom effect memory interface
func (cs *ChannelState[TMemory]) SetMemory(mem *TMemory) {
	cs.Memory = mem
}

// GetActiveVolume returns the current active volume on the channel
func (cs *ChannelState[TMemory]) GetActiveVolume() volume.Volume {
	return cs.activeState.Volume
}

// SetActiveVolume sets the active volume on the channel
func (cs *ChannelState[TMemory]) SetActiveVolume(vol volume.Volume) {
	if vol != volume.VolumeUseInstVol {
		cs.activeState.Volume = vol
	}
}

func (cs *ChannelState[TMemory]) SetSongDataInterface(s song.Data) {
	cs.s = s
}

// GetData returns the interface to the current channel song pattern data
func (cs *ChannelState[TMemory]) GetData() song.ChannelData {
	if cs.txn == nil {
		return nil
	}

	return cs.txn.GetData()
}

func (cs *ChannelState[TMemory]) SetData(cdata song.ChannelData) error {
	if cs.txn == nil {
		return nil
	}

	return cs.txn.SetData(cdata, cs.s, cs)
}

func (cs *ChannelState[TMemory]) GetTxn() ChannelDataTransaction[TMemory] {
	return cs.txn
}

// GetPortaTargetPeriod returns the current target portamento (to note) sampler period
func (cs *ChannelState[TMemory]) GetPortaTargetPeriod() period.Period {
	if p, ok := cs.PortaTargetPeriod.Get(); ok {
		return p
	}
	return nil
}

// SetPortaTargetPeriod sets the current target portamento (to note) sampler period
func (cs *ChannelState[TMemory]) SetPortaTargetPeriod(period period.Period) {
	if period != nil {
		cs.PortaTargetPeriod.Set(period)
	} else {
		cs.PortaTargetPeriod.Reset()
	}
}

// GetTargetPeriod returns the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState[TMemory]) GetTargetPeriod() period.Period {
	return cs.targetState.Period
}

// SetTargetPeriod sets the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState[TMemory]) SetTargetPeriod(period period.Period) {
	cs.targetState.Period = period
}

// GetTargetPeriod returns the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState[TMemory]) GetPeriodOverride() period.Period {
	return cs.periodOverride
}

// SetTargetPeriod sets the soon-to-be-committed sampler period (when the note retriggers)
func (cs *ChannelState[TMemory]) SetPeriodOverride(period period.Period) {
	cs.periodOverride = period
	cs.UsePeriodOverride = true
}

// SetPeriodDelta sets the vibrato (ephemeral) delta sampler period
func (cs *ChannelState[TMemory]) SetPeriodDelta(delta period.PeriodDelta) {
	cs.activeState.PeriodDelta = delta
}

// GetPeriodDelta gets the vibrato (ephemeral) delta sampler period
func (cs *ChannelState[TMemory]) GetPeriodDelta() period.PeriodDelta {
	return cs.activeState.PeriodDelta
}

// SetVolumeActive enables or disables the sample of the instrument
func (cs *ChannelState[TMemory]) SetVolumeActive(on bool) {
	cs.volumeActive = on
}

// GetInstrument returns the interface to the active instrument
func (cs *ChannelState[TMemory]) GetInstrument() *instrument.Instrument {
	return cs.activeState.Instrument
}

// SetInstrument sets the interface to the active instrument
func (cs *ChannelState[TMemory]) SetInstrument(inst *instrument.Instrument) {
	cs.activeState.Instrument = inst
	if inst != nil {
		if inst == cs.prevState.Instrument {
			cs.activeState.Voice = cs.prevState.Voice
		} else {
			cs.activeState.Voice = voiceImpl.New(inst, cs.RenderChannel)
		}
	}
}

// GetVoice returns the active voice interface
func (cs *ChannelState[TMemory]) GetVoice() voice.Voice {
	return cs.activeState.Voice
}

// GetTargetInst returns the interface to the soon-to-be-committed active instrument (when the note retriggers)
func (cs *ChannelState[TMemory]) GetTargetInst() *instrument.Instrument {
	return cs.targetState.Instrument
}

// SetTargetInst sets the soon-to-be-committed active instrument (when the note retriggers)
func (cs *ChannelState[TMemory]) SetTargetInst(inst *instrument.Instrument) {
	cs.targetState.Instrument = inst
}

// GetPrevInst returns the interface to the last row's active instrument
func (cs *ChannelState[TMemory]) GetPrevInst() *instrument.Instrument {
	return cs.prevState.Instrument
}

// GetPrevVoice returns the interface to the last row's active voice
func (cs *ChannelState[TMemory]) GetPrevVoice() voice.Voice {
	return cs.prevState.Voice
}

// GetNoteSemitone returns the note semitone for the channel
func (cs *ChannelState[TMemory]) GetNoteSemitone() note.Semitone {
	return cs.StoredSemitone
}

// GetTargetPos returns the soon-to-be-committed sample position of the instrument
func (cs *ChannelState[TMemory]) GetTargetPos() sampling.Pos {
	return cs.targetState.Pos
}

// SetTargetPos sets the soon-to-be-committed sample position of the instrument
func (cs *ChannelState[TMemory]) SetTargetPos(pos sampling.Pos) {
	cs.targetState.Pos = pos
}

// GetPeriod returns the current sampler period of the active instrument
func (cs *ChannelState[TMemory]) GetPeriod() period.Period {
	return cs.activeState.Period
}

// SetPeriod sets the current sampler period of the active instrument
func (cs *ChannelState[TMemory]) SetPeriod(period period.Period) {
	cs.activeState.Period = period
}

// GetPos returns the sample position of the active instrument
func (cs *ChannelState[TMemory]) GetPos() sampling.Pos {
	return cs.activeState.Pos
}

// SetPos sets the sample position of the active instrument
func (cs *ChannelState[TMemory]) SetPos(pos sampling.Pos) {
	cs.activeState.Pos = pos
}

// SetNotePlayTick sets the tick on which the note will retrigger
func (cs *ChannelState[TMemory]) SetNotePlayTick(enabled bool, action note.Action, tick int) {
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
func (cs *ChannelState[TMemory]) GetRetriggerCount() uint8 {
	return cs.RetriggerCount
}

// SetRetriggerCount sets the current count of the retrigger counter
func (cs *ChannelState[TMemory]) SetRetriggerCount(cnt uint8) {
	cs.RetriggerCount = cnt
}

// SetPanEnabled activates or deactivates the panning. If enabled, then pan updates work (see SetPan)
func (cs *ChannelState[TMemory]) SetPanEnabled(on bool) {
	cs.PanEnabled = on
}

// SetPan sets the active panning value of the channel
func (cs *ChannelState[TMemory]) SetPan(pan panning.Position) {
	if cs.PanEnabled {
		cs.activeState.Pan = pan
	}
}

// GetPan gets the active panning value of the channel
func (cs *ChannelState[TMemory]) GetPan() panning.Position {
	return cs.activeState.Pan
}

// SetTargetSemitone sets the target semitone for the channel
func (cs *ChannelState[TMemory]) SetTargetSemitone(st note.Semitone) {
	if cs.txn != nil {
		cs.txn.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetTargetPeriod))
	}
}

func (cs *ChannelState[TMemory]) SetOverrideSemitone(st note.Semitone) {
	if cs.txn != nil {
		cs.txn.AddNoteOp(cs.SemitoneSetterFactory(st, cs.SetPeriodOverride))
	}
}

// SetStoredSemitone sets the stored semitone for the channel
func (cs *ChannelState[TMemory]) SetStoredSemitone(st note.Semitone) {
	cs.StoredSemitone = st
}

// SetRenderChannel sets the output channel for the channel
func (cs *ChannelState[TMemory]) SetRenderChannel(outputCh *render.Channel) {
	cs.RenderChannel = outputCh
}

// GetRenderChannel returns the output channel for the channel
func (cs *ChannelState[TMemory]) GetRenderChannel() *render.Channel {
	return cs.RenderChannel
}

// SetGlobalVolume sets the last-known global volume on the channel
func (cs *ChannelState[TMemory]) SetGlobalVolume(gv volume.Volume) {
	cs.RenderChannel.LastGlobalVolume = gv
	cs.RenderChannel.SetGlobalVolume(gv)
}

// SetChannelVolume sets the channel volume on the channel
func (cs *ChannelState[TMemory]) SetChannelVolume(cv volume.Volume) {
	cs.RenderChannel.ChannelVolume = cv
}

// GetChannelVolume gets the channel volume on the channel
func (cs *ChannelState[TMemory]) GetChannelVolume() volume.Volume {
	return cs.RenderChannel.ChannelVolume
}

// SetEnvelopePosition sets the envelope position for the active instrument
func (cs *ChannelState[TMemory]) SetEnvelopePosition(ticks int) {
	if nc := cs.GetVoice(); nc != nil {
		voice.SetVolumeEnvelopePosition(nc, ticks)
		voice.SetPanEnvelopePosition(nc, ticks)
		voice.SetPitchEnvelopePosition(nc, ticks)
		voice.SetFilterEnvelopePosition(nc, ticks)
	}
}

// TransitionActiveToPastState will transition the current active state to the 'past' state
// and will activate the specified New-Note Action on it
func (cs *ChannelState[TMemory]) TransitionActiveToPastState() {
	if cs.PastNotes != nil {
		switch cs.NewNoteAction {
		case note.ActionCut:
			// reset at end

		case note.ActionContinue:
			// nothing
			pn := cs.activeState.Clone()
			if nc := pn.Voice; nc != nil {
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}

		case note.ActionRelease:
			pn := cs.activeState.Clone()
			if nc := pn.Voice; nc != nil {
				nc.Release()
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}

		case note.ActionFadeout:
			pn := cs.activeState.Clone()
			if nc := pn.Voice; nc != nil {
				nc.Release()
				nc.Fadeout()
				cs.PastNotes.Add(cs.RenderChannel.ChannelNum, pn)
			}
		}
	}
	cs.activeState.Reset()
}

// DoPastNoteEffect performs an action on all past-note playbacks associated with the channel
func (cs *ChannelState[TMemory]) DoPastNoteEffect(action note.Action) {
	cs.PastNotes.Do(cs.RenderChannel.ChannelNum, action)
}

// SetNewNoteAction sets the New-Note Action on the channel
func (cs *ChannelState[TMemory]) SetNewNoteAction(nna note.Action) {
	cs.NewNoteAction = nna
}

// GetNewNoteAction gets the New-Note Action on the channel
func (cs *ChannelState[TMemory]) GetNewNoteAction() note.Action {
	return cs.NewNoteAction
}

// SetVolumeEnvelopeEnable sets the enable flag on the active volume envelope
func (cs *ChannelState[TMemory]) SetVolumeEnvelopeEnable(enabled bool) {
	voice.EnableVolumeEnvelope(cs.activeState.Voice, enabled)
}

// SetPanningEnvelopeEnable sets the enable flag on the active panning envelope
func (cs *ChannelState[TMemory]) SetPanningEnvelopeEnable(enabled bool) {
	voice.EnablePanEnvelope(cs.activeState.Voice, enabled)
}

// SetPitchEnvelopeEnable sets the enable flag on the active pitch/filter envelope
func (cs *ChannelState[TMemory]) SetPitchEnvelopeEnable(enabled bool) {
	voice.EnablePitchEnvelope(cs.activeState.Voice, enabled)
}

func (cs *ChannelState[TMemory]) NoteCut() {
	cs.activeState.Period = nil
}
