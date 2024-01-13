package voice

import (
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

type Period interface {
	period.Period
}

type s3mVoice struct {
	ms     *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]
	config voice.InstrumentConfig[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]

	component.KeyModulator

	voicer component.Voicer[period.Amiga, s3mVolume.Volume]
	component.AmpModulator[s3mVolume.FineVolume, s3mVolume.Volume]
	component.FreqModulator[period.Amiga]
	component.PanModulator[s3mPanning.Panning]
	vol0Opt component.Vol0Optimization
}

var (
	_ voice.Sampler                                                                = (*s3mVoice)(nil)
	_ voice.AmpModulator[s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume] = (*s3mVoice)(nil)
	_ voice.FreqModulator[period.Amiga]                                            = (*s3mVoice)(nil)
	_ voice.PanModulator[s3mPanning.Panning]                                       = (*s3mVoice)(nil)
)

func New(config voice.VoiceConfig[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], ms *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) voice.RenderVoice[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	v := &s3mVoice{
		ms: ms,
	}

	v.KeyModulator.Setup(component.KeyModulatorSettings{
		Attack:          v.doAttack,
		Release:         v.doRelease,
		Fadeout:         v.doFadeout,
		DeferredAttack:  v.doDeferredAttack,
		DeferredRelease: v.doDeferredRelease,
	})

	v.AmpModulator.Setup(component.AmpModulatorSettings[s3mVolume.FineVolume, s3mVolume.Volume]{
		Active:              true,
		DefaultMixingVolume: config.InitialMixing,
		DefaultVolume:       config.InitialVolume,
	})

	v.PanModulator.Setup(component.PanModulatorSettings[s3mPanning.Panning]{
		Enabled:    config.PanEnabled,
		InitialPan: config.InitialPan,
	})

	v.vol0Opt.Setup(config.Vol0Optimization)
	return v
}

func (v *s3mVoice) doAttack() {
	v.vol0Opt.Reset()

	if v.voicer != nil {
		v.voicer.Attack()
	}
}

func (v *s3mVoice) doRelease() {
	if v.voicer != nil {
		v.voicer.Release()
	}
}

func (v *s3mVoice) doFadeout() {
}

func (v *s3mVoice) doDeferredAttack() {
	if v.voicer != nil {
		v.voicer.DeferredAttack()
	}
}

func (v *s3mVoice) doDeferredRelease() {
	if v.voicer != nil {
		v.voicer.DeferredRelease()
	}
}

func (v *s3mVoice) Setup(config voice.InstrumentConfig[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) {
	v.config = config

	v.FreqModulator.Setup(component.FreqModulatorSettings[period.Amiga]{})
	v.KeyModulator.Release()
	v.Reset()
}

func (v *s3mVoice) Reset() {
	v.vol0Opt.Reset()
}

func (v *s3mVoice) Stop() {
	v.voicer = nil
}

func (v *s3mVoice) IsDone() bool {
	if v.voicer == nil {
		return true
	}

	return v.vol0Opt.IsDone()
}

func (v *s3mVoice) Advance() {
	// has to be after the mod/env updates
	v.KeyModulator.DeferredUpdate()

	v.vol0Opt.ObserveVolume(v.GetFinalVolume())
	v.KeyModulator.Advance()
}

func (v *s3mVoice) Clone() voice.Voice {
	vv := s3mVoice{
		ms:            v.ms,
		config:        v.config,
		AmpModulator:  v.AmpModulator.Clone(),
		FreqModulator: v.FreqModulator.Clone(),
		PanModulator:  v.PanModulator.Clone(),
		vol0Opt:       v.vol0Opt.Clone(),
	}

	vv.KeyModulator = v.KeyModulator.Clone(component.KeyModulatorSettings{
		Attack:          vv.doAttack,
		Release:         vv.doRelease,
		Fadeout:         vv.doFadeout,
		DeferredAttack:  vv.doDeferredAttack,
		DeferredRelease: vv.doDeferredRelease,
	})

	if v.voicer != nil {
		vv.voicer = v.voicer.Clone()
	}

	if v.config.VoiceFilter != nil {
		vv.config.VoiceFilter = v.config.VoiceFilter.Clone()
	}

	return &vv
}

func (v *s3mVoice) SetPCM(samp pcm.Sample, wholeLoop, sustainLoop loop.Loop, defVol s3mVolume.Volume) {
	var s component.Sampler[period.Amiga, s3mVolume.Volume]
	s.Setup(component.SamplerSettings[period.Amiga, s3mVolume.Volume]{
		Sample:        samp,
		DefaultVolume: defVol,
		WholeLoop:     wholeLoop,
		SustainLoop:   sustainLoop,
	})
	v.voicer = &s
}
