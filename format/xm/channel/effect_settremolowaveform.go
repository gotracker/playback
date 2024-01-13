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

// SetTremoloWaveform defines a set tremolo waveform effect
type SetTremoloWaveform[TPeriod period.Period] DataEffect // 'E7x'

func (e SetTremoloWaveform[TPeriod]) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e SetTremoloWaveform[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]) error {
	return m.SetChannelOscillatorWaveform(ch, machine.OscillatorTremolo, oscillator.WaveTableSelect(e&0x0F))
}

func (e SetTremoloWaveform[TPeriod]) TraceData() string {
	return e.String()
}
