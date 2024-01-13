package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/voice/oscillator"
)

// SetVibratoWaveform defines a set vibrato waveform effect
type SetVibratoWaveform[TPeriod period.Period] DataEffect // 'S3x'

func (e SetVibratoWaveform[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetVibratoWaveform[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	x := DataEffect(e) & 0x0f
	return m.SetChannelOscillatorWaveform(ch, machine.OscillatorVibrato, oscillator.WaveTableSelect(x))
}

func (e SetVibratoWaveform[TPeriod]) TraceData() string {
	return e.String()
}
