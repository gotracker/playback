package channel

import (
	"fmt"

	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
)

// SetTremoloWaveform defines a set tremolo waveform effect
type SetTremoloWaveform ChannelCommand // 'S4x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetTremoloWaveform) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	mem := cs.GetMemory()
	trem := mem.TremoloOscillator()
	trem.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetTremoloWaveform) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
