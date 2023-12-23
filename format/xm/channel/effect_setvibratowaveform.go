package channel

import (
	"fmt"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
)

// SetVibratoWaveform defines a set vibrato waveform effect
type SetVibratoWaveform[TPeriod period.Period] DataEffect // 'E4x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetVibratoWaveform[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	mem := cs.GetMemory()
	vib := mem.VibratoOscillator()
	vib.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetVibratoWaveform[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}
