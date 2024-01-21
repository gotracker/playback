package channel

import (
	"fmt"

	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// PanSlide defines a pan slide effect
type PanSlide[TPeriod period.Period] DataEffect // 'Pxx'

func (e PanSlide[TPeriod]) String() string {
	return fmt.Sprintf("P%0.2x", DataEffect(e))
}

func (e PanSlide[TPeriod]) Tick(ch index.Channel, m machine.Machine[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning], tick int) error {
	xx := DataEffect(e)
	x, y := xx>>4, xx&0x0F

	if x == 0 {
		// slide left y units
		if err := m.SlideChannelPan(ch, 1, -float32(y)); err != nil {
			return err
		}
	} else if y == 0 {
		// slide right x units
		if err := m.SlideChannelPan(ch, 1, float32(x)); err != nil {
			return err
		}
	}
	return nil
}

func (e PanSlide[TPeriod]) TraceData() string {
	return e.String()
}
