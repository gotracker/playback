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

func (v *itVoice[TPeriod]) GetPos() (sampling.Pos, error) {
	if vp, ok := v.voicer.(voicerPos); ok {
		return vp.GetPos(), nil
	}
	return sampling.Pos{}, nil
}

func (v *itVoice[TPeriod]) SetPos(pos sampling.Pos) error {
	if vp, ok := v.voicer.(voicerPos); ok {
		vp.SetPos(pos)
	}
	return nil
}

func (v *itVoice[TPeriod]) GetSample(pos sampling.Pos) volume.Matrix {
	var samp volume.Matrix
	if sampler, ok := v.voicer.(voicerSampler); ok {
		samp = sampler.GetSample(pos)
		if samp.Channels == 0 {
			samp.Channels = v.voicer.GetNumChannels()
		}
	}

	vol := v.GetFinalVolume()
	wet := samp.Apply(vol)
	if v.voiceFilter != nil {
		wet = v.voiceFilter.Filter(wet)
	}
	return wet
}

func (v itVoice[TPeriod]) GetSampleRate() frequency.Frequency {
	return v.inst.SampleRate
}
