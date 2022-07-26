package voice

import (
	"github.com/gotracker/playback/voice/oscillator"
)

// AutoVibrato is the setting and memory for the auto-vibrato system
type AutoVibrato struct {
	Enabled           bool
	Sweep             int
	WaveformSelection uint8
	Depth             float32
	Rate              int
	Factory           func() oscillator.Oscillator
}

// Generate creates an AutoVibrato waveform oscillator and configures it with the inital values
func (a *AutoVibrato) Generate() oscillator.Oscillator {
	if a.Factory == nil {
		return nil
	}
	o := a.Factory()
	o.SetWaveform(oscillator.WaveTableSelect(a.WaveformSelection))
	return o
}
