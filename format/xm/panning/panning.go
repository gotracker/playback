package panning

import (
	"math"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/voice/types"
)

const (
	DefaultPanningLeft  = Panning(0x30)
	DefaultPanning      = Panning(0x80)
	DefaultPanningRight = Panning(0xC0)

	MaxPanning = Panning(0xFF)
)

var (
	// DefaultPanningLeftPosition is the default panning value for left channels
	DefaultPanningLeftPosition = PanningFromXm(DefaultPanningLeft)
	// DefaultPanningPosition is the default panning value for unconfigured channels
	DefaultPanningPosition = PanningFromXm(DefaultPanning)
	// DefaultPanningRightPosition is the default panning value for right channels
	DefaultPanningRightPosition = PanningFromXm(DefaultPanningRight)
)

type Panning uint8

var (
	_ types.PanningInformationer[Panning] = Panning(0)
	_ types.PanningDeltaer[Panning]       = Panning(0)
)

func (p Panning) IsInvalid() bool {
	return false
}

func (p Panning) ToPosition() panning.Position {
	return panning.MakeStereoPosition(float32(p), 0, 0xFF)
}

func (Panning) GetDefault() Panning {
	return DefaultPanning
}

func (Panning) GetMax() Panning {
	return MaxPanning
}

func (p Panning) FMA(multiplier, add float32) Panning {
	return Panning(min(max(math.FMA(float64(p), float64(multiplier), float64(add)), 0), 0xFF))
}

func (p Panning) AddDelta(d types.PanDelta) Panning {
	return Panning(min(max(int16(p)+int16(d), 0), int16(MaxPanning)))
}

// PanningFromXm returns a radian panning position from an xm panning value
func PanningFromXm(pos Panning) panning.Position {
	return pos.ToPosition()
}

// PanningToXm returns the xm panning value for a radian panning position
func PanningToXm(pan panning.Position) uint8 {
	val := math.Round(float64(panning.FromStereoPosition(pan, 0, 0xFF)))
	switch {
	case val < 0:
		val = 0
	case val > 0xFF:
		val = 0xFF
	}
	return uint8(val)
}
