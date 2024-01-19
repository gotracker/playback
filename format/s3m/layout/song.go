package layout

import (
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
)

type Song struct {
	common.BaseSong[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]

	ChannelSettings []ChannelSetting
	ChannelOrders   []index.Channel
	NumChannels     int
}

// GetNumChannels returns the number of channels the song has
func (s Song) GetNumChannels() int {
	return s.NumChannels
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

func (s Song) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	nch := min(s.NumChannels, channels)
	rt := render.NewRowText[channel.Data](nch, longFormat)
	rowData := make([]channel.Data, 0, nch)
	_ = song.ForEachRowChannel[s3mVolume.Volume](row, func(ch index.Channel, d song.ChannelData[s3mVolume.Volume]) (bool, error) {
		if int(ch) >= nch || !s.ChannelSettings[ch].Enabled || s.ChannelSettings[ch].Muted {
			return true, nil
		}
		rowData = append(rowData, d.(channel.Data))
		return true, nil
	})
	for len(rowData) < nch {
		rowData = append(rowData, channel.Data{})
	}
	rt.Channels = rowData
	return rt
}

func (s Song) ForEachChannel(enabledOnly bool, fn func(ch index.Channel) (bool, error)) error {
	for _, ch := range s.ChannelOrders {
		cs := &s.ChannelSettings[ch]
		if enabledOnly {
			if !cs.Enabled || (cs.Muted && s.MS.Quirks.DoNotProcessEffectsOnMutedChannels) {
				continue
			}
		}
		cont, err := fn(ch)
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}

func (s Song) IsOPL2Enabled() bool {
	for _, cs := range s.ChannelSettings {
		if !cs.Enabled || cs.Muted {
			continue
		}

		if cs.GetOPLChannel().IsValid() {
			return true
		}
	}
	return false
}
