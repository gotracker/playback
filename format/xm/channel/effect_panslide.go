package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	"github.com/gotracker/playback/period"
)

// PanSlide defines a pan slide effect
type PanSlide[TPeriod period.Period] DataEffect // 'Pxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e PanSlide[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	xx := DataEffect(e)
	x := xx >> 4
	y := xx & 0x0F

	xp := DataEffect(xmPanning.PanningToXm(cs.GetPan()))
	if x == 0 {
		// slide left y units
		if xp < y {
			xp = 0
		} else {
			xp -= y
		}
	} else if y == 0 {
		// slide right x units
		if xp > 0xFF-x {
			xp = 0xFF
		} else {
			xp += x
		}
	}
	cs.SetPan(xmPanning.PanningFromXm(uint8(xp)))
	return nil
}

func (e PanSlide[TPeriod]) String() string {
	return fmt.Sprintf("P%0.2x", DataEffect(e))
}
