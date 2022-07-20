package pcm

import (
	"io"
	"math"

	"github.com/gotracker/gomixing/volume"
)

const (
	//cSample32BitFloatVolumeCoeff = volume.Volume(1)
	cSample32BitFloatBytes = 4
)

// Sample32BitFloat is a 32-bit floating-point sample
type Sample32BitFloat struct{}

// Size returns the size of the sample in bytes
func (s Sample32BitFloat) Size() int {
	return cSample32BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample32BitFloat) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs)+(cSample32BitFloatBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := math.Float32frombits(d.byteOrder.Uint32(d.data[ofs:]))
	return volume.Volume(v), nil
}
