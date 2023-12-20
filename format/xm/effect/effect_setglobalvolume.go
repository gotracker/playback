package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume[TPeriod period.Period] channel.DataEffect // 'G'

// PreStart triggers when the effect enters onto the channel state
func (e SetGlobalVolume[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	v := xmVolume.XmVolume(e)
	p.SetGlobalVolume(v.Volume())
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetGlobalVolume[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetGlobalVolume[TPeriod]) String() string {
	return fmt.Sprintf("G%0.2x", channel.DataEffect(e))
}
