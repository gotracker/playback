package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote ChannelCommand // 'G'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaToNote) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	if cmd := cs.GetChannelData(); cmd.HasNote() {
		cs.SetPortaTargetPeriod(cs.GetTargetPeriod())
		cs.SetNotePlayTick(false, note.ActionContinue, 0)
	}
	return nil
}

// Tick is called on every tick
func (e PortaToNote) Tick(cs S3MChannel, p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaToNote(DataEffect(e))

	// vibrato modifies current period for portamento
	cur := cs.GetPeriod()
	if cur.IsInvalid() {
		return nil
	}
	d := cs.GetPeriodDelta()
	cur = period.AddDelta(cur, d)
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
	return fmt.Sprintf("G%0.2x", DataEffect(e))
}
