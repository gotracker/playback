package voice

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/volume"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	playerRender "github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/state/render"
)

func (v *xmVoice[TPeriod]) Render(centerAheadPan volume.Matrix, details render.Details, renderChannel *playerRender.Channel[xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) (*mixing.Data, error) {
	if v.IsDone() {
		return nil, nil
	}

	if !v.IsActive() {
		return nil, nil
	}

	sampler, err := v.GetSampler(details.SamplerSpeed, renderChannel)
	if err != nil || sampler == nil {
		return nil, err
	}

	// ... so grab the new value now.
	pan := v.GetFinalPan()

	// make a stand-alone data buffer for this channel for this tick
	sampleData := mixing.SampleMixIn{
		Sample:    sampler,
		StaticVol: volume.Volume(1.0),
		VolMatrix: centerAheadPan,
		MixPos:    0,
		MixLen:    details.Samples,
	}

	mixBuffer := details.Mix.NewMixBuffer(details.Samples)
	mixBuffer.MixInSample(sampleData)
	data := &mixing.Data{
		Data:       mixBuffer,
		Pan:        pan.ToPosition(),
		Volume:     volume.Volume(1.0),
		Pos:        0,
		SamplesLen: details.Samples,
	}

	if err := v.SetPos(sampler.GetPosition()); err != nil {
		return nil, err
	}

	return data, nil
}
