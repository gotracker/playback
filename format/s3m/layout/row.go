package layout

import (
	"github.com/gotracker/playback/format/s3m/channel"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
)

type Row []channel.Data

func (r Row) Len() int {
	return len(r)
}

func (r Row) ForEach(fn func(ch index.Channel, cd song.ChannelData[s3mVolume.Volume]) (bool, error)) error {
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
