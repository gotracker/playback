package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/voice/oscillator"
)

// SetVibratoWaveform defines a set vibrato waveform effect
type SetVibratoWaveform[TPeriod period.Period] DataEffect // 'E4x'

func (e SetVibratoWaveform[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e SetVibratoWaveform[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	return m.SetChannelOscillatorWaveform(ch, machine.OscillatorVibrato, oscillator.WaveTableSelect(e&0xf))
}

func (e SetVibratoWaveform[TPeriod]) TraceData() string {
	return e.String()
}
