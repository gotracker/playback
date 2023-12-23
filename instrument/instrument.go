package instrument

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/heucuva/optional"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
)

// StaticValues are the static values associated with an instrument
type StaticValues struct {
	Filename           string
	Name               string
	ID                 ID
	Volume             volume.Volume
	RelativeNoteNumber int8
	AutoVibrato        voice.AutoVibrato
	NewNoteAction      note.Action
	Finetune           note.Finetune
	FilterFactory      filter.Factory
	PluginFilter       filter.Factory
}

// Instrument is the mildly-decoded instrument/sample header
type Instrument struct {
	Static     StaticValues
	Inst       Data
	SampleRate period.Frequency
	Finetune   optional.Value[note.Finetune]
}

// IsInvalid always returns false (valid)
func (inst Instrument) IsInvalid() bool {
	return false
}

// GetSampleRate returns the central-note sample rate value for the instrument
// This may get mutated if a finetune effect is processed
func (inst Instrument) GetSampleRate() period.Frequency {
	return inst.SampleRate
}

// SetSampleRate sets the central-note sample rate value for the instrument
func (inst *Instrument) SetSampleRate(sampleRate period.Frequency) {
	inst.SampleRate = sampleRate
}

// GetDefaultVolume returns the default volume value for the instrument
func (inst Instrument) GetDefaultVolume() volume.Volume {
	return inst.Static.Volume
}

// GetLength returns the length of the instrument
func (inst Instrument) GetLength() sampling.Pos {
	return inst.Inst.GetLength()
}

// SetFinetune sets the finetune value on the instrument
func (inst *Instrument) SetFinetune(ft note.Finetune) {
	inst.Finetune.Set(ft)
}

// GetFinetune returns the finetune value on the instrument
func (inst Instrument) GetFinetune() note.Finetune {
	if ft, ok := inst.Finetune.Get(); ok {
		return ft
	}
	return inst.Static.Finetune
}

// GetID returns the instrument number (1-based)
func (inst Instrument) GetID() ID {
	return inst.Static.ID
}

// GetSemitoneShift returns the amount of semitones worth of shift to play the instrument at
func (inst Instrument) GetSemitoneShift() int8 {
	return inst.Static.RelativeNoteNumber
}

// GetKind returns the kind of the instrument
func (inst Instrument) GetKind() Kind {
	return inst.Inst.GetKind()
}

// GetNewNoteAction returns the NewNoteAction associated to the instrument
func (inst Instrument) GetNewNoteAction() note.Action {
	return inst.Static.NewNoteAction
}

// GetData returns the instrument-specific data interface
func (inst Instrument) GetData() Data {
	return inst.Inst
}

// GetFilterFactory returns the factory for the channel filter
func (inst Instrument) GetFilterFactory() filter.Factory {
	return inst.Static.FilterFactory
}

// GetPluginFilterFactory returns the factory for the channel plugin filter
func (inst Instrument) GetPluginFilterFactory() filter.Factory {
	return inst.Static.PluginFilter
}

// GetAutoVibrato returns the settings for the autovibrato system
func (inst Instrument) GetAutoVibrato() voice.AutoVibrato {
	return inst.Static.AutoVibrato
}

// IsReleaseNote returns true if the note is a release (Note-Off)
func (inst Instrument) IsReleaseNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		return inst.GetKind() == KindOPL2
	}
	return note.IsRelease(n)
}

// IsStopNote returns true if the note is a stop (Note-Cut)
func (inst Instrument) IsStopNote(n note.Note) bool {
	switch n.Type() {
	case note.SpecialTypeStopOrRelease:
		return inst.GetKind() == KindPCM
	}
	return note.IsRelease(n)
}
