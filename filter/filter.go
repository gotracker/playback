package filter

import (
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/frequency"
)

// Filter is an interface to a filter
type Filter interface {
	Filter(volume.Matrix) volume.Matrix
	SetPlaybackRate(playbackRate frequency.Frequency)
	UpdateEnv(uint8)
	Clone() Filter
}

// Factory is a function type that builds a filter with an input parameter taking a value between 0 and 1
type Factory func(instrument frequency.Frequency) Filter
