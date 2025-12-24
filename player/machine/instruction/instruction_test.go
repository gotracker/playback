package instruction

import "testing"

func TestValueStoresTypedValue(t *testing.T) {
	v := Value[string]{Value: "hello"}
	if v.Value != "hello" {
		t.Fatalf("Value stored %q, want %q", v.Value, "hello")
	}
}
