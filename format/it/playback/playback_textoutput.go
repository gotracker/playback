package playback

import (
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/player/render"
)

func (m *manager[TPeriod]) getRowText() *render.RowDisplay[channel.Data] {
	nCh := 0
	for ch := range m.channels {
		if !m.song.IsChannelEnabled(ch) {
			continue
		}
		nCh++
	}
	rowText := render.NewRowText[channel.Data](nCh, m.longChannelOutput)
	for ch, cs := range m.channels {
		if !m.song.IsChannelEnabled(ch) {
			continue
		}

		rowText.Channels[ch] = cs.GetChannelData()
	}
	return &rowText
}
