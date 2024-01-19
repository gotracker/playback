package machine

import "github.com/gotracker/playback/index"

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) resetPatternLoop() {
	c.patternLoop.Total = 0
	c.patternLoop.Count = 0
}

// ContinueLoop returns the next expected row if a loop occurs
func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPatternLoop(ch index.Channel, m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	if c.patternLoop.Total == 0 {
		return nil
	}

	if m.ticker.current.Row != c.patternLoop.End {
		return nil
	}

	newCount := c.patternLoop.Count + 1
	doLoop := newCount <= c.patternLoop.Total
	if !doLoop {
		newCount = 0
	}

	traceChannelValueChangeWithComment(m, ch, "patternLoopCount", c.patternLoop.Count, newCount, "doPatternLoop")
	c.patternLoop.Count = newCount

	if doLoop {
		return m.SetRow(c.patternLoop.Start, false)
	}

	traceChannelValueChangeWithComment(m, ch, "patternLoopTotal", c.patternLoop.Total, 0, "doPatternLoop")
	c.patternLoop.Total = 0
	return nil
}
