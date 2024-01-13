package voice

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
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
	if v.config.VoiceFilter != nil {
		wet = v.config.VoiceFilter.Filter(wet)
	}
	if v.config.PluginFilter != nil {
		wet = v.config.PluginFilter.Filter(wet)
	}
	return wet
}

func (v *itVoice[TPeriod]) GetSampler(samplerRate float32, renderChannel *playerRender.Channel[itVolume.FineVolume, itVolume.FineVolume, itPanning.Panning]) (sampling.Sampler, error) {
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
