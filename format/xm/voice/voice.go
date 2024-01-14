package voice

import (
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

type Period interface {
	period.Period
}

type xmVoice[TPeriod Period] struct {
	config voice.InstrumentConfig[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]

	fadeoutMode fadeout.Mode

	component.KeyModulator

	voicer      component.Voicer[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume]
	amp         component.AmpModulator[xmVolume.XmVolume, xmVolume.XmVolume]
	fadeout     component.FadeoutModulator
	freq        component.FreqModulator[TPeriod]
	autoVibrato component.AutoVibratoModulator[TPeriod]
	pan         component.PanModulator[xmPanning.Panning]
	volEnv      component.VolumeEnvelope[xmVolume.XmVolume]
	panEnv      component.PanEnvelope[xmPanning.Panning]
	vol0Opt     component.Vol0Optimization
}

var (
	_ voice.Sampler                                                                 = (*xmVoice[period.Linear])(nil)
	_ voice.AmpModulator[xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume]   = (*xmVoice[period.Linear])(nil)
	_ voice.FadeoutModulator                                                        = (*xmVoice[period.Linear])(nil)
	_ voice.FreqModulator[period.Linear]                                            = (*xmVoice[period.Linear])(nil)
	_ voice.PanModulator[xmPanning.Panning]                                         = (*xmVoice[period.Linear])(nil)
	_ voice.VolumeEnvelope[xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume] = (*xmVoice[period.Linear])(nil)
	_ voice.PanEnvelope[xmPanning.Panning]                                          = (*xmVoice[period.Linear])(nil)
)

func New[TPeriod Period](config voice.VoiceConfig[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) voice.RenderVoice[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	v := &xmVoice[TPeriod]{}

	v.KeyModulator.Setup(component.KeyModulatorSettings{
		Attack:          v.doAttack,
		Release:         v.doRelease,
		Fadeout:         v.doFadeout,
		DeferredAttack:  v.doDeferredAttack,
		DeferredRelease: v.doDeferredRelease,
	})

	v.amp.Setup(component.AmpModulatorSettings[xmVolume.XmVolume, xmVolume.XmVolume]{
		Active:              true,
		DefaultMixingVolume: config.InitialMixing,
		DefaultVolume:       config.InitialVolume,
	})

	v.pan.Setup(component.PanModulatorSettings[xmPanning.Panning]{
		Enabled:    config.PanEnabled,
		InitialPan: config.InitialPan,
	})

	v.vol0Opt.Setup(config.Vol0Optimization)

	return v
}

func (v *xmVoice[TPeriod]) doAttack() {
	v.vol0Opt.Reset()
	v.autoVibrato.ResetAutoVibrato()

	v.SetVolumeEnvelopePosition(0)
	v.SetPitchEnvelopePosition(0)
	v.SetPanEnvelopePosition(0)
	v.SetFilterEnvelopePosition(0)

	v.fadeout.Reset()
	v.volEnv.Attack()
	v.panEnv.Attack()
	if v.voicer != nil {
		v.voicer.Attack()
	}
}

func (v *xmVoice[TPeriod]) doRelease() {
	v.volEnv.Release()
	v.panEnv.Release()
	if v.voicer != nil {
		v.voicer.Release()
	}
}

func (v *xmVoice[TPeriod]) doFadeout() {
	if v.voicer != nil {
		v.voicer.Fadeout()
	}
}

func (v *xmVoice[TPeriod]) doDeferredAttack() {
	if v.voicer != nil {
		v.voicer.DeferredAttack()
	}
}

func (v *xmVoice[TPeriod]) doDeferredRelease() {
	if v.voicer != nil {
		v.voicer.DeferredRelease()
	}
}

func (v xmVoice[TPeriod]) getFadeoutEnabled() bool {
	return v.fadeoutMode.IsFadeoutActive(v.IsKeyFadeout(), v.volEnv.IsEnabled(), v.volEnv.IsDone())
}

func (v *xmVoice[TPeriod]) Setup(config voice.InstrumentConfig[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) {
	v.config = config
	v.fadeout.Setup(component.FadeoutModulatorSettings{
		Enabled:   v.config.FadeOut.Mode != fadeout.ModeDisabled,
		GetActive: v.getFadeoutEnabled,
		Amount:    config.FadeOut.Amount,
	})
	v.freq.Setup(component.FreqModulatorSettings[TPeriod]{})
	v.autoVibrato.Setup(config.AutoVibrato)
	volEnvSettings := component.EnvelopeSettings[xmVolume.XmVolume, xmVolume.XmVolume]{
		Envelope: config.VolEnv,
	}
	if config.VolEnvFinishFadesOut {
		volEnvSettings.OnFinished = func(v voice.Voice) {
			v.Fadeout()
		}
	}
	v.volEnv.Setup(volEnvSettings)
	v.panEnv.Setup(component.EnvelopeSettings[xmPanning.Panning, xmPanning.Panning]{
		Envelope: config.PanEnv,
	})
	v.KeyModulator.Release()
	v.Reset()
}

func (v *xmVoice[TPeriod]) Reset() {
	v.fadeoutMode = v.config.FadeOut.Mode

	v.fadeout.Reset()

	v.volEnv.Reset()
	v.panEnv.Reset()

	v.autoVibrato.Reset()

	v.volEnv.Reset()
	v.panEnv.Reset()
	v.vol0Opt.Reset()
}

func (v *xmVoice[TPeriod]) Stop() {
	v.voicer = nil
}

func (v *xmVoice[TPeriod]) IsDone() bool {
	if v.voicer == nil {
		return true
	}

	if v.fadeout.IsActive() {
		return v.fadeout.GetVolume() <= 0
	}

	return v.vol0Opt.IsDone()
}

func (v *xmVoice[TPeriod]) Advance() {
	v.fadeout.Advance()
	v.autoVibrato.Advance()
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

	// has to be after the mod/env updates
	v.KeyModulator.DeferredUpdate()

	v.vol0Opt.ObserveVolume(v.GetFinalVolume())
	v.KeyModulator.Advance()
}

func (v *xmVoice[TPeriod]) Clone() voice.Voice {
	vv := xmVoice[TPeriod]{
		config:      v.config,
		fadeoutMode: v.fadeoutMode,
		amp:         v.amp.Clone(),
		fadeout:     v.fadeout.Clone(),
		freq:        v.freq.Clone(),
		autoVibrato: v.autoVibrato.Clone(),
		pan:         v.pan.Clone(),
		panEnv:      v.panEnv.Clone(nil),
		vol0Opt:     v.vol0Opt.Clone(),
	}

	vv.volEnv = v.volEnv.Clone(func(v voice.Voice) {
		vv.Fadeout()
	})

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

func (v *xmVoice[TPeriod]) SetPCM(samp pcm.Sample, wholeLoop, sustainLoop loop.Loop, mixVol, defVol xmVolume.XmVolume) {
	var s component.Sampler[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume]
	s.Setup(component.SamplerSettings[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume]{
		Sample:        samp,
		DefaultVolume: defVol,
		MixVolume:     mixVol,
		WholeLoop:     wholeLoop,
		SustainLoop:   sustainLoop,
	})
	v.voicer = &s
}
