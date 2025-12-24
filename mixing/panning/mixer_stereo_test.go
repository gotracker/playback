package panning

import (
	"testing"

	"github.com/gotracker/playback/mixing/volume"
)

func TestMixerStereoApplyUsesSeparation(t *testing.T) {
	mtx := volume.Matrix{StaticMatrix: volume.StaticMatrix{2, 3}, Channels: 2}
	m := MixerStereo{
		Matrix:               volume.Matrix{StaticMatrix: volume.StaticMatrix{1, 2}, Channels: 2},
		StereoSeparationFunc: func(in volume.Matrix) volume.Matrix { return in.Apply(volume.Volume(0.5)) },
	}

	got := m.ApplyToMatrix(mtx)
	if got.Channels != 2 {
		t.Fatalf("expected stereo output, got %d", got.Channels)
	}
	if got.StaticMatrix[0] != volume.Volume(1) || got.StaticMatrix[1] != volume.Volume(3) {
		t.Fatalf("unexpected matrix result: %+v", got.StaticMatrix)
	}
}
