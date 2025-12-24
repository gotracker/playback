package autovibrato

import (
	"errors"
	"testing"

	"github.com/gotracker/playback/voice/oscillator"
)

type fakeOsc struct {
	wave  oscillator.WaveTableSelect
	clone bool
}

func (f *fakeOsc) Clone() oscillator.Oscillator                 { f.clone = true; return f }
func (f *fakeOsc) GetWave(depth float32) float32                { return float32(f.wave) + depth }
func (f *fakeOsc) Advance(speed int)                            {}
func (f *fakeOsc) SetWaveform(table oscillator.WaveTableSelect) { f.wave = table }
func (f *fakeOsc) GetWaveform() oscillator.WaveTableSelect      { return f.wave }
func (f *fakeOsc) HardReset()                                   {}
func (f *fakeOsc) Reset()                                       {}

type testPeriod struct{}

func (testPeriod) IsInvalid() bool { return false }

func TestGenerateUsesFactoryName(t *testing.T) {
	cfg := AutoVibratoConfig[testPeriod]{FactoryName: "sine", WaveformSelection: 2}
	osc, err := cfg.Generate(func(name string) (oscillator.Oscillator, error) {
		if name != "sine" {
			return nil, errors.New("bad name")
		}
		return &fakeOsc{}, nil
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if osc == nil {
		t.Fatalf("expected oscillator instance")
	}
	if osc.GetWaveform() != 2 {
		t.Fatalf("expected waveform 2, got %d", osc.GetWaveform())
	}
}

func TestGenerateReturnsNilWhenFactoryNil(t *testing.T) {
	cfg := AutoVibratoConfig[testPeriod]{}
	osc, err := cfg.Generate(nil)
	if err != nil {
		t.Fatalf("expected nil error when factory nil")
	}
	if osc != nil {
		t.Fatalf("expected nil oscillator when factory nil")
	}
}

func TestGeneratePropagatesFactoryError(t *testing.T) {
	expect := errors.New("boom")
	cfg := AutoVibratoConfig[testPeriod]{FactoryName: "saw"}
	_, err := cfg.Generate(func(string) (oscillator.Oscillator, error) {
		return nil, expect
	})
	if !errors.Is(err, expect) {
		t.Fatalf("expected propagated error, got %v", err)
	}
}
