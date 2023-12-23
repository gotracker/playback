package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// SetEnvelopePosition defines a set envelope position effect
type SetEnvelopePosition[TPeriod period.Period] DataEffect // 'Lxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetEnvelopePosition[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := DataEffect(e)

	cs.SetEnvelopePosition(int(xx))
	return nil
}

func (e SetEnvelopePosition[TPeriod]) String() string {
	return fmt.Sprintf("L%0.2x", DataEffect(e))
}
