package voice

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
	"github.com/gotracker/playback/voice/types"
)

type Voice interface {
	Clone() Voice
	DumpState(ch index.Channel, t tracing.Tracer)

	// Configuration
	Reset()

	// Actions
	Attack()
	Release()
	Fadeout()
	Stop()

	// State Machine Update
	Advance()

	// General Parameters
	IsDone() bool
	GetSampleRate() period.Frequency
}

type RenderVoice[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	Voice

	// Configuration
	Setup(config InstrumentConfig[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	SetPCM(sample pcm.Sample, wholeLoop loop.Loop, sustainLoop loop.Loop, mixVolume TMixingVolume, defaultVolume TVolume)
}

type AmpModulator[TGlobalVolume, TMixingVolume, TVolume Volume] interface {
	// Amp/Volume Parameters
	IsActive() bool
	SetActive(active bool)
	GetMixingVolume() TMixingVolume
	SetMixingVolume(v TMixingVolume)
	GetVolume() TVolume
	SetVolume(v TVolume)
	GetVolumeDelta() types.VolumeDelta
	SetVolumeDelta(d types.VolumeDelta)
	GetFinalVolume() volume.Volume
}

type FadeoutModulator interface {
	IsFadeout() bool
	GetFadeoutVolume() volume.Volume
}

type FreqModulator[TPeriod Period] interface {
	// Frequency/Pitch Parameters
	GetPeriod() TPeriod
	SetPeriod(p TPeriod)
	GetPeriodDelta() period.Delta
	SetPeriodDelta(delta period.Delta)
	GetFinalPeriod() TPeriod
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

	GetFinalPeriod() TPeriod
	GetFinalVolume() volume.Volume
	GetFinalPan() panning.Position
}

type PanModulator[TPanning Panning] interface {
	// Pan Parameters
	GetPan() TPanning
	SetPan(pan TPanning)
	GetPanDelta() types.PanDelta
	SetPanDelta(d types.PanDelta)
	GetFinalPan() panning.Position
}

type PitchPanModulator[TPanning Panning] interface {
	SetPitchPanNote(st note.Semitone)
	IsPitchPanEnabled() bool
	EnablePitchPan(enabled bool)
	GetPanSeparation() float32
}

type VolumeEnvelope[TGlobalVolume, TMixingVolume, TVolume Volume] interface {
	// Amp/Volume Envelope Parameters
	IsVolumeEnvelopeEnabled() bool
	EnableVolumeEnvelope(enabled bool)
	GetVolumeEnvelopePosition() int
	SetVolumeEnvelopePosition(pos int)
	GetCurrentVolumeEnvelope() TVolume
}

type PitchEnvelope[TPeriod Period] interface {
	// Frequency/Pitch Envelope Parameters
	IsPitchEnvelopeEnabled() bool
	EnablePitchEnvelope(enabled bool)
	GetPitchEnvelopePosition() int
	SetPitchEnvelopePosition(pos int)
	GetCurrentPitchEnvelope() period.Delta
}

type PanEnvelope[TPanning Panning] interface {
	// Pan Envelope Parameters
	IsPanEnvelopeEnabled() bool
	EnablePanEnvelope(enabled bool)
	GetPanEnvelopePosition() int
	SetPanEnvelopePosition(pos int)
	GetCurrentPanEnvelope() TPanning
}

type FilterEnvelope interface {
	// Filter Envelope Parameters
	IsFilterEnvelopeEnabled() bool
	EnableFilterEnvelope(enabled bool)
	GetFilterEnvelopePosition() int
	SetFilterEnvelopePosition(pos int)
	GetCurrentFilterEnvelope() uint8
}
