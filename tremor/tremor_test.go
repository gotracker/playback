package tremor

import "testing"

func TestTremorToggleAdvanceReset(t *testing.T) {
    var tr Tremor

    if !tr.IsActive() {
        t.Fatalf("expected tremor to start active")
    }

    tr.ToggleAndReset()
    if tr.IsActive() {
        t.Fatalf("expected tremor to toggle off")
    }
    if tr.Advance() != 1 {
        t.Fatalf("expected tick to increment after toggle")
    }

    tr.Reset()
    if tr.tick != 0 || !tr.IsActive() {
        t.Fatalf("expected reset to clear tick and enable: %+v", tr)
    }
}
