package mixer

import (
	"time"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/playback/period"
)

type Details struct {
	Mix          *mixing.Mixer
	Panmixer     mixing.PanMixer
	SampleRate   period.Frequency
	SamplerSpeed float32
	Samples      int
	Duration     time.Duration
}
