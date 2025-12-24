package load

import (
	"bytes"
	"testing"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"

	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
)

func TestModuleHeaderToHeaderNil(t *testing.T) {
	if _, err := moduleHeaderToHeader(nil); err == nil {
		t.Fatalf("expected error when module header is nil")
	}
}

func TestConvertSCRSFullToInstrumentNilAncillary(t *testing.T) {
	scrs := &s3mfile.SCRSFull{}
	if _, err := convertSCRSFullToInstrument(scrs, false, nil); err == nil {
		t.Fatalf("expected error when SCRS ancillary is nil")
	}
}

func TestConvertSCRSFullToInstrumentClampsVolume(t *testing.T) {
	scrs := &s3mfile.SCRSFull{
		SCRS: s3mfile.SCRS{
			Ancillary: &s3mfile.SCRSNoneHeader{Volume: s3mfile.Volume(0x80), C2Spd: s3mfile.HiLo32{Lo: 1234}},
		},
	}

	inst, err := convertSCRSFullToInstrument(scrs, false, nil)
	if err != nil {
		t.Fatalf("unexpected error converting SCRS: %v", err)
	}
	if inst.Static.Volume != s3mVolume.MaxVolume {
		t.Fatalf("expected volume clamp to MaxVolume, got %v", inst.Static.Volume)
	}
	if inst.SampleRate != frequency.Frequency(1234) {
		t.Fatalf("expected sample rate 1234, got %v", inst.SampleRate)
	}
}

func TestMODLoaderRejectsInvalidData(t *testing.T) {
	if _, err := MOD(bytes.NewReader([]byte("bad")), nil); err == nil {
		t.Fatalf("expected error for invalid MOD data")
	}
}
