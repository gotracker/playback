package voice

import (
	"time"

	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/render"

	"github.com/gotracker/playback/instrument"
)

// OPL2 is an OPL2 voice interface
type OPL2[TPeriod period.Period] interface {
	voice.Voice
	voice.FreqModulator[TPeriod]
	voice.AmpModulator
	voice.VolumeEnveloper
	voice.PitchEnveloper[TPeriod]
}

// OPL2Registers is a set of OPL operator configurations
type OPL2Registers component.OPL2Registers

// OPLConfiguration is the information needed to configure an OPL2 voice
type OPLConfiguration[TPeriod period.Period] struct {
	Chip          render.OPL2Chip
	Channel       int
	SampleRate    period.Frequency
	InitialVolume volume.Volume
	InitialPeriod TPeriod
	AutoVibrato   voice.AutoVibrato
	Data          instrument.Data
}

// == the actual opl2 voice ==

type opl2Voice[TPeriod period.Period] struct {
	sampleRate    period.Frequency
	initialVolume volume.Volume

	active    bool
	keyOn     bool
	prevKeyOn bool

	fadeoutMode fadeout.Mode

	o        component.OPL2[TPeriod]
	amp      component.AmpModulator
	freq     component.FreqModulator[TPeriod]
	volEnv   component.VolumeEnvelope
	pitchEnv component.PitchEnvelope[TPeriod]
}

// NewOPL2 creates a new OPL2 voice
func NewOPL2[TPeriod period.Period](config OPLConfiguration[TPeriod]) voice.Voice {
	v := opl2Voice[TPeriod]{
		sampleRate:    config.SampleRate,
		initialVolume: config.InitialVolume,
		fadeoutMode:   fadeout.ModeDisabled,
		active:        true,
	}

	var regs component.OPL2Registers

	switch d := config.Data.(type) {
	case *instrument.OPL2:
		v.amp.Setup(1)
		v.amp.ResetFadeoutValue(0)
		v.volEnv.SetEnabled(false)
		v.volEnv.Reset(nil)
		v.pitchEnv.SetEnabled(false)
		v.pitchEnv.Reset(nil)
		regs.Mod.Reg20 = d.Modulator.GetReg20()
		regs.Mod.Reg40 = d.Modulator.GetReg40()
		regs.Mod.Reg60 = d.Modulator.GetReg60()
		regs.Mod.Reg80 = d.Modulator.GetReg80()
		regs.Mod.RegE0 = d.Modulator.GetRegE0()
		regs.Car.Reg20 = d.Carrier.GetReg20()
		regs.Car.Reg40 = d.Carrier.GetReg40()
		regs.Car.Reg60 = d.Carrier.GetReg60()
		regs.Car.Reg80 = d.Carrier.GetReg80()
		regs.Car.RegE0 = d.Carrier.GetRegE0()
		regs.RegC0 = d.GetRegC0()
	default:
		_ = d
	}

	v.o.Setup(config.Chip, config.Channel, regs, config.SampleRate)
	v.amp.SetVolume(config.InitialVolume)
	v.freq.SetPeriod(config.InitialPeriod)
	v.freq.SetAutoVibratoEnabled(config.AutoVibrato.Enabled)
	if config.AutoVibrato.Enabled {
		v.freq.ConfigureAutoVibrato(config.AutoVibrato)
		v.freq.ResetAutoVibrato(config.AutoVibrato.Sweep)
	}

	return &v
}

// == Controller ==

func (v *opl2Voice[TPeriod]) Attack() {
	v.keyOn = true
	v.amp.Attack()
	v.freq.ResetAutoVibrato()
	v.SetVolumeEnvelopePosition(0)
	v.SetPitchEnvelopePosition(0)

}

func (v *opl2Voice[TPeriod]) Release() {
	v.keyOn = false
	v.amp.Release()
	v.o.Release()
}

func (v *opl2Voice[TPeriod]) Fadeout() {
	switch v.fadeoutMode {
	case fadeout.ModeAlwaysActive:
		v.amp.Fadeout()
	case fadeout.ModeOnlyIfVolEnvActive:
		if v.IsVolumeEnvelopeEnabled() {
			v.amp.Fadeout()
		}
	}
}

func (v *opl2Voice[TPeriod]) IsKeyOn() bool {
	return v.keyOn
}

func (v *opl2Voice[TPeriod]) IsFadeout() bool {
	return v.amp.IsFadeoutEnabled()
}

