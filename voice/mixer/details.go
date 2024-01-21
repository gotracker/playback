package mixer

import (
	"time"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/playback/frequency"
)

type Details struct {
	Mix        *mixing.Mixer
	Panmixer   mixing.PanMixer
	SampleRate frequency.Frequency
	Samples    int
	Duration   time.Duration
}
