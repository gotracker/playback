package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
)

// SetFinetune defines a mod-style set finetune effect
type SetFinetune ChannelCommand // 'S2x'

// PreStart triggers when the effect enters onto the channel state
func (e SetFinetune) PreStart(cs S3MChannel, p playback.Playback) error {
	x := DataEffect(e) & 0xf

	inst := cs.GetTargetInst()
	if inst != nil {
		ft := (note.Finetune(x) - 8) * 4
		inst.SetFinetune(ft)
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e SetFinetune) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()
	return nil
}

func (e SetFinetune) String() string {
	return fmt.Sprintf("S%0.2x", DataEffect(e))
}
