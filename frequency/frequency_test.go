package frequency

import "testing"

func TestFrequencyGoStringFormatsUnits(t *testing.T) {
	cases := []struct {
		name string
		in   Frequency
		exp  string
	}{
		{"hz", 440, "440.000000Hz"},
		{"khz", 20000, "20000kHz"},
		{"mhz", 3_200_000, "3200000MHz"},
		{"ghz", 5_000_000_000, "5000000000GHz"},
		{"thz", 7_000_000_000_000, "7000000000000THz"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.GoString(); got != tt.exp {
				t.Fatalf("GoString() = %q, want %q", got, tt.exp)
			}
		})
	}
}
