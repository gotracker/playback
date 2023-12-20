package effect

import (
	"fmt"

	"github.com/gotracker/playback/voice/oscillator"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
)

// SetVibratoWaveform defines a set vibrato waveform effect
type SetVibratoWaveform ChannelCommand // 'S3x'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetVibratoWaveform) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	mem := cs.GetMemory()
	vib := mem.VibratoOscillator()
	vib.SetWaveform(oscillator.WaveTableSelect(x))
	return nil
}

func (e SetVibratoWaveform) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
