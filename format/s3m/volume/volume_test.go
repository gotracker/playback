package volume

import (
	"testing"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	mixvol "github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/voice/types"
)

func TestVolumeConversionsAndSentinels(t *testing.T) {
	if got := VolumeFromS3M(MaxVolume); got < 0.98 || got > 1.0 {
		t.Fatalf("expected max S3M volume near 1, got %v", got)
	}
	if Volume(s3mfile.EmptyVolume).IsUseInstrumentVol() == false {
		t.Fatalf("expected empty volume sentinel to indicate instrument volume")
	}
	if Volume(0x41).IsInvalid() == false {
		t.Fatalf("expected >MaxVolume (except sentinel) invalid")
	}
	if Volume(s3mfile.EmptyVolume).IsInvalid() {
		t.Fatalf("sentinel should not be invalid")
	}

	if got := VolumeToS3M(mixvol.VolumeUseInstVol); got != Volume(s3mfile.EmptyVolume) {
		t.Fatalf("expected use-instrument sentinel")
	}
	if got := VolumeToS3M(2); got == MaxVolume*2 {
		t.Fatalf("expected 2*MaxVolume, got %d", got)
	}
}

func TestVolumeArithmeticClamps(t *testing.T) {
	if got := Volume(60).FMA(2, 0); got != MaxVolume {
		t.Fatalf("expected FMA clamp to MaxVolume, got %d", got)
	}
	if got := Volume(s3mfile.EmptyVolume).FMA(2, 1); got != Volume(s3mfile.EmptyVolume) {
		t.Fatalf("expected FMA to preserve sentinel, got %d", got)
	}
	if got := Volume(2).AddDelta(types.VolumeDelta(-5)); got != 0 {
		t.Fatalf("expected AddDelta clamp to 0, got %d", got)
	}
	if got := Volume(100).AddDelta(types.VolumeDelta(10)); got != MaxVolume {
		t.Fatalf("expected AddDelta clamp to max, got %d", got)
	}
}

func TestFineVolumeConversionsAndClamps(t *testing.T) {
	if FineVolume(s3mfile.EmptyVolume).IsUseInstrumentVol() == false {
		t.Fatalf("expected fine volume sentinel")
	}
	if FineVolume(0x80).IsInvalid() == false {
		t.Fatalf("expected fine volume over max invalid")
	}
	if FineVolume(s3mfile.EmptyVolume).IsInvalid() {
		t.Fatalf("sentinel fine volume should not be invalid")
	}
	if got := FineVolume(s3mfile.EmptyVolume).FMA(2, 1); got != FineVolume(s3mfile.EmptyVolume) {
		t.Fatalf("expected fine FMA to preserve sentinel, got %d", got)
	}
	if got := FineVolume(0x7e).FMA(2, 0); got != MaxFineVolume {
		t.Fatalf("expected fine FMA clamp to max, got %d", got)
	}
	if got := FineVolume(1).AddDelta(types.VolumeDelta(-5)); got != 0 {
		t.Fatalf("expected fine AddDelta clamp to 0, got %d", got)
	}
	if got := FineVolume(0x7e).AddDelta(types.VolumeDelta(10)); got != MaxFineVolume {
		t.Fatalf("expected fine AddDelta clamp to max, got %d", got)
	}
}

func TestSampleVolumeConversions(t *testing.T) {
	if got := VolumeFromS3M8BitSample(128); got != 0 {
		t.Fatalf("expected centered 8-bit sample to map to 0, got %v", got)
	}
	if got := VolumeFromS3M16BitSample(32768); got != 0 {
		t.Fatalf("expected centered 16-bit sample to map to 0, got %v", got)
	}
}
