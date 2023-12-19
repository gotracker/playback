package song

import "github.com/gotracker/playback/index"

// Row is a structure containing a single row
type Row[TChannelData ChannelData] []TChannelData

func (r Row[TChannelData]) GetNumChannels() int {
	return len(r)
}

func (r Row[TChannelData]) GetChannel(ch index.Channel) TChannelData {
	return r[ch]
}
