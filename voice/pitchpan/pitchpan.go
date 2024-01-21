package pitchpan

import "github.com/gotracker/playback/note"

type PitchPan struct {
	Enabled    bool
	Center     note.Semitone
	Separation float32
}
