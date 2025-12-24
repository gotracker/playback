package instrument

import (
	"github.com/heucuva/optional"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/types"
)

type InstrumentIntf interface {
	IsInvalid() bool
	GetSampleRate() frequency.Frequency
	SetSampleRate(sampleRate frequency.Frequency)
	GetLength() sampling.Pos
	SetFinetune(ft note.Finetune)
	GetFinetune() note.Finetune
	GetID() ID
	GetSemitoneShift() int8
	GetNewNoteAction() note.Action
	GetData() Data
	GetVoiceFilterInfo() filter.Info
	GetPluginFilterInfo() filter.Info
	IsReleaseNote(n note.Note) bool
	IsStopNote(n note.Note) bool
	GetDefaultVolumeGeneric() volume.Volume
}

// StaticValues are the static values associated with an instrument
type StaticValues[TPeriod types.Period, TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	PC                 period.PeriodConverter[TPeriod]
	Filename           string
	Name               string
	ID                 ID
	Volume             TVolume
	MixingVolume       optional.Value[TMixingVolume]
	Panning            optional.Value[TPanning]
	RelativeNoteNumber int8
	AutoVibrato        autovibrato.AutoVibratoConfig[TPeriod]
	NewNoteAction      note.Action
	Finetune           note.Finetune
	VoiceFilter        filter.Info
	PluginFilter       filter.Info
}

// Instrument is the mildly-decoded instrument/sample header
type Instrument[TPeriod types.Period, TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	Static     StaticValues[TPeriod, TMixingVolume, TVolume, TPanning]
	Inst       Data
	SampleRate frequency.Frequency
	Finetune   optional.Value[note.Finetune]
}

// IsInvalid always returns false (valid)
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) IsInvalid() bool {
	return false
}

// GetSampleRate returns the central-note sample rate value for the instrument
// This may get mutated if a finetune effect is processed
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetSampleRate() frequency.Frequency {
	return inst.SampleRate
}

// SetSampleRate sets the central-note sample rate value for the instrument
func (inst *Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) SetSampleRate(sampleRate frequency.Frequency) {
	inst.SampleRate = sampleRate
}

// GetDefaultVolume returns the default volume value for the instrument
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetDefaultVolume() TVolume {
	return inst.Static.Volume
}

// GetDefaultVolumeGeneric returns the default volume value for the instrument
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetDefaultVolumeGeneric() volume.Volume {
	return inst.Static.Volume.ToVolume()
}

// GetLength returns the length of the instrument
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetLength() sampling.Pos {
	return inst.Inst.GetLength()
}

// SetFinetune sets the finetune value on the instrument
func (inst *Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) SetFinetune(ft note.Finetune) {
	inst.Finetune.Set(ft)
}

// GetFinetune returns the finetune value on the instrument
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetFinetune() note.Finetune {
	if ft, ok := inst.Finetune.Get(); ok {
		return ft
	}
	return inst.Static.Finetune
}

// GetID returns the instrument number (1-based)
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetID() ID {
	return inst.Static.ID
}

// GetSemitoneShift returns the amount of semitones worth of shift to play the instrument at
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetSemitoneShift() int8 {
	return inst.Static.RelativeNoteNumber
}

// GetNewNoteAction returns the NewNoteAction associated to the instrument
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetNewNoteAction() note.Action {
	return inst.Static.NewNoteAction
}

// GetData returns the instrument-specific data interface
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetData() Data {
	return inst.Inst
}

// GetVoiceFilterInfo returns the factory for the channel filter
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetVoiceFilterInfo() filter.Info {
	return inst.Static.VoiceFilter
}

// GetPluginFilterInfo returns the factory for the channel plugin filter
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) GetPluginFilterInfo() filter.Info {
	return inst.Static.PluginFilter
}

// IsReleaseNote returns true if the note is a release (Note-Off)
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) IsReleaseNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		switch inst.GetData().(type) {
		case *OPL2:
			return true
		}
	}
	return note.IsRelease(n)
}

// IsStopNote returns true if the note is a stop (Note-Cut)
func (inst Instrument[TPeriod, TMixingVolume, TVolume, TPanning]) IsStopNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		switch inst.GetData().(type) {
		case *OPL2:
			return true
		}
	}
	return note.IsRelease(n)
}
