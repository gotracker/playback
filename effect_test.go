package playback

import (
	"fmt"
	"testing"
)

type stubEffecter struct{}
type stubEffect struct{ label string }

func (s stubEffect) TraceData() string { return s.label }
func (s stubEffect) String() string    { return s.label }

type nameEffect struct{ names []string }

func (n nameEffect) TraceData() string { return "" }
func (n nameEffect) Names() []string   { return n.names }

func TestGetEffectNamesPrefersNames(t *testing.T) {
	e := nameEffect{names: []string{"x", "y"}}
	names := GetEffectNames(e)
	if fmt.Sprint(names) != "[x y]" {
		t.Fatalf("unexpected names: %v", names)
	}

	s := stubEffect{label: "label"}
	names = GetEffectNames(s)
	if fmt.Sprint(names) != "[label]" {
		t.Fatalf("expected stringer name fallback, got %v", names)
	}
}
