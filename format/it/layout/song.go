package layout

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/song"
)

// Song is the full definition of the song data of an Song file
type Song struct {
	Head              Header
	Instruments       map[uint8]*instrument.Instrument
	InstrumentNoteMap map[uint8]map[note.Semitone]NoteInstrument
	Patterns          []pattern.Pattern[channel.Data]
	ChannelSettings   []ChannelSetting
	OrderList         []index.Pattern
	FilterPlugins     map[int]filter.Factory
}

// GetOrderList returns the list of all pattern orders for the song
func (s Song) GetOrderList() []index.Pattern {
	return s.OrderList
}

// GetPattern returns an interface to a specific pattern indexed by `patNum`
func (s Song) GetPattern(patNum index.Pattern) song.Pattern[channel.Data] {
	if int(patNum) >= len(s.Patterns) {
		return nil
	}
	return &s.Patterns[patNum]
}

// IsChannelEnabled returns true if the channel at index `channelNum` is enabled
func (s Song) IsChannelEnabled(channelNum int) bool {
	return s.ChannelSettings[channelNum].Enabled
}

// GetRenderChannel returns the output channel for the channel at index `channelNum`
func (s Song) GetRenderChannel(channelNum int) int {
	return s.ChannelSettings[channelNum].OutputChannelNum
}

// NumInstruments returns the number of instruments in the song
func (s Song) NumInstruments() int {
	return len(s.Instruments)
}

// IsValidInstrumentID returns true if the instrument exists
func (s Song) IsValidInstrumentID(instNum instrument.ID) bool {
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
func (s Song) GetInstrument(instNum instrument.ID) (*instrument.Instrument, note.Semitone) {
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
func (s Song) GetName() string {
	return s.Head.Name
}
