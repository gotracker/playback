package effect

import (
	"fmt"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// SetTremoloWaveform defines a set tremolo waveform effect
type SetTremoloWaveform[TPeriod period.Period] channel.DataEffect // 'E7x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTremoloWaveform[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	mem := cs.GetMemory()
	trem := mem.TremoloOscillator()
	trem.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetTremoloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", channel.DataEffect(e))
}
