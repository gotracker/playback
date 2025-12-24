package volume

import "testing"

func TestWithOverflowProtectionClamps(t *testing.T) {
	cases := []struct {
		name string
		in   Volume
		exp  float64
	}{
		{"within", 0.5, 0.5},
		{"pos overflow", 2, 1},
		{"neg overflow", -2, -1},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.WithOverflowProtection(); got != tt.exp {
				t.Fatalf("WithOverflowProtection() = %v, want %v", got, tt.exp)
			}
		})
	}
}

func TestToSampleConversions(t *testing.T) {
	v := Volume(0.5)
	if got := v.ToSample(8).(int8); got != 64 {
		t.Fatalf("8-bit sample = %d, want 64", got)
	}
	if got := v.ToSample(16).(int16); got != 16339 {
		t.Fatalf("16-bit sample = %d, want 16339", got)
	}
	if got := Volume(2).ToSample(16).(int16); got != 32678 {
		t.Fatalf("16-bit clamped sample = %d, want 32678", got)
	}
}

func TestToUintSampleConversions(t *testing.T) {
	if got := Volume(0).ToUintSample(8); got != 0x80 {
		t.Fatalf("uint8 sample = %d, want 128", got)
	}
	if got := Volume(1).ToUintSample(8); got != 0xFF {
		t.Fatalf("uint8 sample at max = %d, want 255", got)
	}
}

func TestApplyHelpers(t *testing.T) {
	v := Volume(0.5)
	if got := v.ApplySingle(0.8); got != Volume(0.4) {
		t.Fatalf("ApplySingle = %v, want 0.4", got)
	}

	in := []Volume{1, -1}
	out := v.ApplyMultiple(in)
	if len(out) != 2 || out[0] != 0.5 || out[1] != -0.5 {
		t.Fatalf("ApplyMultiple = %#v, want [0.5 -0.5]", out)
	}
}
