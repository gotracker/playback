package voice

import (
	"time"

	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/component"
	"github.com/gotracker/playback/voice/fadeout"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/pan"
)

// PCM is a PCM voice interface
type PCM[TPeriod period.Period] interface {
	voice.Voice
	voice.Positioner
	voice.FreqModulator[TPeriod]
	voice.AmpModulator
	voice.PanModulator
	voice.VolumeEnveloper
	voice.PitchEnveloper[TPeriod]
	voice.PanEnveloper
	voice.FilterEnveloper
}

// PCMConfiguration is the information needed to configure an PCM2 voice
type PCMConfiguration[TPeriod period.Period] struct {
	SampleRate    period.Frequency
	InitialVolume volume.Volume
	InitialPeriod TPeriod
	AutoVibrato   voice.AutoVibrato
	Data          instrument.Data
	OutputFilter  voice.FilterApplier
	VoiceFilter   filter.Filter
	PluginFilter  filter.Filter
}

// == the actual pcm voice ==

type pcmVoice[TPeriod period.Period] struct {
	sampleRate    period.Frequency
	initialVolume volume.Volume
	outputFilter  voice.FilterApplier
	voiceFilter   filter.Filter
	pluginFilter  filter.Filter
	fadeoutMode   fadeout.Mode
	channels      int

	active    bool
	keyOn     bool
	prevKeyOn bool

	pitchAndFilterEnvShared bool
	filterEnvActive         bool // if pitchAndFilterEnvShared is true, this dictates which is active initially - true=filter, false=pitch

	sampler   component.Sampler
	amp       component.AmpModulator
	freq      component.FreqModulator[TPeriod]
	pan       component.PanModulator
	volEnv    component.VolumeEnvelope
	pitchEnv  component.PitchEnvelope[TPeriod]
	panEnv    component.PanEnvelope
	filterEnv component.FilterEnvelope
	vol0ticks int
	done      bool

	periodConverter period.PeriodConverter[TPeriod]
}

// NewPCM creates a new PCM voice
func NewPCM[TPeriod period.Period](periodConverter period.PeriodConverter[TPeriod], config PCMConfiguration[TPeriod]) voice.Voice {
	v := pcmVoice[TPeriod]{
		sampleRate:      config.SampleRate,
		initialVolume:   config.InitialVolume,
		outputFilter:    config.OutputFilter,
		voiceFilter:     config.VoiceFilter,
		pluginFilter:    config.PluginFilter,
		active:          true,
		periodConverter: periodConverter,
	}

	switch d := config.Data.(type) {
	case *instrument.PCM:
		v.pitchAndFilterEnvShared = true
		v.filterEnvActive = d.PitchFiltMode
		v.sampler.Setup(d.Sample, d.Loop, d.SustainLoop)
		v.fadeoutMode = d.FadeOut.Mode
		//v.sampler.SetPos(d.InitialPos)
		v.amp.Setup(d.MixingVolume)
		v.amp.ResetFadeoutValue(d.FadeOut.Amount)
		v.pan.SetPan(d.Panning)
		v.volEnv.SetEnabled(d.VolEnv.Enabled)
		v.volEnv.Init(&d.VolEnv)
		v.pitchEnv.SetEnabled(d.PitchFiltEnv.Enabled)
		v.pitchEnv.Init(&d.PitchFiltEnv)
		v.panEnv.SetEnabled(d.PanEnv.Enabled)
		v.panEnv.Init(&d.PanEnv)
		v.filterEnv.SetEnabled(d.PitchFiltEnv.Enabled)
		v.filterEnv.Init(&d.PitchFiltEnv)
		v.channels = d.Sample.Channels()
	}

	switch v.fadeoutMode {
	case fadeout.ModeAlwaysActive:
		v.amp.SetFadeoutEnabled(true)
	case fadeout.ModeOnlyIfVolEnvActive:
		v.amp.SetFadeoutEnabled(v.volEnv.IsEnabled())
	default:
	}

	v.amp.SetVolume(config.InitialVolume)
	v.freq.SetPeriod(config.InitialPeriod)
	v.freq.SetAutoVibratoEnabled(config.AutoVibrato.Enabled)
	if config.AutoVibrato.Enabled {
		v.freq.ConfigureAutoVibrato(config.AutoVibrato)
		v.freq.ResetAutoVibratoAndSweep(config.AutoVibrato.Sweep)
	}

	return &v
}

// == Controller ==

func (v *pcmVoice[TPeriod]) Attack() {
	v.keyOn = true
	v.vol0ticks = 0
	v.done = false
	v.amp.Attack()
	v.freq.ResetAutoVibrato()
	v.sampler.Attack()
	v.SetVolumeEnvelopePosition(0)
	v.SetPitchEnvelopePosition(0)
	v.SetPanEnvelopePosition(0)
	v.SetFilterEnvelopePosition(0)
}

func (v *pcmVoice[TPeriod]) Release() {
	v.keyOn = false
	v.amp.Release()
	v.sampler.Release()
}

