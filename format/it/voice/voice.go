package voice

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

type Period interface {
	period.Period
}

type itVoice[TPeriod Period] struct {
	ms     *settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]
	config voice.InstrumentConfig[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]

	pitchAndFilterEnvShared bool
	filterEnvActive         bool // if pitchAndFilterEnvShared is true, this dictates which is active initially - true=filter, false=pitch
	fadeoutMode             fadeout.Mode

	component.KeyModulator

	voicer      component.Voicer[TPeriod, itVolume.Volume]
	amp         component.AmpModulator[itVolume.FineVolume, itVolume.Volume]
	fadeout     component.FadeoutModulator
	freq        component.FreqModulator[TPeriod]
	autoVibrato component.AutoVibratoModulator[TPeriod]
	pan         component.PanModulator[itPanning.Panning]
	pitchPan    component.PitchPanModulator[itPanning.Panning]
	volEnv      component.VolumeEnvelope[itVolume.Volume]
	pitchEnv    component.PitchEnvelope
	panEnv      component.PanEnvelope[itPanning.Panning]
	filterEnv   component.FilterEnvelope
	vol0Opt     component.Vol0Optimization
}

var (
	_ voice.Sampler                                                                   = (*itVoice[period.Linear])(nil)
	_ voice.AmpModulator[itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume]   = (*itVoice[period.Linear])(nil)
	_ voice.FadeoutModulator                                                          = (*itVoice[period.Linear])(nil)
	_ voice.FreqModulator[period.Linear]                                              = (*itVoice[period.Linear])(nil)
	_ voice.PanModulator[itPanning.Panning]                                           = (*itVoice[period.Linear])(nil)
	_ voice.PitchPanModulator[itPanning.Panning]                                      = (*itVoice[period.Linear])(nil)
	_ voice.VolumeEnvelope[itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume] = (*itVoice[period.Linear])(nil)
	_ voice.PitchEnvelope[period.Linear]                                              = (*itVoice[period.Linear])(nil)
	_ voice.PanEnvelope[itPanning.Panning]                                            = (*itVoice[period.Linear])(nil)
	_ voice.FilterEnvelope                                                            = (*itVoice[period.Linear])(nil)
)

