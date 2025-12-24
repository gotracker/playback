package volume

import (
	"testing"

	mixvol "github.com/gotracker/playback/mixing/volume"
)

func TestXmVolumeConversionsAndClamps(t *testing.T) {
	if got := XmVolume(0x40).ToVolume(); got != 1 {
		t.Fatalf("expected max xm volume to map to 1, got %v", got)
	}
	if XmVolume(0xff).IsUseInstrumentVol() == false {
		t.Fatalf("expected use-instrument sentinel")
	}
	if XmVolume(0x41).IsInvalid() == false {
		t.Fatalf("expected >0x40 (non-sentinel) invalid")
	}

	if got := XmVolume(0x7f).FMA(2, 0); got != 0x40 {
		t.Fatalf("expected FMA clamp to 0x40, got %d", got)
	}
	if got := XmVolume(1).AddDelta(-5); got != 0 {
		t.Fatalf("expected AddDelta clamp to 0, got %d", got)
	}
	if got := XmVolume(0x40).AddDelta(5); got != 0x40 {
		t.Fatalf("expected AddDelta clamp to max, got %d", got)
	}

	if got := ToVolumeXM(mixvol.VolumeUseInstVol); got != XmVolume(0xff) {
		t.Fatalf("expected sentinel round trip for VolumeUseInstVol")
	}
}
