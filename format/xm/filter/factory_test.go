package filter

import "testing"

func TestFilterFactoryKnownAndUnknown(t *testing.T) {
    if f, err := Factory("amigalpf", 8363, nil); err != nil || f == nil {
        t.Fatalf("expected amigalpf filter, got %v, err=%v", f, err)
    }
    if f, err := Factory("", 8363, nil); err != nil || f != nil {
        t.Fatalf("expected empty filter to return nil without error, got %v, err=%v", f, err)
    }
    if _, err := Factory("nope", 8363, nil); err == nil {
        t.Fatalf("expected error for unknown filter")
    }
}