func New[TPeriod Period](config voice.VoiceConfig[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], ms *settings.MachineSettings[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) voice.RenderVoice[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	v := &itVoice[TPeriod]{
		ms:                      ms,
		pitchAndFilterEnvShared: true,
	}

	v.KeyModulator.Setup(component.KeyModulatorSettings{
		Attack:          v.doAttack,
		Release:         v.doRelease,
		Fadeout:         v.doFadeout,
		DeferredAttack:  v.doDeferredAttack,
		DeferredRelease: v.doDeferredRelease,
	})

	v.amp.Setup(component.AmpModulatorSettings[itVolume.FineVolume, itVolume.Volume]{
		Active:              true,
		DefaultMixingVolume: config.InitialMixing,
		DefaultVolume:       config.InitialVolume,
	})

	v.pan.Setup(component.PanModulatorSettings[itPanning.Panning]{
		Enabled:    config.PanEnabled,
		InitialPan: config.InitialPan,
	})

	v.vol0Opt.Setup(config.Vol0Optimization)

	return v
}

func (v *itVoice[TPeriod]) doAttack() {
	v.vol0Opt.Reset()
	v.autoVibrato.Reset()

	v.SetVolumeEnvelopePosition(0)
	v.SetPitchEnvelopePosition(0)
	v.SetPanEnvelopePosition(0)
	v.SetFilterEnvelopePosition(0)

	v.fadeout.Reset()
	v.volEnv.Attack()
	v.pitchEnv.Attack()
	v.panEnv.Attack()
	v.filterEnv.Attack()
	if v.voicer != nil {
		v.voicer.Attack()
	}
}

func (v *itVoice[TPeriod]) doRelease() {
	v.volEnv.Release()
	v.pitchEnv.Release()
	v.panEnv.Release()
	v.filterEnv.Release()
	if v.voicer != nil {
		v.voicer.Release()
	}
}

func (v *itVoice[TPeriod]) doFadeout() {
	if v.voicer != nil {
		v.voicer.Fadeout()
	}

	v.fadeout.Fadeout()
}

func (v *itVoice[TPeriod]) doDeferredAttack() {
	if v.voicer != nil {
		v.voicer.DeferredAttack()
	}
}

func (v *itVoice[TPeriod]) doDeferredRelease() {
	if v.voicer != nil {
		v.voicer.DeferredRelease()
	}
}

func (v itVoice[TPeriod]) getFadeoutEnabled() bool {
	return v.fadeoutMode.IsFadeoutActive(v.IsKeyFadeout(), v.volEnv.IsEnabled(), v.volEnv.IsDone())
}

func (v *itVoice[TPeriod]) Setup(config voice.InstrumentConfig[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) {
	v.config = config
	v.fadeout.Setup(component.FadeoutModulatorSettings{
		GetEnabled: v.getFadeoutEnabled,
		Amount:     config.FadeOut.Amount,
	})
	v.freq.Setup(component.FreqModulatorSettings[TPeriod]{})
	v.autoVibrato.Setup(config.AutoVibrato)
	v.pitchPan.Setup(component.PitchPanModulatorSettings[itPanning.Panning]{
		PitchPanEnable:     config.PitchPan.Enabled,
		PitchPanCenter:     config.PitchPan.Center,
		PitchPanSeparation: config.PitchPan.Separation,
	})
	volEnvSettings := component.EnvelopeSettings[itVolume.Volume, itVolume.Volume]{
		Envelope: config.VolEnv,
	}
	if config.VolEnvFinishFadesOut {
		volEnvSettings.OnFinished = func(v voice.Voice) {
			v.Fadeout()
		}
	}
	v.volEnv.Setup(volEnvSettings)
	v.pitchEnv.Setup(component.EnvelopeSettings[int8, period.Delta]{
		Envelope: config.PitchFiltEnv,
	})
	v.panEnv.Setup(component.EnvelopeSettings[itPanning.Panning, itPanning.Panning]{
		Envelope: config.PanEnv,
	})
	v.filterEnv.Setup(component.EnvelopeSettings[int8, uint8]{
		Envelope: config.PitchFiltEnv,
	})
	v.KeyModulator.Release()
	v.Reset()
}

func (v *itVoice[TPeriod]) Reset() {
	v.filterEnvActive = v.config.PitchFiltMode
	v.fadeoutMode = v.config.FadeOut.Mode

	v.fadeout.Reset()

	v.volEnv.Reset()
	v.pitchEnv.Reset()
	v.panEnv.Reset()
	v.filterEnv.Reset()

	v.autoVibrato.Reset()

	v.volEnv.Reset()
	v.pitchEnv.Reset()
	v.panEnv.Reset()
	v.filterEnv.Reset()
	v.vol0Opt.Reset()
}

func (v *itVoice[TPeriod]) Stop() {
	v.voicer = nil
}

func (v *itVoice[TPeriod]) IsDone() bool {
	if v.voicer == nil {
		return true
	}

	if v.fadeout.IsActive() {
		return v.fadeout.GetVolume() <= 0
	}

	return v.vol0Opt.IsDone()
}

func (v *itVoice[TPeriod]) Advance() {
	v.fadeout.Advance()
	v.autoVibrato.Advance()
	v.pitchPan.Advance()
	if v.IsVolumeEnvelopeEnabled() {
		if doneCB := v.volEnv.Advance(); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsPanEnvelopeEnabled() {
		if doneCB := v.panEnv.Advance(); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsPitchEnvelopeEnabled() {
		if doneCB := v.pitchEnv.Advance(); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsFilterEnvelopeEnabled() {
		if doneCB := v.filterEnv.Advance(); doneCB != nil {
			doneCB(v)
		}
	}

	if v.config.VoiceFilter != nil && v.IsFilterEnvelopeEnabled() {
		fval := v.GetCurrentFilterEnvelope()
		v.config.VoiceFilter.UpdateEnv(fval)
	}

	// has to be after the mod/env updates
	v.KeyModulator.DeferredUpdate()

	v.vol0Opt.ObserveVolume(v.GetFinalVolume())
	v.KeyModulator.Advance()
}

func (v *itVoice[TPeriod]) Clone() voice.Voice {
	vv := itVoice[TPeriod]{
		ms:                      v.ms,
		config:                  v.config,
		pitchAndFilterEnvShared: v.pitchAndFilterEnvShared,
		filterEnvActive:         v.filterEnvActive,
		fadeoutMode:             v.fadeoutMode,
		amp:                     v.amp.Clone(),
		freq:                    v.freq.Clone(),
		pan:                     v.pan.Clone(),
		volEnv:                  v.volEnv.Clone(),
		pitchEnv:                v.pitchEnv.Clone(),
		panEnv:                  v.panEnv.Clone(),
		filterEnv:               v.filterEnv.Clone(),
		vol0Opt:                 v.vol0Opt.Clone(),
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

	if v.config.PluginFilter != nil {
		vv.config.PluginFilter = v.config.PluginFilter.Clone()
	}

	return &vv
}

func (v *itVoice[TPeriod]) SetPCM(samp pcm.Sample, wholeLoop, sustainLoop loop.Loop, defVol itVolume.Volume) {
	var s component.Sampler[TPeriod, itVolume.Volume]
	s.Setup(component.SamplerSettings[TPeriod, itVolume.Volume]{
		Sample:        samp,
		DefaultVolume: defVol,
		WholeLoop:     wholeLoop,
		SustainLoop:   sustainLoop,
	})
	v.voicer = &s
}
