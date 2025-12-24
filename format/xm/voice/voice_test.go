package voice

import (
	"testing"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmSystem "github.com/gotracker/playback/format/xm/system"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/period"
	voiceCore "github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/autovibrato"
	"github.com/gotracker/playback/voice/pcm"
)

func makeXMInstrument() instrument.Instrument[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning] {
	sample := pcm.NewSampleNative([]volume.Matrix{{StaticMatrix: volume.StaticMatrix{1, 1}, Channels: 2}}, 1, 2)

	return instrument.Instrument[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		Static: instrument.StaticValues[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
			PC:     xmPeriod.LinearConverter,
			Volume: xmVolume.XmVolume(32),
			AutoVibrato: autovibrato.AutoVibratoConfig[period.Linear]{
				PC:          xmPeriod.LinearConverter,
				FactoryName: "vibrato",
			},
		},
		Inst:       &instrument.PCM[xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{Sample: sample},
		SampleRate: xmSystem.DefaultC4SampleRate,
	}
}

type xmTestVoice interface {
	voiceCore.RenderSampler[period.Linear]
	Setup(*instrument.Instrument[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error
}

func makeXMVoice() xmTestVoice {
	cfg := voiceCore.VoiceConfig[period.Linear, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]{
		PC:            xmPeriod.LinearConverter,
		InitialVolume: xmVolume.XmVolume(32),
		InitialMixing: xmVolume.DefaultXmMixingVolume,
		PanEnabled:    true,
		InitialPan:    xmPanning.DefaultPanning,
	}
	return New[period.Linear](cfg).(xmTestVoice)
}

func TestXMVoiceSetupAndSample(t *testing.T) {
	v := makeXMVoice()
	inst := makeXMInstrument()

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

func TestXMVoiceStopMarksDone(t *testing.T) {
	v := makeXMVoice()
	inst := makeXMInstrument()

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
