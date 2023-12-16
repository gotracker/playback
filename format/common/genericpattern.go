package common

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
)

type Pattern[TChannelData song.ChannelData] []Row[TChannelData]

func (p Pattern[TChannelData]) GetRow(row index.Row) song.Row {
	return p[row]
}

func (p Pattern[TChannelData]) NumRows() int {
	return len(p)
}

func (p Pattern[TChannelData]) GetRows() song.Rows {
	return p
}
