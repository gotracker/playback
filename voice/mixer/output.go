package mixer

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"
)

// Output applies a filter to a sample stream
type Output struct {
	Input  sampling.SampleStream
	Output ApplyFilter
}

// GetSample operates the filter
// must be pointer receiver
func (o *Output) GetSample(pos sampling.Pos) volume.Matrix {
	dry := o.Input.GetSample(pos)
	return o.Output.ApplyFilter(dry)
}

type ApplyFilter interface {
	ApplyFilter(dry volume.Matrix) volume.Matrix
}
