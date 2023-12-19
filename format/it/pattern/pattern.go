package pattern

import (
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/song"
)

type Pattern = song.Pattern[channel.Data]
