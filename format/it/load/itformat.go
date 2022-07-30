package load

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	itblock "github.com/gotracker/goaudiofile/music/tracked/it/block"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/format/it/layout"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/pattern"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/voice/pcm"
)

func moduleHeaderToHeader(fh *itfile.ModuleHeader) (*layout.Header, error) {
	if fh == nil {
		return nil, errors.New("file header is nil")
	}
	head := layout.Header{
		Name:         fh.GetName(),
		InitialSpeed: int(fh.InitialSpeed),
		InitialTempo: int(fh.InitialTempo),
		GlobalVolume: volume.Volume(fh.GlobalVolume.Value()),
	}
	switch {
	case fh.TrackerCompatVersion < 0x200:
		head.MixingVolume = volume.Volume(fh.MixingVolume.Value())
	case fh.TrackerCompatVersion >= 0x200:
		head.MixingVolume = volume.Volume(fh.MixingVolume) / 128
	}
	return &head, nil
}

func convertItPattern(pkt itfile.PackedPattern, channels int) (*pattern.Pattern[channel.Data], int, error) {
	pat := &pattern.Pattern[channel.Data]{
		Orig: pkt,
	}

	channelMem := make([]itfile.ChannelData, channels)
	maxCh := uint8(0)
	pos := 0
	for rowNum := 0; rowNum < int(pkt.Rows); rowNum++ {
		pat.Rows = append(pat.Rows, pattern.RowData[channel.Data]{})
		row := &pat.Rows[rowNum]
		row.Channels = make([]channel.Data, channels)
	channelLoop:
		for {
			sz, chn, err := pkt.ReadChannelData(pos, channelMem)
			if err != nil {
				return nil, 0, err
			}
			pos += sz
			if chn == nil {
				break channelLoop
			}

			channelNum := int(chn.ChannelNumber)

			cd := channel.Data{
				What:            chn.Flags,
				Note:            chn.Note,
				Instrument:      chn.Instrument,
				VolPan:          chn.VolPan,
				Effect:          channel.Command(chn.Command),
				EffectParameter: channel.DataEffect(chn.CommandData),
			}

			row.Channels[channelNum] = cd
			if maxCh < uint8(channelNum) {
				maxCh = uint8(channelNum)
			}
		}
	}

	return pat, int(maxCh), nil
}

