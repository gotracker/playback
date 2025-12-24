package sampler

import "testing"

func TestSamplerMixerAccessors(t *testing.T) {
	s := NewSampler(44100, 2, 0.5, nil)

	m := s.Mixer()
	if m.Channels != 2 {
		t.Fatalf("expected mixer channels 2, got %d", m.Channels)
	}

	// mutating through Mixer pointer updates internal mixer
	m.Channels = 4
	if s.MixerConfig().Channels != 4 {
		t.Fatalf("expected MixerConfig to reflect pointer mutation to 4, got %d", s.MixerConfig().Channels)
	}

	// modifying returned config copy should not alter stored mixer
	cfg := s.MixerConfig()
	cfg.Channels = 1
	if s.Mixer().Channels != 4 {
		t.Fatalf("expected internal mixer to remain 4, got %d", s.Mixer().Channels)
	}
}

func TestSamplerGetPanMixer(t *testing.T) {
	s := NewSampler(44100, 2, 1, nil)
	pm := s.GetPanMixer()
	if pm == nil {
		t.Fatalf("expected stereo pan mixer")
	}
	if pm.NumChannels() != 2 {
		t.Fatalf("expected pan mixer with 2 channels, got %d", pm.NumChannels())
	}

	s.mixer.Channels = 3
	if s.GetPanMixer() != nil {
		t.Fatalf("expected nil pan mixer for unsupported channel count")
	}
}
