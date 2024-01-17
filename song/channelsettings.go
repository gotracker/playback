package song

import (
	"errors"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/voice/vol0optimization"
)

type ChannelSettings interface {
	IsMuted() bool
	GetOutputChannelNum() int
	GetMemory() ChannelMemory
	IsPanEnabled() bool
	GetDefaultFilterName() string
	IsDefaultFilterEnabled() bool
	GetVol0OptimizationSettings() vol0optimization.Vol0OptimizationSettings
	GetOPLChannel() index.OPLChannel
}

type channelInitialVolumeGetter[TVolume Volume] interface {
	GetInitialVolume() TVolume
}

func GetChannelInitialVolume[TVolume Volume](c ChannelSettings) (TVolume, error) {
	gicv, ok := c.(channelInitialVolumeGetter[TVolume])
	if !ok {
		var empty TVolume
		return empty, errors.New("could not identify channel initial volume interface")
	}

	return gicv.GetInitialVolume(), nil
}

func GetChannelMixingVolume[TMixingVolume Volume](c ChannelSettings) (TMixingVolume, error) {
	gcmv, ok := c.(mixingVolumeGetter[TMixingVolume])
	if !ok {
		var empty TMixingVolume
		return empty, errors.New("could not identify channel volume interface")
	}

	return gcmv.GetMixingVolume(), nil
}

type channelInitialPanningGetter[TPanning Panning] interface {
	GetInitialPanning() TPanning
}

func GetChannelInitialPanning[TPanning Panning](c ChannelSettings) (TPanning, error) {
	gicp, ok := c.(channelInitialPanningGetter[TPanning])
	if !ok {
		var empty TPanning
		return empty, errors.New("could not identify channel initial panning interface")
	}

	return gicp.GetInitialPanning(), nil
}
