package voice

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/frequency"
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
	var dry volume.Matrix
	if sampler, ok := v.voicer.(voicerSampler); ok {
		dry = sampler.GetSample(pos)
		if dry.Channels == 0 {
			dry.Channels = v.voicer.GetNumChannels()
		}
	}

	vol := v.GetFinalVolume()
	wet := dry.Apply(vol)
	if v.config.VoiceFilter != nil {
		wet = v.config.VoiceFilter.Filter(wet)
	}
	return wet
}

func (v s3mVoice) GetSampleRate() frequency.Frequency {
	return v.config.SampleRate
}
