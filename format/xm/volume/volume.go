package volume

import (
	"math"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/types"
)

const (
	DefaultXmVolume       = XmVolume(0x40)
	DefaultXmMixingVolume = XmVolume(0x18)
)

var (
	// DefaultVolume is the default volume value for most everything in xm format
	DefaultVolume = ToVolume(0x10 + VolEffect(DefaultXmVolume))

	// DefaultMixingVolume is the default mixing volume
	DefaultMixingVolume = volume.Volume(0x30) / 0x80
)

// XmVolume is a helpful converter from the XM range of 0-64 into a volume
type XmVolume uint8

var (
	_ types.VolumeMaxer[XmVolume]   = XmVolume(0)
	_ types.VolumeDeltaer[XmVolume] = XmVolume(0)
)

const cVolumeXMCoeff = volume.Volume(1) / 0x40

// Volume returns the volume from the internal format
func (v XmVolume) ToVolume() volume.Volume {
	if v != 0xff {
		return volume.Volume(v) * cVolumeXMCoeff
	}
	return volume.VolumeUseInstVol
}

func (v XmVolume) IsInvalid() bool {
	return v > 0x40 && v != 0xff
}

func (v XmVolume) IsUseInstrumentVol() bool {
	return v == 0xff
}

func (XmVolume) GetMax() XmVolume {
	return 0x40
}

func (v XmVolume) FMA(multiplier, add float32) XmVolume {
	if v == XmVolume(0xff) {
		return v
	}

	return XmVolume(min(max(math.FMA(float64(v), float64(multiplier), float64(add)), 0), 64))
}

func (v XmVolume) AddDelta(d types.VolumeDelta) XmVolume {
	return XmVolume(min(max(int16(v)+int16(d), 0), 0x40))
}

// ToVolumeXM returns the VolumeXM representation of a volume
func ToVolumeXM(v volume.Volume) XmVolume {
	if v != volume.VolumeUseInstVol {
		return XmVolume(v * 0x40)
	}
	return XmVolume(0xff)
}
