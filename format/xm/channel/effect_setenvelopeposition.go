package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// SetEnvelopePosition defines a set envelope position effect
type SetEnvelopePosition[TPeriod period.Period] DataEffect // 'Lxx'

func (e SetEnvelopePosition[TPeriod]) String() string {
	return fmt.Sprintf("L%0.2x", DataEffect(e))
}

func (e SetEnvelopePosition[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	return m.SetChannelEnvelopePositions(ch, int(e))
}

func (e SetEnvelopePosition[TPeriod]) TraceData() string {
	return e.String()
}
