package pcm

import (
	"io"
	"math"

	"github.com/gotracker/gomixing/volume"
)

const (
	//cSample64BitFloatVolumeCoeff = volume.Volume(1)
	cSample64BitFloatBytes = 8
)

// Sample64BitFloat is a 64-bit floating-point sample
type Sample64BitFloat struct{}

// Size returns the size of the sample in bytes
func (s Sample64BitFloat) Size() int {
	return cSample64BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample64BitFloat) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs)+(cSample64BitFloatBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	f := math.Float64frombits(d.byteOrder.Uint64(d.data[ofs:]))
	return volume.Volume(f), nil
}
