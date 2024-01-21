package layout

import (
	"github.com/gotracker/playback/format/it/channel"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type Row[TPeriod period.Period] []channel.Data[TPeriod]

func (r Row[TPeriod]) Len() int {
	return len(r)
}

func (r Row[TPeriod]) ForEach(fn func(ch index.Channel, cd song.ChannelData[itVolume.Volume]) (bool, error)) error {
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
