package mixer

import (
	"testing"

	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
)

type stubStream struct {
	gotPos sampling.Pos
	sample volume.Matrix
}

func (s *stubStream) GetSample(pos sampling.Pos) volume.Matrix {
	s.gotPos = pos
	return s.sample
}

type stubFilter struct {
	gotDry volume.Matrix
	out    volume.Matrix
	called int
}

func (f *stubFilter) ApplyFilter(dry volume.Matrix) volume.Matrix {
	f.called++
	f.gotDry = dry
	return f.out
}

func TestOutputGetSampleAppliesFilter(t *testing.T) {
	dry := volume.Matrix{StaticMatrix: volume.StaticMatrix{volume.Volume(0.25), volume.Volume(-0.5)}, Channels: 2}
	filtered := volume.Matrix{StaticMatrix: volume.StaticMatrix{volume.Volume(0.5), volume.Volume(-1)}, Channels: 2}

	stream := &stubStream{sample: dry}
	filter := &stubFilter{out: filtered}

	o := Output{Input: stream, Output: filter}

	pos := sampling.Pos{Pos: 3, Frac: 0.75}

	got := o.GetSample(pos)

	if stream.gotPos != pos {
		t.Fatalf("expected sample request at %+v, got %+v", pos, stream.gotPos)
	}

	if filter.called != 1 {
		t.Fatalf("expected filter to be called once, got %d", filter.called)
	}

	if filter.gotDry != dry {
		t.Fatalf("expected filter to receive dry sample %v, got %v", dry, filter.gotDry)
	}

	if got != filtered {
		t.Fatalf("expected filtered sample %v, got %v", filtered, got)
	}
}
