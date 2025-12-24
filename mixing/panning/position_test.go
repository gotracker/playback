package panning

import (
	"math"
	"testing"
)

func TestMakeStereoPositionRoundTrip(t *testing.T) {
	pos := MakeStereoPosition(0, -1, 1)
	if pos.Distance != 1 {
		t.Fatalf("expected distance 1, got %v", pos.Distance)
	}
	if math.Abs(float64(pos.Angle-math.Pi/4)) > 1e-6 {
		t.Fatalf("unexpected angle: %v", pos.Angle)
	}

	value := FromStereoPosition(pos, -1, 1)
	if math.Abs(float64(value)) > 1e-6 {
		t.Fatalf("expected value round trip to 0, got %v", value)
	}
}

func TestMakeStereoPositionPanicsOnEqualBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when bounds are equal")
		}
	}()
	_ = MakeStereoPosition(0, 1, 1)
}

func TestFromStereoPositionPanicsOnEqualBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when bounds are equal")
		}
	}()
	_ = FromStereoPosition(Position{}, 2, 2)
}
