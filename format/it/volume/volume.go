package volume

import (
	"math"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/types"
)

var (
	MaxItVolume     = itfile.Volume(0x40)
	MaxItFineVolume = FineVolume(0x7f)
	DefaultItVolume = itfile.DefaultVolume

	// DefaultVolume is the default volume value for most everything in IT format
	DefaultVolume = FromItVolume(DefaultItVolume)

	// DefaultMixingVolume is the default mixing volume
	DefaultMixingVolume = itfile.FineVolume(0x30).Value()
)

type Volume itfile.Volume

var (
	_ types.VolumeMaxer[Volume]   = Volume(0)
	_ types.VolumeDeltaer[Volume] = Volume(0)
)

func (v Volume) ToVolume() volume.Volume {
	return volume.Volume(itfile.Volume(v).Value())
}

func (v Volume) IsInvalid() bool {
	return v > 64 && v != 0xff
}

func (v Volume) IsUseInstrumentVol() bool {
	return v == 0xff
}

func (Volume) GetMax() Volume {
	return Volume(MaxItVolume)
}

func (v Volume) FMA(multiplier, add float32) Volume {
	if v == Volume(0xff) {
		return v
	}

	return Volume(min(max(math.FMA(float64(v), float64(multiplier), float64(add)), 0), float64(MaxItVolume)))
}

func (v Volume) AddDelta(d types.VolumeDelta) Volume {
	return Volume(min(max(int16(v)+int16(d), 0), int16(MaxItVolume)))
}

// FromItVolume converts an it volume to a player volume
func FromItVolume(vol itfile.Volume) volume.Volume {
	return volume.Volume(vol.Value())
}

// FromVolPan converts an it volume-pan to a player volume
func FromVolPan(vp uint8) volume.Volume {
	switch {
	case vp <= uint8(MaxItVolume):
		return volume.Volume(vp) / volume.Volume(MaxItVolume)
	default:
		return volume.VolumeUseInstVol
	}
}

// ToItVolume converts a player volume to an it volume
func ToItVolume(v volume.Volume) Volume {
	switch {
	case v == volume.VolumeUseInstVol:
		return Volume(0xff)
	case v < 0.0:
		return 0
	case v > 1.0:
		return Volume(MaxItVolume)
	default:
		return Volume(v * volume.Volume(MaxItVolume))
	}
}
