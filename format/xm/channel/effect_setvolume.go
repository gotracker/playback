package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

// SetVolume defines a volume slide effect
type SetVolume[TPeriod period.Period] DataEffect // 'C'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetVolume[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := xmVolume.XmVolume(e)

	active := cs.GetActiveState()
	active.SetVolume(xx.Volume())
	return nil
}

func (e SetVolume[TPeriod]) String() string {
	return fmt.Sprintf("C%0.2x", DataEffect(e))
}
