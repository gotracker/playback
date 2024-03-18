package mixer

import (
	"time"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/mixing"
)

type Details struct {
	Mix              *mixing.Mixer
	Panmixer         mixing.PanMixer
	SampleRate       frequency.Frequency
	StereoSeparation float32
	Samples          int
	Duration         time.Duration
}
