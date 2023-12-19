package pattern

import (
	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/song"
)

type Pattern = song.Pattern[channel.Data]
