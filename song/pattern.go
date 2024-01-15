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
	GetRow(row index.Row) Row
	NumRows() int
}

// Pattern is structure containing the pattern data
type Pattern []Row

// GetRow returns a single row of channel data
func (p Pattern) GetRow(row index.Row) Row {
	return p[row]
}

// NumRows returns the number of rows contained within the pattern
func (p Pattern) NumRows() int {
	return len(p)
}

// Patterns is an array of pattern interfaces
type Patterns[TChannelData ChannelData[TVolume], TVolume Volume] []Pattern
