package period

import "testing"

func TestLinearAddClampsAndPorta(t *testing.T) {
	p := Linear{Finetune: 5}
	if got := p.Add(Delta(-10)); got.Finetune != 1 {
		t.Fatalf("expected clamp to 1, got %d", got.Finetune)
	}

	a := Linear{Finetune: 10}
	b := Linear{Finetune: 12}
	if a.PortaUp(1).Finetune != 11 {
		t.Fatalf("porta up failed")
	}
	if b.PortaDown(2).Finetune != 10 {
		t.Fatalf("porta down failed")
	}

	c := Linear{Finetune: 10}
	target := Linear{Finetune: 12}
	if res := c.PortaTo(1, target); res.Finetune != 11 {
		t.Fatalf("porta to upward failed, got %d", res.Finetune)
	}
	if res := target.PortaTo(1, c); res.Finetune != 11 {
		t.Fatalf("porta to downward failed, got %d", res.Finetune)
	}
}

func TestLinearCompareAndLerp(t *testing.T) {
	a := Linear{Finetune: 5}
	b := Linear{Finetune: 6}
	if cmp := a.Compare(b); cmp != -1 {
		t.Fatalf("expected a < b")
	}
	if cmp := b.Compare(a); cmp != 1 {
		t.Fatalf("expected b > a")
	}
	if cmp := a.Compare(a); cmp != 0 {
		t.Fatalf("expected equal compare")
	}

	lerped := a.Lerp(0.5, b).(Linear)
	if lerped.Finetune != 5 {
		t.Fatalf("unexpected lerp result: %v", lerped.Finetune)
	}
}
