package filter

import (
	"testing"

	"github.com/gotracker/playback/mixing/volume"
)

type fakeApplier struct {
	applied volume.Matrix
	setEnv  uint8
}

func (f *fakeApplier) ApplyFilter(dry volume.Matrix) volume.Matrix {
	f.applied = dry
	return dry.Apply(volume.Volume(0.5))
}

func (f *fakeApplier) SetFilterEnvelopeValue(envVal uint8) {
	f.setEnv = envVal
}

func TestApplierInterface(t *testing.T) {
	var _ Applier = (*fakeApplier)(nil)
	fa := &fakeApplier{}
	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2}
	wet := fa.ApplyFilter(dry)
	if fa.applied.Channels != 2 || fa.applied.StaticMatrix[0] != 1 {
		t.Fatalf("expected applied matrix recorded")
	}
	if wet.StaticMatrix[0] != 0.5 || wet.StaticMatrix[1] != 0.5 {
		t.Fatalf("expected halved volumes, got %#v", wet.StaticMatrix)
	}
	fa.SetFilterEnvelopeValue(7)
	if fa.setEnv != 7 {
		t.Fatalf("expected env set to 7, got %d", fa.setEnv)
	}
}