func convertItFileToSong(f *itfile.File, features []feature.Feature) (*layout.Layout, error) {
	h, err := moduleHeaderToHeader(&f.Head)
	if err != nil {
		return nil, err
	}

	linearFrequencySlides := f.Head.Flags.IsLinearSlides()
	oldEffectMode := f.Head.Flags.IsOldEffects()
	efgLinkMode := f.Head.Flags.IsEFGLinking()
	extendedFilterRange := f.Head.Flags.IsExtendedFilterRange()

	sharedMem := channel.SharedMemory{
		LinearFreqSlides:           linearFrequencySlides,
		OldEffectMode:              oldEffectMode,
		EFGLinkMode:                efgLinkMode,
		ExtendedFilterRange:        extendedFilterRange,
		ResetMemoryAtStartOfOrder0: true,
	}

	song := layout.Layout{
		Head:          *h,
		Instruments:   make(map[uint8]instrument.Keyboard[channel.SemitoneAndSampleID]),
		Samples:       make(map[uint16]*instrument.Instrument),
		Patterns:      make([]pattern.Pattern[channel.Data], len(f.Patterns)),
		OrderList:     make([]index.Pattern, int(f.Head.OrderCount)),
		FilterPlugins: make(map[int]filter.Factory),
		Flags:         &sharedMem,
	}

	for _, block := range f.Blocks {
		switch t := block.(type) {
		case *itblock.FX:
			if filter, err := decodeFilter(t); err == nil {
				if i, err := strconv.Atoi(string(t.Identifier[2:])); err == nil {
					song.FilterPlugins[i] = filter
				}
			}
		}
	}

	for i := 0; i < int(f.Head.OrderCount); i++ {
		song.OrderList[i] = index.Pattern(f.OrderList[i])
	}

	sampleMap := make(map[uint8]pcm.Sample)

	if f.Head.Flags.IsUseInstruments() {
		for instNum, inst := range f.Instruments {
			convSettings := convertITInstrumentSettings{
				linearFrequencySlides: linearFrequencySlides,
				extendedFilterRange:   extendedFilterRange,
				useHighPassFilter:     false,
			}
			switch ii := inst.(type) {
			case *itfile.IMPIInstrumentOld:
				instMap, err := convertITInstrumentOldToInstrument(ii, f.Samples, sampleMap, convSettings, features, instNum, f.Head.TrackerCompatVersion)
				if err != nil {
					return nil, err
				}

				for _, ci := range instMap {
					addSampleWithNoteMapToSong(&song, ci.NR, instNum)
				}

			case *itfile.IMPIInstrument:
				instMap, err := convertITInstrumentToInstrument(ii, f.Samples, sampleMap, convSettings, song.FilterPlugins, features, instNum, f.Head.TrackerCompatVersion)
				if err != nil {
					return nil, err
				}

				for _, ci := range instMap {
					addSampleWithNoteMapToSong(&song, ci.NR, instNum)
				}
			}
		}
	}

	lastEnabledChannel := 0
	song.Patterns = make([]pattern.Pattern[channel.Data], len(f.Patterns))
	for patNum, pkt := range f.Patterns {
		pattern, maxCh, err := convertItPattern(pkt, len(f.Head.ChannelVol))
		if err != nil {
			return nil, err
		}
		if pattern == nil {
			continue
		}
		if lastEnabledChannel < maxCh {
			lastEnabledChannel = maxCh
		}
		song.Patterns[patNum] = *pattern
	}

	channels := make([]layout.ChannelSetting, lastEnabledChannel+1)
	for chNum := range channels {
		cs := layout.ChannelSetting{
			OutputChannelNum: chNum,
			Enabled:          true,
			InitialVolume:    volume.Volume(1),
			ChannelVolume:    volume.Volume(f.Head.ChannelVol[chNum].Value()),
			InitialPanning:   itPanning.FromItPanning(f.Head.ChannelPan[chNum]),
			Memory: channel.Memory{
				Shared: &sharedMem,
			},
		}

		cs.Memory.ResetOscillators()

		channels[chNum] = cs
	}

	song.ChannelSettings = channels

	return &song, nil
}

func decodeFilter(f *itblock.FX) (filter.Factory, error) {
	lib := f.LibraryName.String()
	name := f.UserPluginName.String()
	switch {
	case lib == "Echo" && name == "Echo":
		r := bytes.NewReader(f.Data)
		e := filter.EchoFilterFactory{}
		if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
			return nil, err
		}
		return e.Factory(), nil
	default:
		return nil, fmt.Errorf("unhandled fx lib[%s] name[%s]", lib, name)
	}
}

type noteRemap struct {
	Orig  note.Semitone
	Remap channel.SemitoneAndSampleID
	Inst  *instrument.Instrument
}

func addSampleWithNoteMapToSong(song *layout.Layout, sts []noteRemap, instNum int) {
	keyboard := instrument.Keyboard[channel.SemitoneAndSampleID]{}

	for _, st := range sts {
		sample := st.Inst
		if sample == nil {
			continue
		}
		keyboard.SetRemap(st.Orig, st.Remap)
		sid := uint16(instNum+1)<<8 | uint16(st.Remap.ID)
		song.Samples[sid] = sample
	}
	song.Instruments[uint8(instNum)+1] = keyboard
}

func readIT(r io.Reader, features []feature.Feature) (*layout.Layout, error) {
	f, err := itfile.Read(r)
	if err != nil {
		return nil, err
	}

	return convertItFileToSong(f, features)
}
