package index

import "testing"

func TestChannelValidity(t *testing.T) {
    if !Channel(3).IsValid() {
        t.Fatalf("expected channel 3 to be valid")
    }
    if Channel(InvalidChannel).IsValid() {
        t.Fatalf("expected invalid channel to be invalid")
    }

    if !OPLChannel(1).IsValid() {
        t.Fatalf("expected OPL channel 1 to be valid")
    }
    if OPLChannel(InvalidOPLChannel).IsValid() {
        t.Fatalf("expected invalid OPL channel to be invalid")
    }
}

func TestRowIncrementOverflow(t *testing.T) {
    r := Row(0)
    if overflow := r.Increment(4); overflow {
        t.Fatalf("did not expect overflow on first increment")
    }
    if r != 1 {
        t.Fatalf("expected row to be 1, got %d", r)
    }

    r = 2
    if overflow := r.Increment(3); !overflow {
        t.Fatalf("expected overflow when incrementing at max")
    }
    if r != 3 {
        t.Fatalf("expected row to wrap to 3 (one past max), got %d", r)
    }
}
