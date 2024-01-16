package voice

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback/filter"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mSystem "github.com/gotracker/playback/format/s3m/system"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/opl2"
)

type s3mVoice struct {
	inst        *instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]
	opl2Chip    opl2.Chip
	opl2Channel index.OPLChannel

	component.KeyModulator

	voicer component.Voicer[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]
	component.AmpModulator[s3mVolume.FineVolume, s3mVolume.Volume]
	component.FreqModulator[period.Amiga]
	component.PanModulator[s3mPanning.Panning]
	opl2        component.OPL2Registers
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
	v := &s3mVoice{
		opl2Chip:    config.OPLChip,
		opl2Channel: config.OPLChannel,
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

	v.FreqModulator.Setup(component.FreqModulatorSettings[period.Amiga]{
		PC: config.PC,
	})

	v.PanModulator.Setup(component.PanModulatorSettings[s3mPanning.Panning]{
		Enabled:    config.PanEnabled,
		InitialPan: config.InitialPan,
	})

	v.vol0Opt.Setup(config.Vol0Optimization)

	return v
}

func (v *s3mVoice) SetOPL2Chip(chip opl2.Chip) {
	v.opl2Chip = chip
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

func (v *s3mVoice) Setup(inst *instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], outputRate frequency.Frequency) error {
	v.inst = inst

	switch d := inst.GetData().(type) {
	case *instrument.PCM[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]:
		v.AmpModulator.SetMixingVolumeOverride(d.MixingVolume)

		var s component.Sampler[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]
		s.Setup(component.SamplerSettings[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]{
			Sample:        d.Sample,
			DefaultVolume: inst.GetDefaultVolume(),
			MixVolume:     s3mVolume.MaxFineVolume,
			WholeLoop:     d.Loop,
			SustainLoop:   d.SustainLoop,
		})
		v.voicer = &s

	case *instrument.OPL2:
		var o component.OPL2[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]
		o.Setup(v.opl2Chip, int(v.opl2Channel), v.opl2, s3mPeriod.S3MAmigaConverter, s3mSystem.S3MBaseClock, inst.GetDefaultVolume())
		v.voicer = &o

	default:
		return fmt.Errorf("unhandled instrument type: %T", d)
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

func (v *s3mVoice) Reset() error {
	return errors.Join(
		v.AmpModulator.Reset(),
		v.FreqModulator.Reset(),
		v.PanModulator.Reset(),
		v.vol0Opt.Reset(),
	)
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

func (v *s3mVoice) Tick() error {
	// has to be after the mod/env updates
	v.KeyModulator.DeferredUpdate()

	if o, ok := v.voicer.(*component.OPL2[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]); ok {
		fp, err := v.GetFinalPeriod()
		if err != nil {
			return err
		}
		o.Advance(v.GetFinalVolume(), fp)
	}

	v.KeyModulator.Advance()
	return nil
}

func (v *s3mVoice) RowEnd() error {
	v.vol0Opt.ObserveVolume(v.GetFinalVolume())
	return nil
}

func (v *s3mVoice) Clone(bool) voice.Voice {
	vv := s3mVoice{
		inst:          v.inst,
		opl2Chip:      v.opl2Chip,
		opl2Channel:   v.opl2Channel,
		AmpModulator:  v.AmpModulator.Clone(),
		FreqModulator: v.FreqModulator.Clone(),
		PanModulator:  v.PanModulator.Clone(),
		opl2:          v.opl2.Clone(),
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
