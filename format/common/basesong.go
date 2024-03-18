package common

import (
	"reflect"
	"time"

	"github.com/gotracker/playback/mixing/volume"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
	"github.com/gotracker/playback/voice/types"
)

type BaseSong[TPeriod types.Period, TGlobalVolume, TMixingVolume, TVolume types.Volume, TPanning types.Panning] struct {
	System system.System
	MS     *settings.MachineSettings[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]

	Name         string
	InitialBPM   int
	InitialTempo int
	GlobalVolume TGlobalVolume
	MixingVolume TMixingVolume
	InitialOrder index.Order

	Instruments []*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning]
	Patterns    []song.Pattern
	OrderList   []index.Pattern
}

func (BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetPeriodType() reflect.Type {
	var p TPeriod
	return reflect.TypeOf(p)
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetGlobalVolumeType() reflect.Type {
	return reflect.TypeOf(s.GlobalVolume)
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelMixingVolumeType() reflect.Type {
	return reflect.TypeOf(s.MixingVolume)
}

func (BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelVolumeType() reflect.Type {
	var v TVolume
	return reflect.TypeOf(v)
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelPanningType() reflect.Type {
	var p TPanning
	return reflect.TypeOf(p)
}

// GetOrderList returns the list of all pattern orders for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetOrderList() []index.Pattern {
	return s.OrderList
}

// GetInitialBPM returns the initial "tempo" (number of beats per minute) for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetInitialBPM() int {
	return s.InitialBPM
}

// GetInitialTempo returns the initial "speed" (number of ticks per row) for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetInitialTempo() int {
	return s.InitialTempo
}

// GetGlobalVolumeGeneric returns the initial global volume for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetGlobalVolumeGeneric() volume.Volume {
	return s.GlobalVolume.ToVolume()
}

// GetGlobalVolume returns the initial global volume for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetGlobalVolume() TGlobalVolume {
	return s.GlobalVolume
}

// GetMixingVolumeGeneric returns the initial mixing volume for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetMixingVolumeGeneric() volume.Volume {
	return s.MixingVolume.ToVolume()
}

// GetMixingVolume returns the initial mixing volume for the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetMixingVolume() TMixingVolume {
	return s.MixingVolume
}

const durationPerBpm = time.Duration(2500) * time.Millisecond

// GetTickDuration calculates the duration of a tick at a particular BPM
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetTickDuration(bpm int) time.Duration {
	if bpm == 0 {
		return 0
	}

	return durationPerBpm / time.Duration(bpm)
}

// GetPattern returns a specific pattern indexed by `patNum`
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetPattern(patNum index.Pattern) (song.Pattern, error) {
	if int(patNum) >= len(s.Patterns) {
		return nil, song.ErrStopSong
	}
	return s.Patterns[patNum], nil
}

// GetPatternByOrder returns the pattern specified by the order index provided
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetPatternByOrder(o index.Order) (song.Pattern, error) {
	if int(o) >= len(s.OrderList) {
		return nil, song.ErrStopSong
	}

	pat := s.OrderList[o]
	switch pat {
	case index.InvalidPattern:
		return nil, song.ErrStopSong
	case index.NextPattern:
		return nil, index.ErrNextPattern
	}

	return s.GetPattern(pat)
}

// GetNumChannels returns the number of channels the song has
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetNumChannels() int {
	panic("unimplemented")
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	panic("unimplemented")
}

// NumInstruments returns the number of instruments in the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) NumInstruments() int {
	return len(s.Instruments)
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetInstrument(instID int, st note.Semitone) (instrument.InstrumentIntf, note.Semitone) {
	if instID == 0 {
		return nil, st
	}
	idx := instID - 1
	if idx >= len(s.Instruments) {
		return nil, st
	}
	return s.Instruments[idx], st
}

// GetName returns the name of the song
func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetName() string {
	return s.Name
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetPeriodCalculator() song.PeriodCalculatorIntf {
	return s.MS.PeriodConverter
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetInitialOrder() index.Order {
	return s.InitialOrder
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	panic("unimplemented")
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetSystem() system.System {
	return s.System
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ForEachChannel(enabledOnly bool, fn func(ch index.Channel) (bool, error)) error {
	panic("unimplemented")
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) IsOPL2Enabled() bool {
	return false
}

func (s BaseSong[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetMachineSettings() any {
	return s.MS
}
