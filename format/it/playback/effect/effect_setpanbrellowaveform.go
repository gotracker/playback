package effect

import (
	"fmt"

	"github.com/gotracker/voice/oscillator"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// SetPanbrelloWaveform defines a set panbrello waveform effect
type SetPanbrelloWaveform channel.DataEffect // 'S5x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanbrelloWaveform) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	mem := cs.GetMemory()
	panb := mem.PanbrelloOscillator()
	panb.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetPanbrelloWaveform) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
