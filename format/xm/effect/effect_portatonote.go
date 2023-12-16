package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/comparison"
)

// PortaToNote defines a portamento-to-note effect
type PortaToNote channel.DataEffect // '3'

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
	if currentTick == 0 {
		return nil
	}

	mem := cs.GetMemory()
	xx := mem.PortaToNote(channel.DataEffect(e))

	current := cs.GetPeriod()
	target := cs.GetPortaTargetPeriod()
	if period.ComparePeriods(current, target) == comparison.SpaceshipRightGreater {
		return doPortaUpToNote(cs, float32(xx), 4, target, mem.Shared.LinearFreqSlides) // subtracts
	} else {
		return doPortaDownToNote(cs, float32(xx), 4, target, mem.Shared.LinearFreqSlides) // adds
	}
}

func (e PortaToNote) String() string {
	return fmt.Sprintf("3%0.2x", channel.DataEffect(e))
}
