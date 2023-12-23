package channel

import (
	"fmt"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
)

// SetPanbrelloWaveform defines a set panbrello waveform effect
type SetPanbrelloWaveform[TPeriod period.Period] DataEffect // 'S5x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanbrelloWaveform[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	mem := cs.GetMemory()
	panb := mem.PanbrelloOscillator()
	panb.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetPanbrelloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
