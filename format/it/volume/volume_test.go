package volume

import (
	"math"
	"testing"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	mixvol "github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/voice/types"
)

func TestVolumeConversionsAndSentinels(t *testing.T) {
	if got := FromItVolume(itfile.Volume(MaxItVolume)); got != 1 {
		t.Fatalf("expected it volume 64 to convert to 1, got %v", got)
	}

	if !Volume(0xff).IsUseInstrumentVol() {
		t.Fatalf("expected 0xff to be use-instrument sentinel")
	}
	if Volume(65).IsInvalid() == false {
		t.Fatalf("expected volumes over 64 (except sentinel) to be invalid")
	}
	if Volume(0xff).IsInvalid() {
		t.Fatalf("sentinel volume should not be invalid")
	}

	if got := ToItVolume(mixvol.VolumeUseInstVol); got != Volume(0xff) {
		t.Fatalf("expected use-instrument sentinel")
	}
	if got := ToItVolume(-0.1); got != 0 {
		t.Fatalf("expected negative volumes to clamp to 0, got %d", got)
	}
	if got := ToItVolume(2); got != Volume(MaxItVolume) {
		t.Fatalf("expected volumes over 1 to clamp to MaxItVolume, got %d", got)
	}
	if got := ToItVolume(0.5); got != Volume(MaxItVolume/2) {
		t.Fatalf("expected 0.5 to convert to 32, got %d", got)
	}
}

func TestVolumeArithmeticClamps(t *testing.T) {
	if got := Volume(60).FMA(2, 0); got != Volume(MaxItVolume) {
		t.Fatalf("expected FMA clamp to MaxItVolume, got %d", got)
	}

	if got := Volume(0xff).FMA(2, 1); got != Volume(0xff) {
		t.Fatalf("expected FMA to preserve instrument sentinel, got %d", got)
	}

	if got := Volume(3).AddDelta(types.VolumeDelta(-5)); got != 0 {
		t.Fatalf("expected AddDelta to clamp at 0, got %d", got)
	}

	if got := Volume(60).AddDelta(types.VolumeDelta(10)); got != Volume(MaxItVolume) {
		t.Fatalf("expected AddDelta to clamp at MaxItVolume, got %d", got)
	}
}

func TestFineVolumeConversionsAndClamps(t *testing.T) {
	if FineVolume(0xff).IsInvalid() {
		t.Fatalf("sentinel fine volume should not be invalid")
	}
	if FineVolume(MaxItFineVolume+1).IsInvalid() == false {
		t.Fatalf("expected fine volumes over max (except sentinel) to be invalid")
	}
	if !FineVolume(0xff).IsUseInstrumentVol() {
		t.Fatalf("expected 0xff to be use-instrument sentinel")
	}

	if got := ToItFineVolume(mixvol.VolumeUseInstVol); got != FineVolume(0xff) {
		t.Fatalf("expected fine volume use-instrument sentinel, got %d", got)
	}
	if got := ToItFineVolume(-0.2); got != 0 {
		t.Fatalf("expected negative fine volumes to clamp to 0, got %d", got)
	}
	if got := ToItFineVolume(2); got != MaxItFineVolume {
		t.Fatalf("expected fine volumes over 1 to clamp to max, got %d", got)
	}

	if got := FineVolume(120).FMA(2, 0); got != MaxItFineVolume {
		t.Fatalf("expected fine FMA clamp to max, got %d", got)
	}
	if got := FineVolume(0xff).FMA(2, 1); got != FineVolume(0xff) {
		t.Fatalf("expected fine FMA to preserve sentinel, got %d", got)
	}

	if got := FineVolume(1).AddDelta(types.VolumeDelta(-5)); got != 0 {
		t.Fatalf("expected fine AddDelta to clamp at 0, got %d", got)
	}
	if got := FineVolume(125).AddDelta(types.VolumeDelta(10)); got != MaxItFineVolume {
		t.Fatalf("expected fine AddDelta to clamp at max, got %d", got)
	}
}

func TestFromVolPan(t *testing.T) {
	const epsilon = 1e-6

	if got := FromVolPan(32); math.Abs(float64(got-0.5)) > epsilon {
		t.Fatalf("expected volpan 32 to be ~0.5, got %v", got)
	}

	if got := FromVolPan(200); got != mixvol.VolumeUseInstVol {
		t.Fatalf("expected volpan over max to use instrument volume sentinel, got %v", got)
	}
}
