package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

// SetEnvelopePosition defines a set envelope position effect
type SetEnvelopePosition channel.DataEffect // 'Lxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetEnvelopePosition) Start(cs *channel.State, p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := channel.DataEffect(e)

	cs.SetEnvelopePosition(int(xx))
	return nil
}

func (e SetEnvelopePosition) String() string {
	return fmt.Sprintf("L%0.2x", channel.DataEffect(e))
}
