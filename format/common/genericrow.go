package common

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
)

type Row[TChannelData song.ChannelData] []TChannelData

func (r Row[TChannelData]) GetNumChannels() int {
	return len(r)
}

func (r Row[TChannelData]) GetChannel(ch index.Channel) song.ChannelData {
	return r[ch]
}
