package component

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/voice"
)

// OutputFilter applies a filter to a sample stream
type OutputFilter struct {
	Input  sampling.SampleStream
	Output voice.FilterApplier
}

// GetSample operates the filter
func (o *OutputFilter) GetSample(pos sampling.Pos) volume.Matrix {
	dry := o.Input.GetSample(pos)
	return o.Output.ApplyFilter(dry)
}
