package util

import "testing"

func TestLerpClampsAndInterpolates(t *testing.T) {
    if got := Lerp(-0.1, 0, 10); got != 0 {
        t.Fatalf("expected clamp to a when t<0, got %d", got)
    }
    if got := Lerp(1.5, 0, 10); got != 10 {
        t.Fatalf("expected clamp to b when t>1, got %d", got)
    }
    if got := Lerp(0.5, 0, 10); got != 5 {
        t.Fatalf("expected midpoint interpolation to be 5, got %d", got)
    }
}
