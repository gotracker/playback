package mixing

import (
	"math"
	"testing"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

func almostEqual(a, b float64, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestPanMixerMonoMatrix(t *testing.T) {
	mixer := GetPanMixer(1)
	if mixer == nil {
		t.Fatalf("expected mono pan mixer")
	}
	if mixer.NumChannels() != 1 {
		t.Fatalf("expected mono mixer to report 1 channel, got %d", mixer.NumChannels())
	}
	pan := panning.Position{Angle: 0, Distance: 1}
	m := mixer.GetMixingMatrix(pan, 0)
	matrix := m.Apply(volume.Volume(1))
	if matrix.Channels != 1 || !almostEqual(float64(matrix.StaticMatrix[0]), 1.0, 1e-6) {
		t.Fatalf("unexpected mono matrix: %+v", matrix)
	}
}

func TestPanMixerStereoCenteredCoefficients(t *testing.T) {
	mixer := GetPanMixer(2)
	pan := panning.Position{Angle: 0, Distance: 1}
	pm := mixer.GetMixingMatrix(pan, 1)

	stereo, ok := pm.(panning.MixerStereo)
	if !ok {
		t.Fatalf("expected MixerStereo")
	}

	matrix := stereo.Apply(volume.Volume(1))
	if matrix.Channels != 2 {
		t.Fatalf("expected 2 channels, got %d", matrix.Channels)
	}
	expected := volume.StereoCoeff
	if !almostEqual(float64(matrix.StaticMatrix[0]), float64(expected), 1e-6) ||
		!almostEqual(float64(matrix.StaticMatrix[1]), float64(expected), 1e-6) {
		t.Fatalf("unexpected coefficients: %+v", matrix.StaticMatrix)
	}
}

func TestPanMixerStereoApplyRespectsSeparation(t *testing.T) {
	mixer := GetPanMixer(2)
	pan := panning.Position{Angle: 0.2, Distance: 1}
	pm := mixer.GetMixingMatrix(pan, 0)

	stereo := pm.(panning.MixerStereo)

	matrix := stereo.Apply(volume.Volume(1))
	if matrix.Channels != 2 {
		t.Fatalf("expected 2 channels, got %d", matrix.Channels)
	}

	if !almostEqual(float64(matrix.StaticMatrix[0]), float64(matrix.StaticMatrix[1]), 1e-6) {
		t.Fatalf("expected separation=0 to collapse to mono, got L=%v R=%v", matrix.StaticMatrix[0], matrix.StaticMatrix[1])
	}
}

func TestPanMixerStereoSeparationZeroCollapsesToMono(t *testing.T) {
	mixer := GetPanMixer(2)
	pm := mixer.GetMixingMatrix(panning.Position{Angle: 0, Distance: 1}, 0)

	stereo := pm.(panning.MixerStereo)

	in := volume.Matrix{StaticMatrix: volume.StaticMatrix{1, -1}, Channels: 2}
	wet := stereo.StereoSeparationFunc(in)

	if wet.Channels != 2 {
		t.Fatalf("expected 2 channels after separation")
	}
	if !almostEqual(float64(wet.StaticMatrix[0]), 0, 1e-6) || !almostEqual(float64(wet.StaticMatrix[1]), 0, 1e-6) {
		t.Fatalf("expected collapse to mono zero: %+v", wet.StaticMatrix)
	}
}

func TestPanMixerStereoSeparationFullKeepsChannels(t *testing.T) {
	mixer := GetPanMixer(2)
	pm := mixer.GetMixingMatrix(panning.Position{Angle: 0, Distance: 1}, 1)

	stereo := pm.(panning.MixerStereo)

	in := volume.Matrix{StaticMatrix: volume.StaticMatrix{0.25, -0.25}, Channels: 2}
	wet := stereo.StereoSeparationFunc(in)

	if wet != in {
		t.Fatalf("expected stereo separation 1 to keep channels: got %+v want %+v", wet, in)
	}
}

func TestPanMixerQuadAngleZero(t *testing.T) {
	mixer := GetPanMixer(4)
	pan := panning.Position{Angle: 0, Distance: 1}
	pm := mixer.GetMixingMatrix(pan, 0)

	matrix := pm.Apply(volume.Volume(1))
	if matrix.Channels != 4 {
		t.Fatalf("expected 4-channel matrix, got %d", matrix.Channels)
	}
	expected := volume.StaticMatrix{0, 1, 0, 1}
	for i := 0; i < 4; i++ {
		if !almostEqual(float64(matrix.StaticMatrix[i]), float64(expected[i]), 1e-6) {
			t.Fatalf("channel %d mismatch: got %v want %v", i, matrix.StaticMatrix[i], expected[i])
		}
	}
}

func TestGetPanMixerUnknownChannels(t *testing.T) {
	if m := GetPanMixer(3); m != nil {
		t.Fatalf("expected nil mixer for unsupported channel count, got %#v", m)
	}
}
