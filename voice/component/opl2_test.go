package component

import (
	"testing"

	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
)

// These tests focus on the math helpers that do not require a live OPL2 chip.
func TestOPL2Calc40(t *testing.T) {
	opl := OPL2[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]{}

	if got := opl.calc40(0x00, volume.Volume(1)); got != 0x00 {
		t.Fatalf("calc40(0x00, 1.0) = %d, want 0", got)
	}

	if got := opl.calc40(0x3f, volume.Volume(1)); got != 0x3f {
		t.Fatalf("calc40(0x3f, 1.0) = %d, want 63 (no attenuation)", got)
	}

	if got := opl.calc40(0x00, volume.Volume(0.5)); got != 0x1f && got != 0x20 {
		t.Fatalf("calc40(0x00, 0.5) = %d, want ~31-32", got)
	}
}

func TestOPL2FreqToFnumBlock(t *testing.T) {
	opl := OPL2[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume]{}

	fnum, block := opl.freqToFnumBlock(440.0)
	if block != 4 {
		t.Fatalf("freqToFnumBlock block = %d, want 4", block)
	}
	if fnum != 580 {
		t.Fatalf("freqToFnumBlock fnum = %d, want 580", fnum)
	}

	if fnum, block := opl.freqToFnumBlock(7000.0); fnum != 0 || block != 0 {
		t.Fatalf("freqToFnumBlock high freq = (%d,%d), want (0,0)", fnum, block)
	}
}
