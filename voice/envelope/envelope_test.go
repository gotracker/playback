package envelope

import "testing"

type dummyLoop struct{ on bool }

func (d dummyLoop) Enabled() bool                           { return d.on }
func (d dummyLoop) Length() int                             { return 3 }
func (d dummyLoop) CalcPos(pos int, length int) (int, bool) { return pos + 1, true }

func TestEnvelopeDefaults(t *testing.T) {
	var e Envelope[int]
	if e.Enabled {
		t.Fatalf("expected default Enabled=false")
	}
	if e.Loop != nil || e.Sustain != nil {
		t.Fatalf("expected nil loops by default")
	}
	if e.Length != 0 || len(e.Values) != 0 {
		t.Fatalf("expected zero length/values")
	}
}

func TestEnvelopeStoresValues(t *testing.T) {
	e := Envelope[int]{
		Enabled: true,
		Loop:    dummyLoop{on: true},
		Sustain: dummyLoop{on: true},
		Length:  2,
		Values: []Point[int]{
			{Pos: 0, Length: 1, Y: 10},
			{Pos: 1, Length: 1, Y: 20},
		},
	}
	if !e.Enabled {
		t.Fatalf("expected Enabled=true")
	}
	if !e.Loop.Enabled() || !e.Sustain.Enabled() {
		t.Fatalf("expected loops enabled")
	}
	if e.Length != 2 || len(e.Values) != 2 {
		t.Fatalf("unexpected length/values")
	}
	if e.Values[0].Y != 10 || e.Values[1].Y != 20 {
		t.Fatalf("unexpected point values: %+v", e.Values)
	}
}
