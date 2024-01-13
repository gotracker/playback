package song

import (
	"errors"

	"github.com/gotracker/playback/index"
)

var (
	// ErrStopSong is a magic error asking to stop the current song
	ErrStopSong = errors.New("stop song")
)

type PatternIntf interface {
	GetRowIntf(row index.Row) RowIntf
	NumRows() int
}

// Pattern is structure containing the pattern data
type Pattern[TChannelData ChannelData[TVolume], TVolume Volume] []Row[TChannelData, TVolume]

// GetRow returns a single row of channel data
func (p Pattern[TChannelData, TVolume]) GetRow(row index.Row) Row[TChannelData, TVolume] {
	return p[row]
}

func (p Pattern[TChannelData, TVolume]) GetRowIntf(row index.Row) RowIntf {
	return p[row]
}

// NumRows returns the number of rows contained within the pattern
func (p Pattern[TChannelData, TVolume]) NumRows() int {
	return len(p)
}

// Patterns is an array of pattern interfaces
type Patterns[TChannelData ChannelData[TVolume], TVolume Volume] []Pattern[TChannelData, TVolume]
