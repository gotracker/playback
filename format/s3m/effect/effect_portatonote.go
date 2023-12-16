package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote ChannelCommand // 'G'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaToNote) Start(cs playback.Channel[channel.Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	if cmd := cs.GetData(); cmd != nil && cmd.HasNote() {
		cs.SetPortaTargetPeriod(cs.GetTargetPeriod())
		cs.SetNotePlayTick(false, note.ActionContinue, 0)
	}
	return nil
}

// Tick is called on every tick
func (e PortaToNote) Tick(cs playback.Channel[channel.Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaToNote(channel.DataEffect(e))

	// vibrato modifies current period for portamento
	cur := cs.GetPeriod()
	if cur == nil {
		return nil
	}
	cur = cur.AddDelta(cs.GetPeriodDelta())
	ptp := cs.GetPortaTargetPeriod()
	if currentTick != 0 {
		if period.ComparePeriods(cur, ptp) == comparison.SpaceshipRightGreater {
			return doPortaUpToNote(cs, float32(xx), 4, ptp) // subtracts
		} else {
			return doPortaDownToNote(cs, float32(xx), 4, ptp) // adds
		}
	}
	return nil
}

func (e PortaToNote) String() string {
	return fmt.Sprintf("G%0.2x", channel.DataEffect(e))
}
