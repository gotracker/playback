package channel

import (
	"fmt"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
)

// SetTremoloWaveform defines a set tremolo waveform effect
type SetTremoloWaveform[TPeriod period.Period] DataEffect // 'S4x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTremoloWaveform[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	mem := cs.GetMemory()
	trem := mem.TremoloOscillator()
	trem.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetTremoloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
