package autovibrato

import (
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/gotracker/playback/voice/types"
)

// AutoVibratoConfig is the setting and memory for the auto-vibrato system
type AutoVibratoConfig[TPeriod types.Period] struct {
	PC                period.PeriodConverter[TPeriod]
	Enabled           bool
	Sweep             int
	WaveformSelection uint8
	Depth             float32
	Rate              int
	FactoryName       string
}

// Generate creates an AutoVibrato waveform oscillator and configures it with the inital values
func (a AutoVibratoConfig[TPeriod]) Generate(factory func(string) (oscillator.Oscillator, error)) (oscillator.Oscillator, error) {
	if factory == nil {
		return nil, nil
	}

	o, err := factory(a.FactoryName)
	if err != nil {
		return nil, err
	}

	o.SetWaveform(oscillator.WaveTableSelect(a.WaveformSelection))
	return o, nil
}
