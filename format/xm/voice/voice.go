package voice

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback/filter"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/opl2"
	"github.com/gotracker/playback/voice/pcm"
)

type Period interface {
	period.Period
}

type xmVoice[TPeriod Period] struct {
	inst *instrument.Instrument[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]

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
	voiceFilter filter.Filter
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

	v.freq.Setup(component.FreqModulatorSettings[TPeriod]{
		PC: config.PC,
	})

	v.pan.Setup(component.PanModulatorSettings[xmPanning.Panning]{
		Enabled:    config.PanEnabled,
		InitialPan: config.InitialPan,
	})

	v.vol0Opt.Setup(config.Vol0Optimization)

	return v
}

func (v *xmVoice[TPeriod]) SetOPL2Chip(chip opl2.Chip) {
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

func (v *xmVoice[TPeriod]) Setup(inst *instrument.Instrument[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], outputRate frequency.Frequency) error {
	v.inst = inst

	switch d := inst.GetData().(type) {
	case *instrument.PCM[xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]:
		v.fadeoutMode = d.FadeOut.Mode

		v.fadeout.Setup(component.FadeoutModulatorSettings{
			Enabled:   d.FadeOut.Mode != fadeout.ModeDisabled,
			GetActive: v.getFadeoutEnabled,
			Amount:    d.FadeOut.Amount,
		})

		volEnvSettings := component.EnvelopeSettings[xmVolume.XmVolume, xmVolume.XmVolume]{
			Envelope: d.VolEnv,
		}
		if d.VolEnvFinishFadesOut {
			volEnvSettings.OnFinished = func(v voice.Voice) {
				v.Fadeout()
			}
		}
		v.volEnv.Setup(volEnvSettings)

		v.panEnv.Setup(component.EnvelopeSettings[xmPanning.Panning, xmPanning.Panning]{
			Envelope: d.PanEnv,
		})

		v.amp.SetMixingVolumeOverride(d.MixingVolume)

		v.setupPCM(d.Sample, d.Loop, d.SustainLoop, xmVolume.DefaultXmMixingVolume, inst.GetDefaultVolume())

	default:
		return fmt.Errorf("unhandled instrument type: %T", d)
	}
	if inst == nil {
		return errors.New("instrument is nil")
	}

	v.autoVibrato.Setup(inst.Static.AutoVibrato)

	if factory := inst.GetFilterFactory(); factory != nil {
		v.voiceFilter = factory(inst.SampleRate)
		v.voiceFilter.SetPlaybackRate(outputRate)
	} else {
		v.voiceFilter = nil
	}

	v.Reset()
	return nil
}

func (v *xmVoice[TPeriod]) Reset() error {
	v.KeyModulator.Release()
	return errors.Join(
		v.amp.Reset(),
		v.fadeout.Reset(),
		v.freq.Reset(),
		v.autoVibrato.Reset(),
		v.pan.Reset(),
		v.volEnv.Reset(),
		v.panEnv.Reset(),
		v.vol0Opt.Reset(),
	)
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

func (v *xmVoice[TPeriod]) Tick() error {
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

	v.KeyModulator.Advance()
	return nil
}

func (v *xmVoice[TPeriod]) RowEnd() error {
	v.vol0Opt.ObserveVolume(v.GetFinalVolume())
	return nil
}

func (v *xmVoice[TPeriod]) Clone(bool) voice.Voice {
	vv := xmVoice[TPeriod]{
		inst:        v.inst,
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

	if v.voiceFilter != nil {
		vv.voiceFilter = v.voiceFilter.Clone()
	}

	return &vv
}

func (v *xmVoice[TPeriod]) setupPCM(samp pcm.Sample, wholeLoop, sustainLoop loop.Loop, mixVol, defVol xmVolume.XmVolume) {
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