func (v *pcmVoice[TPeriod]) Fadeout() {
	switch v.fadeoutMode {
	case fadeout.ModeAlwaysActive:
		v.amp.Fadeout()
	case fadeout.ModeOnlyIfVolEnvActive:
		if v.IsVolumeEnvelopeEnabled() {
			v.amp.Fadeout()
		}
	}

	v.sampler.Fadeout()
}

func (v *pcmVoice[TPeriod]) IsKeyOn() bool {
	return v.keyOn
}

func (v *pcmVoice[TPeriod]) IsFadeout() bool {
	return v.amp.IsFadeoutEnabled()
}

func (v *pcmVoice[TPeriod]) IsDone() bool {
	if v.done {
		return true
	}

	if v.amp.IsFadeoutEnabled() {
		return v.amp.GetFadeoutVolume() <= 0
	}

	return v.vol0ticks >= 3
}

// == SampleStream ==

func (v *pcmVoice[TPeriod]) GetSample(pos sampling.Pos) volume.Matrix {
	samp := v.sampler.GetSample(pos)
	if samp.Channels == 0 {
		v.done = true
		samp.Channels = v.channels
	}
	vol := v.GetFinalVolume()
	wet := samp.Apply(vol)
	if v.voiceFilter != nil {
		wet = v.voiceFilter.Filter(wet)
	}
	if v.pluginFilter != nil {
		wet = v.pluginFilter.Filter(wet)
	}
	return wet
}

// == Positioner ==

func (v *pcmVoice[TPeriod]) SetPos(pos sampling.Pos) {
	v.sampler.SetPos(pos)
}

func (v *pcmVoice[TPeriod]) GetPos() sampling.Pos {
	return v.sampler.GetPos()
}

// == FreqModulator ==

func (v *pcmVoice[TPeriod]) SetPeriod(period TPeriod) {
	v.freq.SetPeriod(period)
}

func (v *pcmVoice[TPeriod]) GetPeriod() TPeriod {
	return v.freq.GetPeriod()
}

func (v *pcmVoice[TPeriod]) SetPeriodDelta(delta period.Delta) {
	v.freq.SetDelta(delta)
}

func (v *pcmVoice[TPeriod]) GetPeriodDelta() period.Delta {
	return v.freq.GetDelta()
}

func (v *pcmVoice[TPeriod]) GetFinalPeriod() TPeriod {
	p := v.freq.GetFinalPeriod()
	if v.IsPitchEnvelopeEnabled() {
		delta := v.GetCurrentPitchEnvelope()
		p = period.AddDelta(p, delta)
	}
	return p
}

// == AmpModulator ==

func (v *pcmVoice[TPeriod]) SetVolume(vol volume.Volume) {
	if vol == volume.VolumeUseInstVol {
		vol = v.initialVolume
	}
	v.amp.SetVolume(vol)
}

func (v *pcmVoice[TPeriod]) GetVolume() volume.Volume {
	return v.amp.GetVolume()
}

func (v *pcmVoice[TPeriod]) GetFinalVolume() volume.Volume {
	vol := v.amp.GetFinalVolume()
	if v.IsVolumeEnvelopeEnabled() {
		vol *= v.GetCurrentVolumeEnvelope()
	}
	return vol
}

// == PanModulator ==

func (v *pcmVoice[TPeriod]) SetPan(pan panning.Position) {
	v.pan.SetPan(pan)
}

func (v *pcmVoice[TPeriod]) GetPan() panning.Position {
	return v.pan.GetPan()
}

func (v *pcmVoice[TPeriod]) GetFinalPan() panning.Position {
	p := v.pan.GetFinalPan()
	if v.IsPanEnvelopeEnabled() {
		p = pan.CalculateCombinedPanning(p, v.panEnv.GetCurrentValue())
	}
	return p
}

// == VolumeEnveloper ==

func (v *pcmVoice[TPeriod]) EnableVolumeEnvelope(enabled bool) {
	v.volEnv.SetEnabled(enabled)
}

func (v *pcmVoice[TPeriod]) IsVolumeEnvelopeEnabled() bool {
	return v.volEnv.IsEnabled()
}

func (v *pcmVoice[TPeriod]) GetCurrentVolumeEnvelope() volume.Volume {
	if v.volEnv.IsEnabled() {
		return v.volEnv.GetCurrentValue()
	}
	return 1
}

