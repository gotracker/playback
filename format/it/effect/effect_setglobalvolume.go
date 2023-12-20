package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume[TPeriod period.Period] channel.DataEffect // 'V'

// PreStart triggers when the effect enters onto the channel state
func (e SetGlobalVolume[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	v := volume.Volume(channel.DataEffect(e)) / 0x80
	if v > 1 {
		v = 1
	}
	cs.SetChannelVolume(v)
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetGlobalVolume[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetGlobalVolume[TPeriod]) String() string {
	return fmt.Sprintf("V%0.2x", channel.DataEffect(e))
}
