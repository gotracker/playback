package output

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/volume"
)

// PremixData is a structure containing the audio pre-mix data for a specific row or buffer
type PremixData struct {
	SamplesLen  int
	Data        []mixing.ChannelData
	MixerVolume volume.Volume
	Userdata    interface{}
}
