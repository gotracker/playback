package period

import "testing"

func TestAmigaPortaDownUp(t *testing.T) {
	minP, maxP := Amiga(100), Amiga(1000)

	if got := Amiga(500).PortaDown(Delta(10), minP, maxP, false); got != 510 {
		t.Fatalf("PortaDown expected 510, got %d", got)
	}

	if got := Amiga(500).PortaUp(Delta(10), minP, maxP, false); got != 490 {
		t.Fatalf("PortaUp expected 490, got %d", got)
	}
}

func TestAmigaIsInvalid(t *testing.T) {
	if !Amiga(0).IsInvalid() {
		t.Fatalf("IsInvalid should be true for 0")
	}

	if Amiga(500).IsInvalid() {
		t.Fatalf("IsInvalid should be false for non-zero period")
	}
}
