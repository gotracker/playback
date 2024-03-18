package sampling

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/gotracker/playback/mixing/volume"
)

const (
	//cSample64BitFloatVolumeCoeff = volume.Volume(1)
	cSample64BitFloatBytes = 8
)

// Sample64BitFloat is a 64-bit floating-point sample
type Sample64BitFloat struct {
	byteOrder binary.ByteOrder
}

// Size returns the size of the sample in bytes
func (Sample64BitFloat) Size() int {
	return cSample64BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample64BitFloat) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs)+(cSample64BitFloatBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	f := math.Float64frombits(s.byteOrder.Uint64(data[ofs:]))
	return volume.Volume(f), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample64BitFloat) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	s.byteOrder.PutUint64(data[ofs:], math.Float64bits(v.WithOverflowProtection()))
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample64BitFloat) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, s.byteOrder, math.Float64bits(v.WithOverflowProtection()))
}
