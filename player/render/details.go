package render

import (
	"time"

	"github.com/gotracker/gomixing/mixing"
)

type Details struct {
	Mix          *mixing.Mixer
	Panmixer     mixing.PanMixer
	SamplerSpeed float32
	Samples      int
	Duration     time.Duration
}
