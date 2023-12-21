package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote[TPeriod period.Period] DataEffect // 'G'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaToNote[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	if cmd := cs.GetData(); cmd != nil && cmd.HasNote() {
		cs.SetPortaTargetPeriod(cs.GetTargetPeriod())
		cs.SetNotePlayTick(false, note.ActionContinue, 0)
	}
	return nil
}

// Tick is called on every tick
func (e PortaToNote[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory], p playback.Playback, currentTick int) error {
	mem := cs.GetMemory()
	xx := mem.PortaToNote(DataEffect(e))

	// vibrato modifies current period for portamento
	cur := cs.GetPeriod()
	if cur == nil {
		return nil
	}
	switch pc := any(cur).(type) {
	case *itPeriod.Linear:
		cur = any(pc.Add(cs.GetPeriodDelta())).(*TPeriod)
	case *itPeriod.Amiga:
		cur = any(pc.Add(cs.GetPeriodDelta())).(*TPeriod)
	default:
		panic("unhandled period type")
	}
	ptp := cs.GetPortaTargetPeriod()
	if !mem.Shared.OldEffectMode || currentTick != 0 {
		if period.ComparePeriods(cur, ptp) == comparison.SpaceshipRightGreater {
			return doPortaUpToNote(cs, float32(xx), 4, ptp) // subtracts
		} else {
			return doPortaDownToNote(cs, float32(xx), 4, ptp) // adds
		}
	}
	return nil
}

func (e PortaToNote[TPeriod]) String() string {
	return fmt.Sprintf("G%0.2x", DataEffect(e))
}
