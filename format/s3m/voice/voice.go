package voice

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback/filter"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

type s3mVoice struct {
	inst *instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]

	component.KeyModulator

	voicer component.Voicer[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]
	component.AmpModulator[s3mVolume.FineVolume, s3mVolume.Volume]
	component.FreqModulator[period.Amiga]
	component.PanModulator[s3mPanning.Panning]
	vol0Opt     component.Vol0Optimization
	voiceFilter filter.Filter
}

var (
	_ voice.Sampler                                                                = (*s3mVoice)(nil)
	_ voice.AmpModulator[s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume] = (*s3mVoice)(nil)
	_ voice.FreqModulator[period.Amiga]                                            = (*s3mVoice)(nil)
	_ voice.PanModulator[s3mPanning.Panning]                                       = (*s3mVoice)(nil)
)

func New(config voice.VoiceConfig[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) voice.RenderVoice[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	v := &s3mVoice{}

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

	v.FreqModulator.Setup(component.FreqModulatorSettings[period.Amiga]{})

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

func (v *s3mVoice) Setup(inst *instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], outputRate frequency.Frequency) error {
	v.inst = inst

	switch d := inst.GetData().(type) {
	case *instrument.PCM[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]:
		v.AmpModulator.SetMixingVolumeOverride(d.MixingVolume)

		v.setupPCM(d.Sample, d.Loop, d.SustainLoop, s3mVolume.MaxFineVolume, inst.GetDefaultVolume())

	default:
		return fmt.Errorf("unhandled instrument type: %T", inst)
	}
	if inst == nil {
		return errors.New("instrument is nil")
	}

	if factory := inst.GetFilterFactory(); factory != nil {
		v.voiceFilter = factory(inst.SampleRate)
		v.voiceFilter.SetPlaybackRate(outputRate)
	} else {
		v.voiceFilter = nil
	}

	v.Reset()
	return nil
}

func (v *s3mVoice) Reset() {
	v.AmpModulator.Reset()
	v.FreqModulator.Reset()
	v.PanModulator.Reset()
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

func (v *s3mVoice) Clone(bool) voice.Voice {
	vv := s3mVoice{
		inst:          v.inst,
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

	if v.voiceFilter != nil {
		vv.voiceFilter = v.voiceFilter.Clone()
	}

	return &vv
}

func (v *s3mVoice) setupPCM(samp pcm.Sample, wholeLoop, sustainLoop loop.Loop, mixVol s3mVolume.FineVolume, defVol s3mVolume.Volume) {
	var s component.Sampler[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]
	s.Setup(component.SamplerSettings[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]{
		Sample:        samp,
		DefaultVolume: defVol,
		MixVolume:     mixVol,
		WholeLoop:     wholeLoop,
		SustainLoop:   sustainLoop,
	})
	v.voicer = &s
}
