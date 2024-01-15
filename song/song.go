package song

import (
	"errors"
	"reflect"
	"time"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/system"
	"github.com/gotracker/playback/voice/types"
)

// Data is an interface to the song data
type Data interface {
	GetPeriodType() reflect.Type
	GetGlobalVolumeType() reflect.Type
	GetChannelMixingVolumeType() reflect.Type
	GetChannelVolumeType() reflect.Type
	GetChannelPanningType() reflect.Type

	GetInitialBPM() int
	GetInitialTempo() int
	GetMixingVolumeGeneric() volume.Volume
	GetTickDuration(bpm int) time.Duration
	GetOrderList() []index.Pattern
	GetNumChannels() int
	GetChannelSettings(index.Channel) ChannelSettings
	NumInstruments() int
	IsValidInstrumentID(instrument.ID) bool
	GetInstrument(instrument.ID) (instrument.InstrumentIntf, note.Semitone)
	GetName() string
	GetPatternByOrder(index.Order) (Pattern, error)
	GetPattern(index.Pattern) (Pattern, error)
	GetPeriodCalculator() PeriodCalculatorIntf
	GetInitialOrder() index.Order
	GetRowRenderStringer(Row, int, bool) RowStringer
	GetSystem() system.System
	ForEachChannel(enabledOnly bool, fn func(ch index.Channel) (bool, error)) error
	IsOPL2Enabled() bool
}

type (
	Volume  = types.Volume
	Panning = types.Panning
)

type globalVolumeGetter[TGlobalVolume Volume] interface {
	GetGlobalVolume() TGlobalVolume
}

func GetGlobalVolume[TGlobalVolume Volume](s Data) (TGlobalVolume, error) {
	ggv, ok := s.(globalVolumeGetter[TGlobalVolume])
	if !ok {
		var empty TGlobalVolume
		return empty, errors.New("could not identify global volume interface")
	}

	return ggv.GetGlobalVolume(), nil
}

type mixingVolumeGetter[TMixingVolume Volume] interface {
	GetMixingVolume() TMixingVolume
}

func GetMixingVolume[TMixingVolume Volume](s Data) (TMixingVolume, error) {
	gmv, ok := s.(mixingVolumeGetter[TMixingVolume])
	if !ok {
		var empty TMixingVolume
		return empty, errors.New("could not identify mixing volume interface")
	}

	return gmv.GetMixingVolume(), nil
}
