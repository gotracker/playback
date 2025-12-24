package panning

import "github.com/gotracker/playback/mixing/volume"

type MixerStereo struct {
	volume.Matrix
	StereoSeparationFunc func(volume.Matrix) volume.Matrix
}

func (m MixerStereo) ApplyToMatrix(mtx volume.Matrix) volume.Matrix {
	dry := m.Matrix.ApplyToMatrix(mtx)
	return m.StereoSeparationFunc(dry)
}

func (m MixerStereo) Apply(vol volume.Volume) volume.Matrix {
	dry := m.Matrix.Apply(vol)
	return m.StereoSeparationFunc(dry)
}
