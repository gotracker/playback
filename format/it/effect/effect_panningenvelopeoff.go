package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// PanningEnvelopeOff defines a panning envelope: off effect
type PanningEnvelopeOff[TPeriod period.Period] channel.DataEffect // 'S79'

// Start triggers on the first tick, but before the Tick() function is called
func (e PanningEnvelopeOff[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetPanningEnvelopeEnable(false)
	return nil
}

func (e PanningEnvelopeOff[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
