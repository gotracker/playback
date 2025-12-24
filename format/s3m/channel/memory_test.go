package channel

import "testing"

func TestMemoryVibratoNibbles(t *testing.T) {
	var mem Memory
	vx, vy := mem.Vibrato(0x4f)
	if vx != 0x04 || vy != 0x0f {
		t.Fatalf("expected vibrato nibbles 0x04/0x0f, got 0x%02x/0x%02x", vx, vy)
	}
	// reuse last values
	vx, vy = mem.Vibrato(0)
	if vx != 0x04 || vy != 0x0f {
		t.Fatalf("expected vibrato memory reuse, got 0x%02x/0x%02x", vx, vy)
	}
}

func TestMemoryStartOrderReset(t *testing.T) {
	mem := Memory{Shared: &SharedMemory{ResetMemoryAtStartOfOrder0: true}}
	mem.Porta(0x22)
	mem.Vibrato(0x34)
	mem.Tremolo(0x56)
	mem.SampleOffset(0x78)
	mem.TempoDecrease(0x9a)
	mem.TempoIncrease(0xbc)
	mem.LastNonZero(0xde)

	mem.StartOrder0()

	if got := mem.Porta(0); got != 0 {
		t.Fatalf("expected porta reset to 0, got 0x%02x", got)
	}
	if vx, vy := mem.Vibrato(0); vx != 0 || vy != 0 {
		t.Fatalf("expected vibrato reset, got 0x%02x/0x%02x", vx, vy)
	}
	if vx, vy := mem.Tremolo(0); vx != 0 || vy != 0 {
		t.Fatalf("expected tremolo reset, got 0x%02x/0x%02x", vx, vy)
	}
	if got := mem.SampleOffset(0); got != 0 {
		t.Fatalf("expected sample offset reset, got 0x%02x", got)
	}
	if got := mem.TempoDecrease(0); got != 0 {
		t.Fatalf("expected tempo decrease reset, got 0x%02x", got)
	}
	if got := mem.TempoIncrease(0); got != 0 {
		t.Fatalf("expected tempo increase reset, got 0x%02x", got)
	}
	if got := mem.LastNonZero(0); got != 0 {
		t.Fatalf("expected last non-zero reset, got 0x%02x", got)
	}
}
