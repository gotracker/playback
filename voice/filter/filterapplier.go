package filter

import (
	"github.com/gotracker/gomixing/volume"
)

// Applier is an interface for applying a filter to a sample stream
type Applier interface {
	ApplyFilter(dry volume.Matrix) volume.Matrix
	SetFilterEnvelopeValue(envVal uint8)
}
