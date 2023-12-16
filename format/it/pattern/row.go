package pattern

import (
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/song"
)

type Row []channel.Data

func (r Row) GetChannels() []song.ChannelData {
	c := make([]song.ChannelData, len(r))
	for i := range r {
		c[i] = &r[i]
	}
	return c
}
