package channel

import (
	"fmt"

	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mSystem "github.com/gotracker/playback/format/s3m/system"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// ExtraFinePortaUp defines an extra-fine portamento up effect
type ExtraFinePortaUp ChannelCommand // 'FEx'

func (e ExtraFinePortaUp) String() string {
	return fmt.Sprintf("F%0.2x", DataEffect(e))
}

func (e ExtraFinePortaUp) Tick(ch index.Channel, m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], tick int) error {
	if tick != 0 {
		return nil
	}

	y := DataEffect(e) & 0x0F
	return m.DoChannelPortaUp(ch, period.Delta(y)*1*s3mSystem.SlideFinesPerSemitone)
}

func (e ExtraFinePortaUp) TraceData() string {
	return e.String()
}
