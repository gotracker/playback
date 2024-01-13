package instrument

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/heucuva/optional"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/types"
)

type InstrumentIntf interface {
	IsInvalid() bool
	GetSampleRate() period.Frequency
	SetSampleRate(sampleRate period.Frequency)
	GetLength() sampling.Pos
	SetFinetune(ft note.Finetune)
	GetFinetune() note.Finetune
	GetID() ID
	GetSemitoneShift() int8
	GetKind() Kind
	GetNewNoteAction() note.Action
	GetData() Data
	GetFilterFactory() filter.Factory
	GetPluginFilterFactory() filter.Factory
	GetAutoVibrato() autovibrato.AutoVibratoSettings
	IsReleaseNote(n note.Note) bool
	IsStopNote(n note.Note) bool
	GetDefaultVolumeGeneric() volume.Volume
}

// StaticValues are the static values associated with an instrument
type StaticValues[TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	Filename           string
	Name               string
	ID                 ID
	Volume             TVolume
	MixingVolume       optional.Value[TMixingVolume]
	Panning            optional.Value[TPanning]
	RelativeNoteNumber int8
	AutoVibrato        autovibrato.AutoVibratoSettings
	NewNoteAction      note.Action
	Finetune           note.Finetune
	FilterFactory      filter.Factory
	PluginFilter       filter.Factory
}

// Instrument is the mildly-decoded instrument/sample header
type Instrument[TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	Static     StaticValues[TMixingVolume, TVolume, TPanning]
	Inst       Data
	SampleRate period.Frequency
	Finetune   optional.Value[note.Finetune]
}

// IsInvalid always returns false (valid)
func (inst Instrument[TMixingVolume, TVolume, TPanning]) IsInvalid() bool {
	return false
}

// GetSampleRate returns the central-note sample rate value for the instrument
// This may get mutated if a finetune effect is processed
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetSampleRate() period.Frequency {
	return inst.SampleRate
}

// SetSampleRate sets the central-note sample rate value for the instrument
func (inst *Instrument[TMixingVolume, TVolume, TPanning]) SetSampleRate(sampleRate period.Frequency) {
	inst.SampleRate = sampleRate
}

// GetDefaultVolume returns the default volume value for the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetDefaultVolume() TVolume {
	return inst.Static.Volume
}

// GetDefaultVolumeGeneric returns the default volume value for the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetDefaultVolumeGeneric() volume.Volume {
	return inst.Static.Volume.ToVolume()
}

// GetLength returns the length of the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetLength() sampling.Pos {
	return inst.Inst.GetLength()
}

// SetFinetune sets the finetune value on the instrument
func (inst *Instrument[TMixingVolume, TVolume, TPanning]) SetFinetune(ft note.Finetune) {
	inst.Finetune.Set(ft)
}

// GetFinetune returns the finetune value on the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetFinetune() note.Finetune {
	if ft, ok := inst.Finetune.Get(); ok {
		return ft
	}
	return inst.Static.Finetune
}

// GetID returns the instrument number (1-based)
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetID() ID {
	return inst.Static.ID
}

// GetSemitoneShift returns the amount of semitones worth of shift to play the instrument at
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetSemitoneShift() int8 {
	return inst.Static.RelativeNoteNumber
}

// GetKind returns the kind of the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetKind() Kind {
	return inst.Inst.GetKind()
}

// GetNewNoteAction returns the NewNoteAction associated to the instrument
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetNewNoteAction() note.Action {
	return inst.Static.NewNoteAction
}

// GetData returns the instrument-specific data interface
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetData() Data {
	return inst.Inst
}

// GetFilterFactory returns the factory for the channel filter
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetFilterFactory() filter.Factory {
	return inst.Static.FilterFactory
}

// GetPluginFilterFactory returns the factory for the channel plugin filter
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetPluginFilterFactory() filter.Factory {
	return inst.Static.PluginFilter
}

// GetAutoVibrato returns the settings for the autovibrato system
func (inst Instrument[TMixingVolume, TVolume, TPanning]) GetAutoVibrato() autovibrato.AutoVibratoSettings {
	return inst.Static.AutoVibrato
}

// IsReleaseNote returns true if the note is a release (Note-Off)
func (inst Instrument[TMixingVolume, TVolume, TPanning]) IsReleaseNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		return inst.GetKind() == KindOPL2
	}
	return note.IsRelease(n)
}

// IsStopNote returns true if the note is a stop (Note-Cut)
func (inst Instrument[TMixingVolume, TVolume, TPanning]) IsStopNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		return inst.GetKind() == KindPCM
	}
	return note.IsRelease(n)
}
