package layout

import "github.com/gotracker/gomixing/volume"

// Header is a mildly-decoded S3M header definition
type Header struct {
	Name         string
	InitialSpeed int
	InitialTempo int
	GlobalVolume volume.Volume
	MixingVolume volume.Volume
	Stereo       bool
}
