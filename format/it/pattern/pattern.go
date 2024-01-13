package pattern

import (
	"github.com/gotracker/playback/format/it/channel"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type Pattern[TPeriod period.Period] struct {
	song.Pattern[channel.Data[TPeriod], itVolume.Volume]
}
