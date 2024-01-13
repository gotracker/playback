package volume

import (
	"math"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/voice/types"
)

type FineVolume itfile.FineVolume

var (
	_ types.VolumeMaxer[FineVolume]   = FineVolume(0)
	_ types.VolumeDeltaer[FineVolume] = FineVolume(0)
)

func (v FineVolume) ToVolume() volume.Volume {
	return volume.Volume(itfile.FineVolume(v).Value())
}

func (v FineVolume) IsInvalid() bool {
	return v > 0x80 && v != 0xFF
}

func (v FineVolume) IsUseInstrumentVol() bool {
	return v == 0xFF
}

func (FineVolume) GetMax() FineVolume {
	return MaxItFineVolume
}

func (v FineVolume) FMA(multiplier, add float32) FineVolume {
	if v == FineVolume(0xff) {
		return v
	}

	return min(FineVolume(max(math.FMA(float64(v), float64(multiplier), float64(add)), 0)), MaxItFineVolume)
}

func (v FineVolume) AddDelta(d types.VolumeDelta) FineVolume {
	return FineVolume(min(max(int16(v)+int16(d), 0), int16(MaxItFineVolume)))
}

// ToItFineVolume converts a player volume to an it fine volume
func ToItFineVolume(v volume.Volume) FineVolume {
	switch {
	case v == volume.VolumeUseInstVol:
		return FineVolume(0xff)
	case v < 0.0:
		return 0
	case v > 1.0:
		return FineVolume(MaxItFineVolume)
	default:
		return FineVolume(v * volume.Volume(MaxItFineVolume))
	}
}
