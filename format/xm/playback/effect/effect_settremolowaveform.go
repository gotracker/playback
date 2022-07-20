package effect

import (
	"fmt"

	"github.com/gotracker/voice/oscillator"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// SetTremoloWaveform defines a set tremolo waveform effect
type SetTremoloWaveform channel.DataEffect // 'E7x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTremoloWaveform) Start(cs playback.Channel[channel.Memory, channel.Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	mem := cs.GetMemory()
	trem := mem.TremoloOscillator()
	trem.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetTremoloWaveform) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
