package channel

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/note"
)

type Stop struct {
	s note.Semitone
	i InstID
}

func stop() EffectS3M {
	return Stop{}
}

// Start triggers on the first tick, but before the Tick() function is called
func (e Stop) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	if inst := cs.GetInstrument(); inst != nil {
		if inst.IsReleaseNote(note.StopOrReleaseNote{}) {
			if v := cs.GetVoice(); v != nil {
				v.Release()
			}
		} else {
			cs.NoteCut()
		}
	}

	return nil
}

func (e Stop) String() string {
	return "=== .."
}
