package layout

import (
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
)

// Header is a mildly-decoded S3M header definition
type Header struct {
	Name         string
	InitialSpeed int
	InitialTempo int
	GlobalVolume s3mVolume.Volume
	MixingVolume s3mVolume.FineVolume
	Stereo       bool
	InitialOrder index.Order
}
