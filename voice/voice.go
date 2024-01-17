package voice

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/opl2"
	"github.com/gotracker/playback/voice/types"
)

type Voice interface {
	Clone(background bool) Voice
	DumpState(ch index.Channel, t tracing.Tracer)

	// Configuration
	Reset() error
	SetOPL2Chip(chip opl2.Chip)

	// Actions
	Attack()
	Release()
	Fadeout()
	Stop()

	// State Machine Update
	Tick() error
	RowEnd() error

	// General Parameters
	IsDone() bool
	SetMuted(muted bool) error
	IsMuted() bool
	GetSampleRate() frequency.Frequency
}

type RenderVoice[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	Voice

	// Configuration
	Setup(inst *instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning], outputRate frequency.Frequency) error
}

type AmpModulator[TGlobalVolume, TMixingVolume, TVolume Volume] interface {
	// Amp/Volume Parameters
	IsActive() bool
	SetActive(active bool) error
	GetMixingVolume() TMixingVolume
	SetMixingVolume(v TMixingVolume) error
	GetVolume() TVolume
	SetVolume(v TVolume) error
	GetVolumeDelta() types.VolumeDelta
	SetVolumeDelta(d types.VolumeDelta) error
	GetFinalVolume() volume.Volume
}

type FadeoutModulator interface {
	IsFadeout() bool
	GetFadeoutVolume() volume.Volume
}

type FreqModulator[TPeriod Period] interface {
	// Frequency/Pitch Parameters
	GetPeriod() TPeriod
	SetPeriod(p TPeriod) error
	GetPeriodDelta() period.Delta
	SetPeriodDelta(delta period.Delta) error
	GetFinalPeriod() (TPeriod, error)
}

type Sampler interface {
	// Sampler Parameters
	SetPos(pos sampling.Pos) error
	GetPos() (sampling.Pos, error)
}

type RenderSampler[TPeriod Period] interface {
	Voice
	sampling.SampleStream

	IsActive() bool

	SetPos(pos sampling.Pos) error
	GetPos() (sampling.Pos, error)

	GetFinalPeriod() (TPeriod, error)
	GetFinalVolume() volume.Volume
	GetFinalPan() panning.Position
}

type PanModulator[TPanning Panning] interface {
	// Pan Parameters
	GetPan() TPanning
	SetPan(pan TPanning) error
	GetPanDelta() types.PanDelta
	SetPanDelta(d types.PanDelta) error
	GetFinalPan() panning.Position
}

type PitchPanModulator[TPanning Panning] interface {
	SetPitchPanNote(st note.Semitone) error
	IsPitchPanEnabled() bool
	EnablePitchPan(enabled bool) error
	GetPanSeparation() float32
}

type VolumeEnvelope[TGlobalVolume, TMixingVolume, TVolume Volume] interface {
	// Amp/Volume Envelope Parameters
	IsVolumeEnvelopeEnabled() bool
	EnableVolumeEnvelope(enabled bool) error
	GetVolumeEnvelopePosition() int
	SetVolumeEnvelopePosition(pos int) error
	GetCurrentVolumeEnvelope() TVolume
}

type PitchEnvelope[TPeriod Period] interface {
	// Frequency/Pitch Envelope Parameters
	IsPitchEnvelopeEnabled() bool
	EnablePitchEnvelope(enabled bool) error
	GetPitchEnvelopePosition() int
	SetPitchEnvelopePosition(pos int) error
	GetCurrentPitchEnvelope() period.Delta
}

type PanEnvelope[TPanning Panning] interface {
	// Pan Envelope Parameters
	IsPanEnvelopeEnabled() bool
	EnablePanEnvelope(enabled bool) error
	GetPanEnvelopePosition() int
	SetPanEnvelopePosition(pos int) error
	GetCurrentPanEnvelope() TPanning
}

type FilterEnvelope interface {
	// Filter Envelope Parameters
	IsFilterEnvelopeEnabled() bool
	EnableFilterEnvelope(enabled bool) error
	GetFilterEnvelopePosition() int
	SetFilterEnvelopePosition(pos int) error
	GetCurrentFilterEnvelope() uint8
}
