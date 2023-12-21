package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// SetFinetune defines a mod-style set finetune effect
type SetFinetune[TPeriod period.Period] DataEffect // 'S2x'

// PreStart triggers when the effect enters onto the channel state
func (e SetFinetune[TPeriod]) PreStart(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	x := DataEffect(e) & 0xf

	inst := cs.GetTargetInst()
	if inst != nil {
		ft := (note.Finetune(x) - 8) * 4
		inst.SetFinetune(ft)
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetFinetune[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetFinetune[TPeriod]) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
