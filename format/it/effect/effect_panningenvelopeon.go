package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/period"
)

// PanningEnvelopeOn defines a panning envelope: on effect
type PanningEnvelopeOn[TPeriod period.Period] channel.DataEffect // 'S7A'

// Start triggers on the first tick, but before the Tick() function is called
func (e PanningEnvelopeOn[TPeriod]) Start(cs playback.Channel[TPeriod, channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	cs.SetPanningEnvelopeEnable(true)
	return nil
}

func (e PanningEnvelopeOn[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
