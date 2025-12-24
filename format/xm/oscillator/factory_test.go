package oscillator

import "testing"

func TestOscillatorFactoryKnownAndUnknown(t *testing.T) {
    if o, err := Factory("vibrato"); err != nil || o == nil {
        t.Fatalf("expected vibrato oscillator, got %v, err=%v", o, err)
    }
    if o, err := Factory(""); err != nil || o != nil {
        t.Fatalf("expected empty name to return nil without error, got %v, err=%v", o, err)
    }
    if _, err := Factory("nope"); err == nil {
        t.Fatalf("expected error for unknown oscillator")
    }
}
