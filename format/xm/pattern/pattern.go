package pattern

import (
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/song"
)

type Pattern = song.Pattern[channel.Data]
