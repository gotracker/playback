package effect

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
)

// SetGlobalVolume defines a set global volume effect
type SetGlobalVolume channel.DataEffect // 'V'

// PreStart triggers when the effect enters onto the channel state
func (e SetGlobalVolume) PreStart(cs *channel.State, p playback.Playback) error {
	v := volume.Volume(channel.DataEffect(e)) / 0x80
	if v > 1 {
		v = 1
	}
	cs.SetChannelVolume(v)
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetGlobalVolume) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetGlobalVolume) String() string {
	return fmt.Sprintf("V%0.2x", channel.DataEffect(e))
}
