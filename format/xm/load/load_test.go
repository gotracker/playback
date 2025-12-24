package load

import (
	"testing"

	xmfile "github.com/gotracker/goaudiofile/music/tracked/xm"

	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

func TestModuleHeaderToHeaderNil(t *testing.T) {
	if _, err := moduleHeaderToHeader(nil, false); err == nil {
		t.Fatalf("expected error when module header is nil")
	}
}

func TestConvertXMInstrumentToInstrumentNil(t *testing.T) {
	if _, _, err := convertXMInstrumentToInstrument[period.Amiga](nil, xmPeriod.AmigaConverter, false, nil); err == nil {
		t.Fatalf("expected error when instrument is nil")
	}
}

func TestXMInstrumentVolumeClampAndVibratoScaling(t *testing.T) {
	inst := &xmfile.InstrumentHeader{
		SamplesCount: 1,
		VibratoDepth: 64,
		VibratoRate:  1,
		Samples: []xmfile.SampleHeader{{
			Length:             1,
			Volume:             0x80, // over max, should clamp to 0x40
			Finetune:           0,
			Flags:              0,
			Panning:            0x40,
			RelativeNoteNumber: 0,
			SampleData:         []byte{0},
		}},
	}

	samplesLinear, _, err := convertXMInstrumentToInstrument(inst, xmPeriod.AmigaConverter, true, nil)
	if err != nil {
		t.Fatalf("unexpected error converting instrument (linear): %v", err)
	}
	if len(samplesLinear) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samplesLinear))
	}
	if got := samplesLinear[0].Static.Volume; got != xmVolume.XmVolume(0x40) {
		t.Fatalf("expected volume clamp to 0x40, got %v", got)
	}
	if got := samplesLinear[0].Static.AutoVibrato.Depth; got != float32(64) {
		t.Fatalf("expected vibrato depth 64 with linear slides, got %v", got)
	}

	samplesAmiga, _, err := convertXMInstrumentToInstrument(inst, xmPeriod.AmigaConverter, false, nil)
	if err != nil {
		t.Fatalf("unexpected error converting instrument (amiga): %v", err)
	}
	if got := samplesAmiga[0].Static.AutoVibrato.Depth; got != float32(1) {
		t.Fatalf("expected vibrato depth scaled to 1 for amiga slides, got %v", got)
	}
}
