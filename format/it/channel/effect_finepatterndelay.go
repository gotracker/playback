package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// FinePatternDelay defines an fine pattern delay effect
type FinePatternDelay[TPeriod period.Period] DataEffect // 'S6x'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePatternDelay[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := DataEffect(e) & 0xf

	m := p.(IT)
	if err := m.AddRowTicks(int(x)); err != nil {
		return err
	}
	return nil
}

func (e FinePatternDelay[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
