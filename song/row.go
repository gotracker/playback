package song

import "github.com/gotracker/playback/index"

// Row is an interface to a row
type Row interface {
	GetNumChannels() int
	GetChannel(index.Channel) ChannelData
}
