package voice

import (
	"testing"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itPeriod "github.com/gotracker/playback/format/it/period"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
	voiceCore "github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/pcm"
)

func makeTestInstrument() instrument.Instrument[period.Linear, itVolume.FineVolume, itVolume.Volume, itPanning.Panning] {
	sample := pcm.NewSampleNative([]volume.Matrix{{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2}}, 1, 2)

	return instrument.Instrument[period.Linear, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		Static: instrument.StaticValues[period.Linear, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
			PC:     itPeriod.LinearConverter,
			Volume: itVolume.Volume(32),
			AutoVibrato: autovibrato.AutoVibratoConfig[period.Linear]{
				PC:          itPeriod.LinearConverter,
				FactoryName: "vibrato",
			},
		},
		Inst:       &instrument.PCM[itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{Sample: sample},
		SampleRate: 8363,
	}
}

type testVoice interface {
	voiceCore.RenderSampler[period.Linear]
	Setup(*instrument.Instrument[period.Linear, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error
}

func makeVoice() testVoice {
	cfg := voiceCore.VoiceConfig[period.Linear, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
		PC:            itPeriod.LinearConverter,
		InitialVolume: itVolume.Volume(32),
		InitialMixing: itVolume.FineVolume(64),
		PanEnabled:    true,
		InitialPan:    itPanning.DefaultPanning,
	}
	return New[period.Linear](cfg).(testVoice)
}

func TestVoiceSetupAndSample(t *testing.T) {
	v := makeVoice()
	inst := makeTestInstrument()

	if err := v.Setup(&inst); err != nil {
		t.Fatalf("voice setup error: %v", err)
	}
	if v.IsDone() {
		t.Fatalf("expected voice active after setup")
	}

	if err := v.Tick(); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	samp := v.GetSample(sampling.Pos{})
	if samp.Channels != 2 {
		t.Fatalf("expected stereo sample, got %d channels", samp.Channels)
	}
	if fv := v.GetFinalVolume(); fv <= 0 {
		t.Fatalf("expected final volume > 0, got %v", fv)
	}
}

func TestVoiceStopMarksDone(t *testing.T) {
	v := makeVoice()
	inst := makeTestInstrument()

	if err := v.Setup(&inst); err != nil {
		t.Fatalf("voice setup error: %v", err)
	}
	v.Stop()
	if !v.IsDone() {
		t.Fatalf("expected voice to be done after stop")
	}
	if fv := v.GetFinalVolume(); fv != 0 {
		t.Fatalf("expected final volume to be 0 after stop, got %v", fv)
	}
}
