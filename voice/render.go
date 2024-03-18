package voice

import (
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/mixer"
)

func RenderAndTick[TPeriod Period](in Voice, pc period.PeriodConverter[TPeriod], centerAheadPan panning.PanMixer, details mixer.Details, out mixer.ApplyFilter) (*mixing.Data, error) {
	if in.IsDone() {
		return nil, nil
	}

	defer in.Tick()

	rs, ok := in.(RenderSampler[TPeriod])
	if !ok {
		return nil, nil
	}

	if !rs.IsActive() {
		return nil, nil
	}

	pos, err := rs.GetPos()
	if err != nil {
		return nil, err
	}

	p, err := rs.GetFinalPeriod()
	if err != nil {
		return nil, err
	}

	if err := in.SetPlaybackRate(details.SampleRate); err != nil {
		return nil, err
	}

	samplerAdd := float32(pc.GetSamplerAdd(p, rs.GetSampleRate(), details.SampleRate))

	o := mixer.Output{
		Input:  rs,
		Output: out,
	}

	sampler := sampling.NewSampler(&o, pos, samplerAdd)

	// ... so grab the new value now.
	pan := rs.GetFinalPan()

	// make a stand-alone data buffer for this channel for this tick
	sampleData := mixing.SampleMixIn{
		Sample:    sampler,
		StaticVol: volume.Volume(1.0),
		PanMatrix: centerAheadPan,
		MixPos:    0,
		MixLen:    details.Samples,
	}

	mixBuffer := details.Mix.NewMixBuffer(details.Samples)
	mixBuffer.MixInSample(sampleData)
	data := &mixing.Data{
		Data:       mixBuffer,
		PanMatrix:  details.Panmixer.GetMixingMatrix(pan, details.StereoSeparation),
		Volume:     volume.Volume(1.0),
		Pos:        0,
		SamplesLen: details.Samples,
	}

	// reflect the sampling position back to the voice
	if err := rs.SetPos(sampler.GetPosition()); err != nil {
		return nil, err
	}

	return data, nil
}
