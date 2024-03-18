package voice

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
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
	if v.voiceFilter != nil {
		wet = v.voiceFilter.Filter(wet)
	}
	return wet
}

func (v s3mVoice) GetSampleRate() frequency.Frequency {
	return v.inst.SampleRate
}
