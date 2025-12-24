package autovibrato

import (
	"testing"

	"github.com/gotracker/playback/voice/oscillator"
)

type stubOsc struct{}

func (s *stubOsc) Clone() oscillator.Oscillator                 { return s }
func (s *stubOsc) GetWave(depth float32) float32                { return depth }
func (s *stubOsc) Advance(speed int)                            {}
func (s *stubOsc) SetWaveform(table oscillator.WaveTableSelect) {}
func (s *stubOsc) GetWaveform() oscillator.WaveTableSelect      { return 0 }
func (s *stubOsc) HardReset()                                   {}
func (s *stubOsc) Reset()                                       {}

func TestAutoVibratoSettingsFactoryStored(t *testing.T) {
	cfg := AutoVibratoSettings[testPeriod]{Factory: func(name string) (oscillator.Oscillator, error) {
		return &stubOsc{}, nil
	}}

	osc, err := cfg.Factory("any")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if osc == nil {
		t.Fatalf("expected oscillator instance")
	}
}
