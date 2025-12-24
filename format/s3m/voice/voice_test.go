package voice

import (
	"testing"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mSystem "github.com/gotracker/playback/format/s3m/system"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
	voiceCore "github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/pcm"
)

func makeS3MInstrument() instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	sample := pcm.NewSampleNative([]volume.Matrix{{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2}}, 1, 2)

	return instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		Static: instrument.StaticValues[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
			PC:     s3mPeriod.S3MAmigaConverter,
			Volume: s3mVolume.Volume(32),
		},
		Inst:       &instrument.PCM[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{Sample: sample},
		SampleRate: s3mSystem.DefaultC4SampleRate,
	}
}

type s3mTestVoice interface {
	voiceCore.RenderSampler[period.Amiga]
	Setup(*instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error
}

func makeS3MVoice() s3mTestVoice {
	cfg := voiceCore.VoiceConfig[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PC:            s3mPeriod.S3MAmigaConverter,
		InitialVolume: s3mVolume.Volume(32),
		InitialMixing: s3mVolume.FineVolume(64),
		PanEnabled:    true,
		InitialPan:    s3mPanning.DefaultPanning,
	}
	return New(cfg).(s3mTestVoice)
}

func TestS3MVoiceSetupAndSample(t *testing.T) {
	v := makeS3MVoice()
	inst := makeS3MInstrument()

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

func TestS3MVoiceStopMarksDone(t *testing.T) {
	v := makeS3MVoice()
	inst := makeS3MInstrument()

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
