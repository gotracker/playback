package render

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice/render"

	"github.com/gotracker/gomixing/volume"
)

// Channel is the important bits to make output to a particular downmixing channel work
type Channel struct {
	ChannelNum       int
	Filter           filter.Filter
	GetSampleRate    func() period.Frequency
	SetGlobalVolume  func(volume.Volume)
	GetOPL2Chip      func() render.OPL2Chip
	ChannelVolume    volume.Volume
	LastGlobalVolume volume.Volume // this is the channel's version of the GlobalVolume
}

// ApplyFilter will apply the channel filter, if there is one.
func (oc *Channel) ApplyFilter(dry volume.Matrix) volume.Matrix {
	if dry.Channels == 0 {
		return volume.Matrix{}
	}
	premix := oc.GetPremixVolume()
	wet := dry.Apply(premix)
	if oc.Filter != nil {
		return oc.Filter.Filter(wet)
	}
	return wet
}

// GetPremixVolume returns the premix volume of the output channel
func (oc *Channel) GetPremixVolume() volume.Volume {
	return oc.LastGlobalVolume * oc.ChannelVolume
}

// SetFilterEnvelopeValue updates the filter on the channel with the new envelope value
func (oc *Channel) SetFilterEnvelopeValue(envVal int8) {
	if oc.Filter != nil {
		oc.Filter.UpdateEnv(envVal)
	}
}
