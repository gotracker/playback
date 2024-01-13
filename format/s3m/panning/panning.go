package panning

import (
	"math"

	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/playback/voice/types"
)

var (
	// DefaultPanningLeft is the default panning value for left channels
	DefaultPanningLeft = Panning(0x03)
	// DefaultPanning is the default panning value for unconfigured channels
	DefaultPanning = Panning(0x08)
	// DefaultPanningRight is the default panning value for right channels
	DefaultPanningRight = Panning(0x0C)

	MaxPanning = Panning(0x0F)
)

type Panning uint8

var (
	_ types.PanningInformationer[Panning] = Panning(0)
	_ types.PanningDeltaer[Panning]       = Panning(0)
)

func (p Panning) IsInvalid() bool {
	return p > 0x0F
}

func (p Panning) ToPosition() panning.Position {
	return panning.MakeStereoPosition(float32(p), 0, 0x0F)
}

func (Panning) GetDefault() Panning {
	return DefaultPanning
}

func (Panning) GetMax() Panning {
	return MaxPanning
}

func (p Panning) FMA(multiplier, add float32) Panning {
	return Panning(min(max(math.FMA(float64(p), float64(multiplier), float64(add)), 0), 0x0F))
}

func (p Panning) AddDelta(d types.PanDelta) Panning {
	return Panning(min(max(int16(p)+int16(d), 0), int16(MaxPanning)))
}

// PanningFromS3M returns a radian panning position from an S3M panning value
func PanningFromS3M(pos uint8) panning.Position {
	return Panning(pos).ToPosition()
}
