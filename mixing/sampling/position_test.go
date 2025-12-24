package sampling

import "testing"

func TestPosAddCarriesFraction(t *testing.T) {
    p := Pos{Pos: 1, Frac: 0.75}
    p.Add(0.5)
    if p.Pos != 2 {
        t.Fatalf("expected pos carry to 2, got %d", p.Pos)
    }
    if p.Frac < 0.24 || p.Frac > 0.26 {
        t.Fatalf("expected frac around 0.25, got %f", p.Frac)
    }
}
