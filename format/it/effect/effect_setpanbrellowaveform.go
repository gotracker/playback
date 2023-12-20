package effect

import (
	"fmt"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// SetPanbrelloWaveform defines a set panbrello waveform effect
type SetPanbrelloWaveform[TPeriod period.Period] channel.DataEffect // 'S5x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetPanbrelloWaveform[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	mem := cs.GetMemory()
	panb := mem.PanbrelloOscillator()
	panb.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetPanbrelloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
