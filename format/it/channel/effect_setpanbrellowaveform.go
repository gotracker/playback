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

// SetPanbrelloWaveform defines a set panbrello waveform effect
type SetPanbrelloWaveform[TPeriod period.Period] DataEffect // 'S5x'

func (e SetPanbrelloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}

func (e SetPanbrelloWaveform[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	x := DataEffect(e) & 0x0f
	return m.SetChannelOscillatorWaveform(ch, machine.OscillatorPanbrello, oscillator.WaveTableSelect(x))
}

func (e SetPanbrelloWaveform[TPeriod]) TraceData() string {
	return e.String()
}
