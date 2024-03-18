package sampling

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/gotracker/playback/mixing/volume"
)

const (
	//cSample32BitFloatVolumeCoeff = volume.Volume(1)
	cSample32BitFloatBytes = 4
)

// Sample32BitFloat is a 32-bit floating-point sample
type Sample32BitFloat struct {
	byteOrder binary.ByteOrder
}

// Size returns the size of the sample in bytes
func (Sample32BitFloat) Size() int {
	return cSample32BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample32BitFloat) ReadAt(data []byte, ofs int64) (volume.Volume, error) {
	if len(data) <= int(ofs)+(cSample32BitFloatBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := math.Float32frombits(s.byteOrder.Uint32(data[ofs:]))
	return volume.Volume(v), nil
}

// WriteAt writes a value to the slice provided in the byte order provided
func (s Sample32BitFloat) WriteAt(data []byte, ofs int64, v volume.Volume) error {
	if len(data) <= int(ofs) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	s.byteOrder.PutUint32(data[ofs:], math.Float32bits(float32(v.WithOverflowProtection())))
	return nil
}

// Write writes a value to the Writer provided in the byte order provided
func (s Sample32BitFloat) Write(out io.Writer, v volume.Volume) error {
	return binary.Write(out, s.byteOrder, math.Float32bits(float32(v.WithOverflowProtection())))
}
