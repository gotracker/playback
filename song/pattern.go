package song

import (
	"errors"

	"github.com/gotracker/playback/index"
)

var (
	// ErrStopSong is a magic error asking to stop the current song
	ErrStopSong = errors.New("stop song")
)

// Pattern is structure containing the pattern data
type Pattern[TChannelData ChannelData] []Row[TChannelData]

// GetRow returns a single row of channel data
func (p Pattern[TChannelData]) GetRow(row index.Row) Row[TChannelData] {
	return p[row]
}

// NumRows returns the number of rows contained within the pattern
func (p Pattern[TChannelData]) NumRows() int {
	return len(p)
}

// Patterns is an array of pattern interfaces
type Patterns[TChannelData ChannelData] []Pattern[TChannelData]
