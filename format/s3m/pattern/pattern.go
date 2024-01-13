package pattern

import (
	"github.com/gotracker/playback/format/s3m/channel"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/song"
)

type Pattern = song.Pattern[channel.Data, s3mVolume.Volume]
