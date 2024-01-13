package pattern

import (
	"github.com/gotracker/playback/format/xm/channel"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type Pattern[TPeriod period.Period] struct {
	song.Pattern[channel.Data[TPeriod], xmVolume.XmVolume]
}
