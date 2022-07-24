package layout

import (
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/song"
)

// Layout is the full definition of the song data of an XM file
type Layout struct {
	Head            Header
	Instruments     map[uint8]*instrument.Keyboard[*instrument.Instrument]
	Patterns        []pattern.Pattern[channel.Data]
	ChannelSettings []ChannelSetting
	OrderList       []index.Pattern
}

// GetOrderList returns the list of all pattern orders for the song
func (s Layout) GetOrderList() []index.Pattern {
	return s.OrderList
}

// GetPattern returns an interface to a specific pattern indexed by `patNum`
func (s Layout) GetPattern(patNum index.Pattern) song.Pattern[channel.Data] {
	if int(patNum) >= len(s.Patterns) {
		return nil
	}
	return &s.Patterns[patNum]
}

// IsChannelEnabled returns true if the channel at index `channelNum` is enabled
func (s Layout) IsChannelEnabled(channelNum int) bool {
	return s.ChannelSettings[channelNum].Enabled
}

// GetRenderChannel returns the output channel for the channel at index `channelNum`
func (s Layout) GetRenderChannel(channelNum int) int {
	return s.ChannelSettings[channelNum].OutputChannelNum
}

// NumInstruments returns the number of instruments in the song
func (s Layout) NumInstruments() int {
	return len(s.Instruments)
}

// IsValidInstrumentID returns true if the instrument exists
func (s Layout) IsValidInstrumentID(instNum instrument.ID) bool {
	if instNum.IsEmpty() {
		return false
	}
	switch id := instNum.(type) {
	case channel.SampleID:
		keyboard, ok := s.Instruments[id.InstID]
		return ok && keyboard.GetInstrument() != nil
	}
	return false
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s Layout) GetInstrument(instNum instrument.ID) (*instrument.Instrument, note.Semitone) {
	if instNum.IsEmpty() {
		return nil, note.UnchangedSemitone
	}
	switch id := instNum.(type) {
	case channel.SampleID:
		keyboard, ok := s.Instruments[id.InstID]
		if !ok {
			return nil, note.UnchangedSemitone
		}

		if remapInst, ok := keyboard.GetRemap(id.Semitone); ok {
			return remapInst, id.Semitone
		}

		inst := keyboard.GetInstrument()
		return inst, id.Semitone
	}
	return nil, note.UnchangedSemitone
}

// GetName returns the name of the song
func (s Layout) GetName() string {
	return s.Head.Name
}
