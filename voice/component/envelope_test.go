package component

import (
	"testing"

	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/loop"
)

func TestVolumeEnvelopeAdvancesAndFinishes(t *testing.T) {
	var env VolumeEnvelope[testVolume]
	finished := false
	onFinished := voice.Callback(func(v voice.Voice) { finished = true })

	env.Setup(EnvelopeSettings[testVolume, testVolume]{
		Envelope: envelope.Envelope[testVolume]{
			Enabled: true,
			Loop:    &loop.Disabled{},
			Sustain: &loop.Disabled{},
			Length:  12,
			Values: []envelope.Point[testVolume]{
				{Pos: 0, Length: 2, Y: testVolume(0)},
				{Pos: 10, Length: 0, Y: testVolume(1)},
			},
		},
		OnFinished: onFinished,
	})

	if v := env.GetCurrentValue(); !almostEqualVol(v.ToVolume(), 0, 1e-6) {
		t.Fatalf("initial envelope value mismatch: got %v", v)
	}

	if cb := env.Advance(); cb != nil {
		t.Fatalf("expected no callback on first advance")
	}
	if v := env.GetCurrentValue(); !almostEqualVol(v.ToVolume(), 0.5, 1e-6) {
		t.Fatalf("after first advance expected 0.5, got %v", v)
	}

	if cb := env.Advance(); cb != nil {
		t.Fatalf("expected no callback on second advance")
	}
	if v := env.GetCurrentValue(); !almostEqualVol(v.ToVolume(), 1, 1e-6) {
		t.Fatalf("after second advance expected 1.0, got %v", v)
	}

	var cb voice.Callback
	for i := 0; i < 16 && cb == nil; i++ {
		cb = env.Advance()
	}
	if cb == nil {
		t.Fatalf("expected callback when envelope finishes")
	}

	cb(nil)
	if !finished {
		t.Fatalf("expected returned callback to set finished flag")
	}
}

func TestEnvelopeHandlesNilLoopsDeterministically(t *testing.T) {
	var env VolumeEnvelope[testVolume]
	finished := false
	onFinished := voice.Callback(func(v voice.Voice) { finished = true })

	env.Setup(EnvelopeSettings[testVolume, testVolume]{
		Envelope: envelope.Envelope[testVolume]{
			Enabled: true,
			Length:  4,
			Values: []envelope.Point[testVolume]{
				{Pos: 0, Length: 2, Y: testVolume(0)},
				{Pos: 2, Length: 1, Y: testVolume(1)},
			},
		},
		OnFinished: onFinished,
	})

	if v := env.GetCurrentValue(); !almostEqualVol(v.ToVolume(), 0, 1e-6) {
		t.Fatalf("initial value mismatch: got %v", v)
	}

	_ = env.Advance()
	if v := env.GetCurrentValue(); v.ToVolume() <= 0 {
		t.Fatalf("expected progression after advance, got %v", v)
	}

	var cb voice.Callback
	for i := 0; i < 8 && cb == nil; i++ {
		cb = env.Advance()
	}
	if cb == nil {
		t.Fatalf("expected envelope to finish")
	}
	cb(nil)
	if !finished {
		t.Fatalf("expected finished callback to run")
	}
}
