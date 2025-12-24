package layout

import (
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/it/channel"
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
)

type SemitoneSample uint16

func (s SemitoneSample) Split() (int, note.Semitone) {
	return int(s >> 8), note.Semitone(s & 0xFF)
}

func NewSemitoneSample(sampIdx int, remap note.Semitone) SemitoneSample {
	return SemitoneSample(uint16(sampIdx<<8) | uint16(remap))
}

type SemitoneSamples [120]SemitoneSample // semitone -> sample + semitone remap

// Song is the full definition of the song data of an IT file
type Song[TPeriod period.Period] struct {
	common.BaseSong[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]

	InstrumentNoteMap map[uint8]SemitoneSamples
	ChannelSettings   []ChannelSetting
	FilterPlugins     map[int]filter.Info
}

// GetNumChannels returns the number of channels the song has
func (s Song[TPeriod]) GetNumChannels() int {
	return len(s.ChannelSettings)
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song[TPeriod]) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s Song[TPeriod]) GetInstrument(instID int, st note.Semitone) (instrument.InstrumentIntf, note.Semitone) {
	if instID == 0 {
		return nil, st
	}

	idx := instID - 1

	if inm, ok := s.InstrumentNoteMap[uint8(instID)]; ok {
		if rm := inm[st]; rm != 0 {
			idx, st = rm.Split()
		}
	}

	if idx < 0 || idx >= len(s.Instruments) {
		return nil, st
	}

	return s.Instruments[idx], st
}

func (s Song[TPeriod]) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	vm := render.NewRowViewModel[channel.Data[TPeriod]](channels)
	rowData := vm.Channels
	song.ForEachRowChannel(row, func(ch index.Channel, d song.ChannelData[itVolume.Volume]) (bool, error) {
		if int(ch) >= channels || !s.ChannelSettings[ch].Enabled || s.ChannelSettings[ch].Muted {
			return true, nil
		}
		rowData[ch] = d.(channel.Data[TPeriod])
		return true, nil
	})
	vm.Channels = rowData
	return render.FormatRowText(vm, longFormat)
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
