package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ExtraFinePortaDown defines an extra-fine portamento down effect
type ExtraFinePortaDown ChannelCommand // 'EEx'

func (e ExtraFinePortaDown) String() string {
	return fmt.Sprintf("E%0.2x", DataEffect(e))
}

func (e ExtraFinePortaDown) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	y := DataEffect(e) & 0x0F
	return m.DoChannelPortaDown(ch, period.Delta(y)*1)
}

func (e ExtraFinePortaDown) TraceData() string {
	return e.String()
}
