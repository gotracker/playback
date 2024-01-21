package layout

import (
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
)

// Header is a mildly-decoded IT header definition
type Header struct {
	Name             string
	InitialSpeed     int
	InitialTempo     int
	GlobalVolume     itVolume.FineVolume
	MixingVolume     itVolume.FineVolume
	LinearFreqSlides bool
	InitialOrder     index.Order
}
