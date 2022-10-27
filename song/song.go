package song

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

// Data is an interface to the song data
type Data interface {
	GetOrderList() []index.Pattern
	IsChannelEnabled(int) bool
	GetRenderChannel(int) int
	NumInstruments() int
	IsValidInstrumentID(instrument.ID) bool
	GetInstrument(instrument.ID) (*instrument.Instrument, note.Semitone)
	GetName() string
}

type DataPattern[TChannelData any] interface {
	GetPattern(index.Pattern) Pattern[TChannelData]
}

type DataFlags[TFlags any] interface {
	GetFlags() *TFlags
}
