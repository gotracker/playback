package voice

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	playerRender "github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/voice/component"
)

type voicerPos interface {
	GetPos() sampling.Pos
	SetPos(pos sampling.Pos)
}

type voicerSampler interface {
	GetSample(pos sampling.Pos) volume.Matrix
}

func (v *s3mVoice) GetPos() (sampling.Pos, error) {
	if vp, ok := v.voicer.(voicerPos); ok {
		return vp.GetPos(), nil
	}
	return sampling.Pos{}, nil
}

func (v *s3mVoice) SetPos(pos sampling.Pos) error {
	if vp, ok := v.voicer.(voicerPos); ok {
		vp.SetPos(pos)
	}
	return nil
}

func (v *s3mVoice) GetSample(pos sampling.Pos) volume.Matrix {
	var samp volume.Matrix
	if sampler, ok := v.voicer.(voicerSampler); ok {
		samp = sampler.GetSample(pos)
		if samp.Channels == 0 {
			samp.Channels = v.voicer.GetNumChannels()
		}
	}

	vol := v.GetFinalVolume()
	wet := samp.Apply(vol)
	if v.config.VoiceFilter != nil {
		wet = v.config.VoiceFilter.Filter(wet)
	}
	return wet
}

func (v *s3mVoice) GetSampler(samplerRate float32, renderChannel *playerRender.Channel[s3mVolume.Volume, s3mVolume.FineVolume, s3mPanning.Panning]) (sampling.Sampler, error) {
	o := component.OutputFilter{
		Input:  v,
		Output: renderChannel,
	}

	pos, err := v.GetPos()
	if err != nil {
		return nil, err
	}

	p := v.GetFinalPeriod()

	samplerAdd := float32(v.ms.PeriodConverter.GetSamplerAdd(p, float64(v.config.SampleRate)*float64(samplerRate)))
	s := sampling.NewSampler(&o, pos, samplerAdd)

	return s, nil
}
