package layout

import (
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
)

// Header is a mildly-decoded XM header definition
type Header struct {
	Name             string
	InitialSpeed     int
	InitialTempo     int
	GlobalVolume     xmVolume.XmVolume
	MixingVolume     xmVolume.XmVolume
	LinearFreqSlides bool
	InitialOrder     index.Order
}
