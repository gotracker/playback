package panning

import "github.com/gotracker/playback/mixing/volume"

type PanMixer interface {
	ApplyToMatrix(mtx volume.Matrix) volume.Matrix
	Apply(vol volume.Volume) volume.Matrix
}
