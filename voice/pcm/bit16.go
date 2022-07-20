package pcm

import (
	"io"

	"github.com/gotracker/gomixing/volume"
)

const (
	cSample16BitVolumeCoeff = volume.Volume(1) / 0x8000
	cSample16BitBytes       = 2
)

// Sample16BitSigned is a signed 16-bit sample
type Sample16BitSigned struct{}

// Volume returns the volume value for the sample
func (s Sample16BitSigned) volume(v int16) volume.Volume {
	return volume.Volume(v) * cSample16BitVolumeCoeff
}

// Size returns the size of the sample in bytes
func (s Sample16BitSigned) Size() int {
	return cSample16BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample16BitSigned) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs)+(cSample16BitBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := int16(d.byteOrder.Uint16(d.data[ofs:]))
	return s.volume(v), nil
}

// Sample16BitUnsigned is an unsigned 16-bit sample
type Sample16BitUnsigned struct{}

// Volume returns the volume value for the sample
func (s Sample16BitUnsigned) volume(v uint16) volume.Volume {
	return volume.Volume(int16(v-0x8000)) * cSample16BitVolumeCoeff
}

// Size returns the size of the sample in bytes
func (s Sample16BitUnsigned) Size() int {
	return cSample16BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample16BitUnsigned) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs)+(cSample16BitBytes-1) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := uint16(d.byteOrder.Uint16(d.data[ofs:]))
	return s.volume(v), nil
}
