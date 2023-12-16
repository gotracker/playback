package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
)

// SetVolume defines a volume slide effect
type SetVolume channel.DataEffect // 'C'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetVolume) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := xmVolume.XmVolume(e)

	cs.SetActiveVolume(xx.Volume())
	return nil
}

func (e SetVolume) String() string {
	return fmt.Sprintf("C%0.2x", channel.DataEffect(e))
}
