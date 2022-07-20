package pcm

import (
	"io"

	"github.com/gotracker/gomixing/volume"
)

const (
	cSample8BitVolumeCoeff = volume.Volume(1) / 0x80
	cSample8BitBytes       = 1
)

// Sample8BitSigned is a signed 8-bit sample
type Sample8BitSigned struct{}

// Volume returns the volume value for the sample
func (s Sample8BitSigned) volume(v int8) volume.Volume {
	return volume.Volume(v) * cSample8BitVolumeCoeff
}

// Size returns the size of the sample in bytes
func (s Sample8BitSigned) Size() int {
	return cSample8BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample8BitSigned) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := int8(d.data[ofs])
	return s.volume(v), nil
}

// Sample8BitUnsigned is an unsigned 8-bit sample
type Sample8BitUnsigned struct{}

// Volume returns the volume value for the sample
func (s Sample8BitUnsigned) volume(v uint8) volume.Volume {
	return volume.Volume(int8(v-0x80)) * cSample8BitVolumeCoeff
}

// Size returns the size of the sample in bytes
func (s Sample8BitUnsigned) Size() int {
	return cSample8BitBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s Sample8BitUnsigned) ReadAt(d *SampleData, ofs int64) (volume.Volume, error) {
	if len(d.data) <= int(ofs) {
		return 0, io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	v := uint8(d.data[ofs])
	return s.volume(v), nil
}
