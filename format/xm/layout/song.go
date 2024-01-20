package layout

import (
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/xm/channel"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
)

// Song is the full definition of the song data of an XM file
type Song[TPeriod period.Period] struct {
	common.BaseSong[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]

	InstrumentNoteMap map[uint8]SemitoneSamples
	ChannelSettings   []ChannelSetting
}

type SemitoneSamples [96]int // semitone -> instrument index

// GetNumChannels returns the number of channels the song has
func (s Song[TPeriod]) GetNumChannels() int {
	return len(s.ChannelSettings)
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song[TPeriod]) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s Song[TPeriod]) GetInstrument(instID instrument.ID) (instrument.InstrumentIntf, note.Semitone) {
	if instID.IsEmpty() {
		return nil, note.UnchangedSemitone
	}

	idx, st := instID.GetIndexAndSemitone()

	i := idx
	if inm, ok := s.InstrumentNoteMap[uint8(idx)]; ok {
		i = inm[st]
	}

	if i < 0 || i >= len(s.Instruments) {
		return nil, st
	}

	return s.Instruments[i], st
}

func (s Song[TPeriod]) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	rt := render.NewRowText[channel.Data[TPeriod]](channels, longFormat)
	rowData := make([]channel.Data[TPeriod], channels)
	song.ForEachRowChannel(row, func(ch index.Channel, d song.ChannelData[xmVolume.XmVolume]) (bool, error) {
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
