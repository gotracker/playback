package output

import (
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/volume"
)

// PremixData is a structure containing the audio pre-mix data for a specific row or buffer
type PremixData struct {
	SamplesLen  int
	Data        []mixing.ChannelData
	MixerVolume volume.Volume
	Userdata    interface{}
}