func (v *pcmVoice[TPeriod]) SetVolumeEnvelopePosition(pos int) {
	if doneCB := v.volEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

// == PitchEnveloper ==

func (v *pcmVoice[TPeriod]) EnablePitchEnvelope(enabled bool) {
	v.pitchEnv.SetEnabled(enabled)
}

func (v *pcmVoice[TPeriod]) IsPitchEnvelopeEnabled() bool {
	if v.pitchAndFilterEnvShared && v.filterEnvActive {
		return false
	}
	return v.pitchEnv.IsEnabled()
}

func (v *pcmVoice[TPeriod]) GetCurrentPitchEnvelope() period.Delta {
	if v.pitchEnv.IsEnabled() {
		return v.pitchEnv.GetCurrentValue()
	}
	var empty period.Delta
	return empty
}

func (v *pcmVoice[TPeriod]) SetPitchEnvelopePosition(pos int) {
	if !v.pitchAndFilterEnvShared || !v.filterEnvActive {
		if doneCB := v.pitchEnv.SetEnvelopePosition(pos); doneCB != nil {
			doneCB(v)
		}
	}
}

// == FilterEnveloper ==

func (v *pcmVoice[TPeriod]) EnableFilterEnvelope(enabled bool) {
	if !v.pitchAndFilterEnvShared {
		v.filterEnv.SetEnabled(enabled)
		return
	}

	// shared filter/pitch envelope
	if !v.filterEnvActive {
		return
	}

	v.filterEnv.SetEnabled(enabled)
}

func (v *pcmVoice[TPeriod]) IsFilterEnvelopeEnabled() bool {
	if v.pitchAndFilterEnvShared && !v.filterEnvActive {
		return false
	}
	return v.filterEnv.IsEnabled()
}

func (v *pcmVoice[TPeriod]) GetCurrentFilterEnvelope() uint8 {
	return v.filterEnv.GetCurrentValue()
}

func (v *pcmVoice[TPeriod]) SetFilterEnvelopePosition(pos int) {
	if !v.pitchAndFilterEnvShared || v.filterEnvActive {
		if doneCB := v.filterEnv.SetEnvelopePosition(pos); doneCB != nil {
			doneCB(v)
		}
	}
}

// == PanEnveloper ==

func (v *pcmVoice[TPeriod]) EnablePanEnvelope(enabled bool) {
	v.panEnv.SetEnabled(enabled)
}

func (v *pcmVoice[TPeriod]) IsPanEnvelopeEnabled() bool {
	return v.panEnv.IsEnabled()
}

func (v *pcmVoice[TPeriod]) GetCurrentPanEnvelope() panning.Position {
	return v.panEnv.GetCurrentValue()
}

func (v *pcmVoice[TPeriod]) SetPanEnvelopePosition(pos int) {
	if doneCB := v.panEnv.SetEnvelopePosition(pos); doneCB != nil {
		doneCB(v)
	}
}

// == required function interfaces ==

func (v *pcmVoice[TPeriod]) Advance(tickDuration time.Duration) {
	v.amp.Advance()
	v.freq.Advance()
	v.pan.Advance()
	if v.IsVolumeEnvelopeEnabled() {
		if doneCB := v.volEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsPanEnvelopeEnabled() {
		if doneCB := v.panEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsPitchEnvelopeEnabled() {
		if doneCB := v.pitchEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}
	if v.IsFilterEnvelopeEnabled() {
		if doneCB := v.filterEnv.Advance(v.keyOn, v.prevKeyOn); doneCB != nil {
			doneCB(v)
		}
	}

	if v.voiceFilter != nil && v.IsFilterEnvelopeEnabled() {
		fval := v.GetCurrentFilterEnvelope()
		v.voiceFilter.UpdateEnv(fval)
	}

	if vol := v.GetFinalVolume(); vol <= 0 {
		v.vol0ticks++
	} else {
		v.vol0ticks = 0
	}
	v.prevKeyOn = v.keyOn
}

func (v *pcmVoice[TPeriod]) GetSampler(samplerRate float32) sampling.Sampler {
	p := v.GetFinalPeriod()
	samplerAdd := float32(v.periodConverter.GetSamplerAdd(p, float64(samplerRate)))
	o := component.OutputFilter{
		Input:  v,
		Output: v.outputFilter,
	}
	return sampling.NewSampler(&o, v.GetPos(), samplerAdd)
}

func (v *pcmVoice[TPeriod]) Clone() voice.Voice {
	p := pcmVoice[TPeriod]{
		sampleRate:              v.sampleRate,
		initialVolume:           v.initialVolume,
		outputFilter:            v.outputFilter,
		fadeoutMode:             v.fadeoutMode,
		channels:                v.channels,
		active:                  true,
		keyOn:                   false,
		prevKeyOn:               false,
		pitchAndFilterEnvShared: v.pitchAndFilterEnvShared,
		filterEnvActive:         v.filterEnvActive,
		sampler:                 v.sampler.Clone(),
		amp:                     v.amp,
		freq:                    v.freq,
		pan:                     v.pan,
		volEnv:                  v.volEnv.Clone(),
		pitchEnv:                v.pitchEnv.Clone(),
		panEnv:                  v.panEnv.Clone(),
		filterEnv:               v.filterEnv.Clone(),
		vol0ticks:               0,
		done:                    false,
		periodConverter:         v.periodConverter,
	}

	if v.voiceFilter != nil {
		p.voiceFilter = v.voiceFilter.Clone()
	}
	if v.pluginFilter != nil {
		p.pluginFilter = v.pluginFilter.Clone()
	}
	return &p
}

func (v *pcmVoice[TPeriod]) StartTransaction() voice.Transaction[TPeriod] {
	t := txn[TPeriod]{
		Voice: v,
	}
	return &t
}

func (v *pcmVoice[TPeriod]) SetActive(active bool) {
	v.active = active
}

func (v *pcmVoice[TPeriod]) IsActive() bool {
	return v.active
}
