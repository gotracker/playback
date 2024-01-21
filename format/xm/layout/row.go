package layout

import (
	"github.com/gotracker/playback/format/xm/channel"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type Row[TPeriod period.Period] []channel.Data[TPeriod]

func (r Row[TPeriod]) Len() int {
	return len(r)
}

func (r Row[TPeriod]) ForEach(fn func(ch index.Channel, cd song.ChannelData[xmVolume.XmVolume]) (bool, error)) error {
	for i, c := range r {
		cont, err := fn(index.Channel(i), c)
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}
