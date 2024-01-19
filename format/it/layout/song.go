package layout

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
)

// Song is the full definition of the song data of an IT file
type Song[TPeriod period.Period] struct {
	common.BaseSong[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]

	InstrumentNoteMap map[uint8]map[note.Semitone]NoteInstrument[TPeriod]
	ChannelSettings   []ChannelSetting
	FilterPlugins     map[int]filter.Factory
}

// GetNumChannels returns the number of channels the song has
func (s Song[TPeriod]) GetNumChannels() int {
	return len(s.ChannelSettings)
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song[TPeriod]) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

func (s Song[TPeriod]) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	rt := render.NewRowText[channel.Data[TPeriod]](channels, longFormat)
	rowData := make([]channel.Data[TPeriod], channels)
	song.ForEachRowChannel(row, func(ch index.Channel, d song.ChannelData[itVolume.Volume]) (bool, error) {
		if int(ch) >= channels || !s.ChannelSettings[ch].Enabled || s.ChannelSettings[ch].Muted {
			return true, nil
		}
		rowData[ch] = d.(channel.Data[TPeriod])
		return true, nil
	})
	rt.Channels = rowData
	return rt
}

func (s Song[TPeriod]) ForEachChannel(enabledOnly bool, fn func(ch index.Channel) (bool, error)) error {
	for i, cs := range s.ChannelSettings {
		if enabledOnly {
			if !cs.Enabled || (cs.Muted && s.MS.Quirks.DoNotProcessEffectsOnMutedChannels) {
				continue
			}
		}
		cont, err := fn(index.Channel(i))
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}
