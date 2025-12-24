package channel

import "testing"

func TestMemoryEFGLinkModeSharesRegisters(t *testing.T) {
	mem := Memory{Shared: &SharedMemory{EFGLinkMode: true}}
	if got := mem.PortaDown(0x31); got != 0x31 {
		t.Fatalf("expected porta down to store 0x31, got 0x%02x", got)
	}
	if got := mem.PortaUp(0); got != 0x31 {
		t.Fatalf("expected porta up to reuse shared value, got 0x%02x", got)
	}
	if got := mem.PortaToNote(0); got != 0x31 {
		t.Fatalf("expected porta-to-note to reuse shared value, got 0x%02x", got)
	}
}

func TestMemoryStartOrderReset(t *testing.T) {
	mem := Memory{Shared: &SharedMemory{ResetMemoryAtStartOfOrder0: true}}
	mem.VolumeSlide(0x12)
	mem.PortaDown(0x34)
	mem.PortaUp(0x56)
	mem.Vibrato(0x78)

	mem.StartOrder0()

	if got, _ := mem.VolumeSlide(0); got != 0 {
		t.Fatalf("expected volume slide memory reset to 0, got 0x%02x", got)
	}
	if got := mem.PortaDown(0); got != 0 {
		t.Fatalf("expected porta down memory reset to 0, got 0x%02x", got)
	}
	if got := mem.PortaUp(0); got != 0 {
		t.Fatalf("expected porta up memory reset to 0, got 0x%02x", got)
	}
	if got, _ := mem.Vibrato(0); got != 0 {
		t.Fatalf("expected vibrato memory reset to 0, got 0x%02x", got)
	}
}
