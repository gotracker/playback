package playback

import (
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/player/render"
)

func (m *Manager) getRowText() *render.RowDisplay[channel.Data] {
	nCh := 0
	for ch := range m.channels {
		if !m.song.IsChannelEnabled(ch) {
			continue
		}
		nCh++
	}
	rowText := render.NewRowText[channel.Data](nCh, true)
	for ch, cs := range m.channels {
		if !m.song.IsChannelEnabled(ch) {
			continue
		}

		if cd := cs.GetData(); cd != nil {
			rowText.Channels[ch] = *cd
		}
	}
	return &rowText
}
