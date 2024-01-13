package render

import (
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/song"
	channelfilter "github.com/gotracker/playback/voice/filter"
)

type ChannelIntf interface {
	channelfilter.Applier
	GetPremixVolume() volume.Volume
}

// Channel is the important bits to make output to a particular downmixing channel work
type Channel[TGlobalVolume, TMixingVolume song.Volume, TPanning song.Panning] struct {
	ChannelNum   int
	Filter       filter.Filter
	GetOPL2Chip  func() OPL2Chip
	GlobalVolume volume.Volume // this is the channel's version of the GlobalVolume
}

// ApplyFilter will apply the channel filter, if there is one.
func (oc *Channel[TGlobalVolume, TMixingVolume, TPanning]) ApplyFilter(dry volume.Matrix) volume.Matrix {
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
func (oc *Channel[TGlobalVolume, TMixingVolume, TPanning]) GetPremixVolume() volume.Volume {
	return oc.GlobalVolume
}

// SetFilterEnvelopeValue updates the filter on the channel with the new envelope value
func (oc *Channel[TGlobalVolume, TMixingVolume, TPanning]) SetFilterEnvelopeValue(envVal uint8) {
	if oc.Filter != nil {
		oc.Filter.UpdateEnv(envVal)
	}
}
