package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/voice/oscillator"
)

// SetVibratoWaveform defines a set vibrato waveform effect
type SetVibratoWaveform ChannelCommand // 'S3x'

func (e SetVibratoWaveform) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetVibratoWaveform) RowStart(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) error {
	x := DataEffect(e) & 0x0f
	return m.SetChannelOscillatorWaveform(ch, machine.OscillatorVibrato, oscillator.WaveTableSelect(x))
}

func (e SetVibratoWaveform) TraceData() string {
	return e.String()
}
