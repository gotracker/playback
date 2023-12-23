package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote[TPeriod period.Period] DataEffect // '3'

// Start triggers on the first tick, but before the Tick() function is called
func (e PortaToNote[TPeriod]) Start(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback) error {
	cs.ResetRetriggerCount()
	cs.UnfreezePlayback()
	if cmd := cs.GetChannelData(); cmd.HasNote() {
		cs.SetPortaTargetPeriod(cs.GetTargetPeriod())
		cs.SetNotePlayTick(false, note.ActionContinue, 0)
	}
	return nil
}

// Tick is called on every tick
func (e PortaToNote[TPeriod]) Tick(cs playback.Channel[TPeriod, Memory, Data], p playback.Playback, currentTick int) error {
	if currentTick == 0 {
		return nil
	}

	mem := cs.GetMemory()
	xx := mem.PortaToNote(DataEffect(e))

	current := cs.GetPeriod()
	target := cs.GetPortaTargetPeriod()
	if period.ComparePeriods(current, target) == comparison.SpaceshipRightGreater {
		return doPortaUpToNote(cs, float32(xx), 4, target) // subtracts
	} else {
		return doPortaDownToNote(cs, float32(xx), 4, target) // adds
	}
}

func (e PortaToNote[TPeriod]) String() string {
	return fmt.Sprintf("3%0.2x", DataEffect(e))
}
