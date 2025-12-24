package period

import (
	"testing"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/period"
)

func TestCalcFinetuneC4SampleRateTable(t *testing.T) {
	cases := []struct {
		name     string
		finetune uint8
		expect   frequency.Frequency
	}{
		{"base", 0x0, 7895},
		{"center", 0x8, 8363},
		{"max", 0xF, 8757},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcFinetuneC4SampleRate(tt.finetune); got != tt.expect {
				t.Fatalf("finetune 0x%X -> %v, want %v", tt.finetune, got, tt.expect)
			}
		})
	}
}

func TestCalcFinetuneC4SampleRatePanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid finetune")
		}
	}()
	_ = CalcFinetuneC4SampleRate(0x10)
}

func TestAmigaConvertersConfigured(t *testing.T) {
	s3m, ok := S3MAmigaConverter.(period.AmigaConverter)
	if !ok {
		t.Fatalf("expected S3MAmigaConverter to be period.AmigaConverter")
	}
	if !s3m.SlideTo0Allowed {
		t.Fatalf("expected SlideTo0Allowed to be true")
	}
	if s3m.MinPeriod != 64 || s3m.MaxPeriod != 32767 {
		t.Fatalf("unexpected S3MAmigaConverter bounds: min %d max %d", s3m.MinPeriod, s3m.MaxPeriod)
	}

	mod, ok := MODAmigaConverter.(period.AmigaConverter)
	if !ok {
		t.Fatalf("expected MODAmigaConverter to be period.AmigaConverter")
	}
	if mod.MinPeriod != 56 || mod.MaxPeriod != 13696 {
		t.Fatalf("unexpected MODAmigaConverter bounds: min %d max %d", mod.MinPeriod, mod.MaxPeriod)
	}
}
