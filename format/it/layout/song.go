package layout

import (
	"reflect"
	"time"

	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itPeriod "github.com/gotracker/playback/format/it/period"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
)

// Song is the full definition of the song data of an IT file
type Song[TPeriod period.Period] struct {
	System            system.System
	Head              Header
	Instruments       map[uint8]*instrument.Instrument[itVolume.FineVolume, itVolume.Volume, itPanning.Panning]
	InstrumentNoteMap map[uint8]map[note.Semitone]NoteInstrument
	Patterns          []song.Pattern
	ChannelSettings   []ChannelSetting
	OrderList         []index.Pattern
	FilterPlugins     map[int]filter.Factory
}

func (s Song[TPeriod]) GetPeriodType() reflect.Type {
	if s.Head.LinearFreqSlides {
		var p period.Linear
		return reflect.TypeOf(p)
	} else {
		var p period.Amiga
		return reflect.TypeOf(p)
	}
}

func (Song[TPeriod]) GetGlobalVolumeType() reflect.Type {
	var gv itVolume.FineVolume
	return reflect.TypeOf(gv)
}

func (Song[TPeriod]) GetChannelMixingVolumeType() reflect.Type {
	var cmv itVolume.FineVolume
	return reflect.TypeOf(cmv)
}

func (Song[TPeriod]) GetChannelVolumeType() reflect.Type {
	var cv itVolume.Volume
	return reflect.TypeOf(cv)
}

func (Song[TPeriod]) GetChannelPanningType() reflect.Type {
	var cp itPanning.Panning
	return reflect.TypeOf(cp)
}

// GetOrderList returns the list of all pattern orders for the song
func (s Song[TPeriod]) GetOrderList() []index.Pattern {
	return s.OrderList
}

// GetInitialBPM returns the initial "tempo" (number of beats per minute) for the song
func (s Song[TPeriod]) GetInitialBPM() int {
	return s.Head.InitialTempo
}

// GetInitialTempo returns the initial "speed" (number of ticks per row) for the song
func (s Song[TPeriod]) GetInitialTempo() int {
	return s.Head.InitialSpeed
}

// GetGlobalVolumeGeneric returns the initial global volume for the song
func (s Song[TPeriod]) GetGlobalVolumeGeneric() volume.Volume {
	return s.Head.GlobalVolume.ToVolume()
}

// GetGlobalVolume returns the initial global volume for the song
func (s Song[TPeriod]) GetGlobalVolume() itVolume.FineVolume {
	return s.Head.GlobalVolume
}

// GetMixingVolumeGeneric returns the initial mixing volume for the song
func (s Song[TPeriod]) GetMixingVolumeGeneric() volume.Volume {
	return s.Head.MixingVolume.ToVolume()
}

// GetMixingVolume returns the initial mixing volume for the song
func (s Song[TPeriod]) GetMixingVolume() itVolume.FineVolume {
	return s.Head.MixingVolume
}

const durationPerBpm = time.Duration(2500) * time.Millisecond

// GetTickDuration calculates the duration of a tick at a particular BPM
func (s Song[TPeriod]) GetTickDuration(bpm int) time.Duration {
	if bpm == 0 {
		return 0
	}

	return durationPerBpm / time.Duration(bpm)
}

// GetPattern returns a specific pattern indexed by `patNum`
func (s Song[TPeriod]) GetPattern(patNum index.Pattern) (song.Pattern, error) {
	if int(patNum) >= len(s.Patterns) {
		return nil, song.ErrStopSong
	}
	return s.Patterns[patNum], nil
}

// GetPatternIntfByOrder returns the pattern specified by the order index provided
func (s Song[TPeriod]) GetPatternByOrder(o index.Order) (song.Pattern, error) {
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
func (s Song[TPeriod]) GetNumChannels() int {
	return len(s.ChannelSettings)
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song[TPeriod]) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

// NumInstruments returns the number of instruments in the song
func (s Song[TPeriod]) NumInstruments() int {
	return len(s.Instruments)
}

// IsValidInstrumentID returns true if the instrument exists
func (s Song[TPeriod]) IsValidInstrumentID(instNum instrument.ID) bool {
	if instNum.IsEmpty() {
		return false
	}
	switch id := instNum.(type) {
	case channel.SampleID:
		_, ok := s.Instruments[id.InstID]
		return ok
	}
	return false
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s Song[TPeriod]) GetInstrument(instNum instrument.ID) (instrument.InstrumentIntf, note.Semitone) {
	if instNum.IsEmpty() {
		return nil, note.UnchangedSemitone
	}
	switch id := instNum.(type) {
	case channel.SampleID:
		if nm, ok1 := s.InstrumentNoteMap[id.InstID]; ok1 {
			if sm, ok2 := nm[id.Semitone]; ok2 {
				return sm.Inst, sm.NoteRemap
			}
		}
	}
	return nil, note.UnchangedSemitone
}

// GetName returns the name of the song
func (s Song[TPeriod]) GetName() string {
	return s.Head.Name
}

func (s Song[TPeriod]) GetPeriodCalculator() song.PeriodCalculatorIntf {
	if s.Head.LinearFreqSlides {
		return itPeriod.LinearConverter
	}
	return itPeriod.AmigaConverter
}

func (s Song[TPeriod]) GetInitialOrder() index.Order {
	return s.Head.InitialOrder
}

func (s Song[TPeriod]) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	rt := render.NewRowText[channel.Data[TPeriod]](channels, longFormat)
	rowData := make([]channel.Data[TPeriod], channels)
	copy(rowData, row.(Row[TPeriod]))
	rt.Channels = rowData
	return rt
}

func (s Song[TPeriod]) GetSystem() system.System {
	return s.System
}

func (s Song[TPeriod]) ForEachChannel(enabledOnly bool, fn func(ch index.Channel) (bool, error)) error {
	for i, cs := range s.ChannelSettings {
		if enabledOnly && !cs.Enabled {
			continue
		}
		cont, err := fn(index.Channel(i))
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}

func (Song[TPeriod]) IsOPL2Enabled() bool {
	return false
}
