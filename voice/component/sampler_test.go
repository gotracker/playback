package component

import (
	"math"
	"testing"

	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/voice/pcm"
	"github.com/gotracker/playback/voice/types"
)

type testVolume float32

func (testVolume) IsInvalid() bool                       { return false }
func (testVolume) IsUseInstrumentVol() bool              { return false }
func (v testVolume) ToVolume() volume.Volume             { return volume.Volume(v) }
func (testVolume) AddDelta(types.VolumeDelta) testVolume { return testVolume(0) }
func (testVolume) GetMax() testVolume                    { return testVolume(1) }

func almostEqualVol(a, b volume.Volume, tol float64) bool {
	return math.Abs(float64(a-b)) <= tol
}

func TestSamplerLerpsBetweenSamples(t *testing.T) {
	data := []volume.Matrix{
		{StaticMatrix: volume.StaticMatrix{1}, Channels: 1},
		{StaticMatrix: volume.StaticMatrix{0}, Channels: 1},
	}
	samp := pcm.NewSampleNative(data, len(data), 1)

	var s Sampler[types.Period, testVolume, testVolume]
	s.Setup(SamplerSettings[types.Period, testVolume, testVolume]{
		Sample:        samp,
		DefaultVolume: testVolume(1),
		MixVolume:     testVolume(1),
	})

	got := s.GetSample(sampling.Pos{Pos: 0, Frac: 0.5})
	want := volume.Matrix{StaticMatrix: volume.StaticMatrix{0.5}, Channels: 1}
	if got.Channels != want.Channels || !almostEqualVol(got.StaticMatrix[0], want.StaticMatrix[0], 1e-6) {
		t.Fatalf("unexpected lerp result: got %+v want %+v", got, want)
	}
}

func TestSamplerFadeoutPastEndWithoutLoop(t *testing.T) {
	data := []volume.Matrix{
		{StaticMatrix: volume.StaticMatrix{0.5}, Channels: 1},
		{StaticMatrix: volume.StaticMatrix{0.25}, Channels: 1},
	}
	samp := pcm.NewSampleNative(data, len(data), 1)

	var s Sampler[types.Period, testVolume, testVolume]
	s.Setup(SamplerSettings[types.Period, testVolume, testVolume]{
		Sample:        samp,
		DefaultVolume: testVolume(1),
		MixVolume:     testVolume(1),
	})

	got := s.GetSample(sampling.Pos{Pos: 3})
	want := volume.Matrix{StaticMatrix: volume.StaticMatrix{0.125}, Channels: 1}
	if got.Channels != want.Channels || !almostEqualVol(got.StaticMatrix[0], want.StaticMatrix[0], 1e-6) {
		t.Fatalf("unexpected fadeout result: got %+v want %+v", got, want)
	}
}
