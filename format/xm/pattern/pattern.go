package pattern

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
)

type Pattern []Row

func (p Pattern) GetRow(row index.Row) song.Row {
	return p[row]
}

func (p Pattern) NumRows() int {
	return len(p)
}

func (p Pattern) GetRows() song.Rows {
	return p
}
