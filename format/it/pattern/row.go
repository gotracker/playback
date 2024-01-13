package pattern

import (
	"github.com/gotracker/playback/format/it/channel"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type Row[TPeriod period.Period] struct {
	song.Row[channel.Data[TPeriod], itVolume.Volume]
}
