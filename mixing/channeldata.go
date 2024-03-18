package mixing

import (
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

// Data is a single buffer of data at a specific panning position
type Data struct {
	Data       MixBuffer
	PanMatrix  panning.PanMixer
	Volume     volume.Volume
	Pos        int
	SamplesLen int
	Flush      func()
}

// ChannelData is a single channel's data
type ChannelData []Data