func (v *opl2Voice[TPeriod]) IsDone() bool {
	if !v.amp.IsFadeoutEnabled() {
		return false
	}
	return v.amp.GetFadeoutVolume() <= 0
}

// == FreqModulator ==

func (v *opl2Voice[TPeriod]) SetPeriod(period TPeriod) {
	v.freq.SetPeriod(period)
}

func (v *opl2Voice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *opl2Voice[TPeriod]) SetPeriodDelta(delta period.Delta) {
	v.freq.SetDelta(delta)
}

func (v *opl2Voice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetDelta()
}

func (v *opl2Voice[TPeriod]) GetFinalPeriod() TPeriod {
	p := v.freq.GetFinalPeriod()
	if v.IsPitchEnvelopeEnabled() {
		d := v.GetCurrentPitchEnvelope()
		p = period.AddDelta(p, d)
	}
	return p
}

// == AmpModulator ==

func (v *opl2Voice[TPeriod]) SetVolume(vol volume.Volume) {
	if vol == volume.VolumeUseInstVol {
		vol = v.initialVolume
	}
	v.amp.SetVolume(vol)
}

func (v *opl2Voice[TPeriod]) GetVolume() volume.Volume {
	return v.amp.GetVolume()
}

func (v *opl2Voice[TPeriod]) GetFinalVolume() volume.Volume {
	vol := v.amp.GetFinalVolume()
	if v.IsVolumeEnvelopeEnabled() {
		vol *= v.GetCurrentVolumeEnvelope()
	}
	return vol
}

// == VolumeEnveloper ==

func (v *opl2Voice[TPeriod]) EnableVolumeEnvelope(enabled bool) {
	v.volEnv.SetEnabled(enabled)
}

func (v *opl2Voice[TPeriod]) IsVolumeEnvelopeEnabled() bool {
	return v.volEnv.IsEnabled()
}

func (v *opl2Voice[TPeriod]) GetCurrentVolumeEnvelope() volume.Volume {
	if v.volEnv.IsEnabled() {
		return v.volEnv.GetCurrentValue()
	}
	return 1
}

func (v *opl2Voice[TPeriod]) SetVolumeEnvelopePosition(pos int) {
	if doneCB := v.volEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

// == PitchEnveloper ==

func (v *opl2Voice[TPeriod]) EnablePitchEnvelope(enabled bool) {
	v.pitchEnv.SetEnabled(enabled)
}

func (v *opl2Voice[TPeriod]) IsPitchEnvelopeEnabled() bool {
	return v.pitchEnv.IsEnabled()
}

func (v *opl2Voice[TPeriod]) GetCurrentPitchEnvelope() period.Delta {
	if v.pitchEnv.IsEnabled() {
		return v.pitchEnv.GetCurrentValue()
	}
	var empty period.Delta
	return empty
}

func (v *opl2Voice[TPeriod]) SetPitchEnvelopePosition(pos int) {
	if doneCB := v.pitchEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

// == required function interfaces ==

func (v *opl2Voice[TPeriod]) Advance(tickDuration time.Duration) {
	defer func() {
		v.prevKeyOn = v.keyOn
	}()
	v.amp.Advance()
	v.freq.Advance()
	if v.IsVolumeEnvelopeEnabled() {
		if doneCB := v.volEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsPitchEnvelopeEnabled() {
		if doneCB := v.pitchEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}

	// has to be after the mod/env updates
	if v.keyOn != v.prevKeyOn {
		if v.keyOn {
			v.o.Attack()
		} else {
			v.o.Release()
		}
	}

	v.o.Advance(v.GetFinalVolume(), v.GetFinalPeriod())
}

func (v *opl2Voice[TPeriod]) GetSample(pos sampling.Pos) volume.Matrix {
	return volume.Matrix{}
}

func (v *opl2Voice[TPeriod]) GetSampler(samplerRate float32) sampling.Sampler {
	return nil
}

func (v *opl2Voice[TPeriod]) Clone() voice.Voice {
	o := *v
	return &o
}

func (v *opl2Voice[TPeriod]) StartTransaction() voice.Transaction[TPeriod] {
	t := txn[TPeriod]{
		Voice: v,
	}
	return &t
}

func (v *opl2Voice[TPeriod]) SetActive(active bool) {
	v.active = active
}

func (v *opl2Voice[TPeriod]) IsActive() bool {
	return v.active
}
