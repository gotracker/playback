package pcm

import "testing"

func TestGetSampleBytes(t *testing.T) {
	cases := []struct {
		fmt SampleDataFormat
		exp int
	}{
		{SampleDataFormat8BitUnsigned, 1},
		{SampleDataFormat8BitSigned, 1},
		{SampleDataFormat16BitLEUnsigned, 2},
		{SampleDataFormat16BitBESigned, 2},
		{SampleDataFormat32BitLEFloat, 4},
		{SampleDataFormat64BitBEFloat, 8},
		{SampleDataFormatNative, 1},
	}

	for _, tt := range cases {
		if got := getSampleBytes(tt.fmt); got != tt.exp {
			t.Fatalf("fmt %v => %d, want %d", tt.fmt, got, tt.exp)
		}
	}
}
