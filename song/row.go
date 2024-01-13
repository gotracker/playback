package song

import "github.com/gotracker/playback/index"

type RowIntf interface {
	GetNumChannels() int
	GetChannelIntf(index.Channel) ChannelDataIntf
}

// Row is a structure containing a single row
type Row[TChannelData ChannelData[TVolume], TVolume Volume] []TChannelData

func (r Row[TChannelData, TVolume]) GetNumChannels() int {
	return len(r)
}

func (r Row[TChannelData, TVolume]) GetChannel(ch index.Channel) TChannelData {
	return r[ch]
}

func (r Row[TChannelData, TVolume]) GetChannelIntf(ch index.Channel) ChannelDataIntf {
	return r[ch]
}

type RowStringer interface {
	String(options ...any) string
}
