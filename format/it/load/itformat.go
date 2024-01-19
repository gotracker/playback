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

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/common"
	"github.com/gotracker/playback/format/it/channel"
	"github.com/gotracker/playback/format/it/layout"
	itPanning "github.com/gotracker/playback/format/it/panning"
	"github.com/gotracker/playback/format/it/settings"
	itSystem "github.com/gotracker/playback/format/it/system"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

func moduleHeaderToHeader(fh *itfile.ModuleHeader) (*layout.Header, error) {
	if fh == nil {
		return nil, errors.New("file header is nil")
	}
	head := layout.Header{
		Name:             fh.GetName(),
		InitialSpeed:     int(fh.InitialSpeed),
		InitialTempo:     int(fh.InitialTempo),
		GlobalVolume:     itVolume.FineVolume(fh.GlobalVolume),
		MixingVolume:     itVolume.FineVolume(fh.MixingVolume),
		LinearFreqSlides: fh.Flags.IsLinearSlides(),
		InitialOrder:     0,
	}
	switch {
	case fh.TrackerCompatVersion < 0x200:
		head.MixingVolume = max(itVolume.FineVolume(fh.MixingVolume*2), itVolume.MaxItFineVolume)
	case fh.TrackerCompatVersion >= 0x200:
		head.MixingVolume = itVolume.FineVolume(fh.MixingVolume)
	}
	return &head, nil
}

func convertItPattern[TPeriod period.Period](pkt itfile.PackedPattern, channels int) (song.Pattern, int, error) {
	pat := make(song.Pattern, pkt.Rows)

	channelMem := make([]itfile.ChannelData, channels)
	maxCh := uint8(0)
	pos := 0
	for rowNum := 0; rowNum < int(pkt.Rows); rowNum++ {
		row := make(layout.Row[TPeriod], channels)
		pat[rowNum] = row
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

			cd := channel.Data[TPeriod]{
				What:            chn.Flags,
				Note:            chn.Note,
				Instrument:      chn.Instrument,
				VolPan:          chn.VolPan,
				Effect:          channel.Command(chn.Command),
				EffectParameter: channel.DataEffect(chn.CommandData),
			}

			row[channelNum] = cd
			if maxCh < uint8(channelNum) {
				maxCh = uint8(channelNum)
			}
		}
	}

	return pat, int(maxCh), nil
}

func convertItFileToSong(f *itfile.File, features []feature.Feature) (song.Data, error) {
	if f.Head.Flags.IsLinearSlides() {
		return convertItFileToTypedSong[period.Linear](f, features)
	} else {
		return convertItFileToTypedSong[period.Amiga](f, features)
	}
}

func convertItFileToTypedSong[TPeriod period.Period](f *itfile.File, features []feature.Feature) (*layout.Song[TPeriod], error) {
	h, err := moduleHeaderToHeader(&f.Head)
	if err != nil {
		return nil, err
	}

	linearFrequencySlides := f.Head.Flags.IsLinearSlides()
	oldEffectMode := f.Head.Flags.IsOldEffects()
	efgLinkMode := f.Head.Flags.IsEFGLinking()
	stereoMode := f.Head.Flags.IsStereo()
	vol0Enabled := f.Head.Flags.IsVol0Optimizations()

	ms := settings.GetMachineSettings[TPeriod]()

	songData := &layout.Song[TPeriod]{
		BaseSong: common.BaseSong[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]{
			System:       itSystem.ITSystem,
			MS:           ms,
			Name:         h.Name,
			InitialBPM:   h.InitialTempo,
			InitialTempo: h.InitialSpeed,
			GlobalVolume: h.GlobalVolume,
			MixingVolume: h.MixingVolume,
			InitialOrder: h.InitialOrder,
			Instruments:  make([]*instrument.Instrument[TPeriod, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], f.Head.InstrumentCount),
			Patterns:     make([]song.Pattern, len(f.Patterns)),
			OrderList:    make([]index.Pattern, int(f.Head.OrderCount)),
		},
		InstrumentNoteMap: make(map[uint8]map[note.Semitone]layout.NoteInstrument[TPeriod]),
		FilterPlugins:     make(map[int]filter.Factory),
	}

	for _, block := range f.Blocks {
		switch t := block.(type) {
		case *itblock.FX:
			if filter, err := decodeFilter(t); err == nil {
				if i, err := strconv.Atoi(string(t.Identifier[2:])); err == nil {
					songData.FilterPlugins[i] = filter
				}
			}
		}
	}

	for i := 0; i < int(f.Head.OrderCount); i++ {
		songData.OrderList[i] = index.Pattern(f.OrderList[i])
	}

	if f.Head.Flags.IsUseInstruments() {
		for instNum, inst := range f.Instruments {
			convSettings := convertITInstrumentSettings{
				linearFrequencySlides: linearFrequencySlides,
				extendedFilterRange:   (f.Head.Flags & 0x1000) != 0, // OpenMPT hack to introduce extended filter ranges
				useHighPassFilter:     false,
			}
			switch ii := inst.(type) {
			case *itfile.IMPIInstrumentOld:
				instMap, err := convertITInstrumentOldToInstrument(ii, ms.PeriodConverter, f.Samples, convSettings, features)
				if err != nil {
					return nil, err
				}

				for _, ci := range instMap {
					addSampleWithNoteMapToSong(songData, ci.Inst, ci.NR, instNum)
				}

			case *itfile.IMPIInstrument:
				instMap, err := convertITInstrumentToInstrument(ii, ms.PeriodConverter, f.Samples, convSettings, songData.FilterPlugins, features)
				if err != nil {
					return nil, err
				}

				for _, ci := range instMap {
					addSampleWithNoteMapToSong(songData, ci.Inst, ci.NR, instNum)
				}
			}
		}
	}

	lastEnabledChannel := 0
	for patNum, pkt := range f.Patterns {
		p, maxCh, err := convertItPattern[TPeriod](pkt, len(f.Head.ChannelVol))
		if err != nil {
			return nil, err
		}
		if p == nil {
			continue
		}
		if lastEnabledChannel < maxCh {
			lastEnabledChannel = maxCh
		}
		songData.Patterns[patNum] = p
	}

	sharedMem := channel.SharedMemory{
		OldEffectMode:              oldEffectMode,
		EFGLinkMode:                efgLinkMode,
		ResetMemoryAtStartOfOrder0: true,
	}

	channels := make([]layout.ChannelSetting, lastEnabledChannel+1)
	for chNum := range channels {
		cs := layout.ChannelSetting{
			OutputChannelNum: chNum,
			Enabled:          true,
			Muted:            false,
			InitialVolume:    itVolume.Volume(itVolume.DefaultItVolume),
			ChannelVolume:    min(itVolume.FineVolume(f.Head.ChannelVol[chNum]*2), itVolume.MaxItFineVolume),
			PanEnabled:       stereoMode,
			InitialPanning:   itPanning.Panning(f.Head.ChannelPan[chNum]),
			Memory: channel.Memory{
				Shared: &sharedMem,
			},
			Vol0OptEnabled: vol0Enabled,
		}

		channels[chNum] = cs
	}

	songData.ChannelSettings = channels
	return songData, nil
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
	Remap note.Semitone
}

func addSampleWithNoteMapToSong[TPeriod period.Period](song *layout.Song[TPeriod], sample *instrument.Instrument[TPeriod, itVolume.FineVolume, itVolume.Volume, itPanning.Panning], sts []noteRemap, instNum int) {
	if sample == nil {
		return
	}
	id := channel.SampleID{
		InstID: uint8(instNum + 1),
	}
	sample.Static.ID = id
	song.Instruments[instNum] = sample

	id, ok := sample.Static.ID.(channel.SampleID)
	if !ok {
		return
	}
	inm, ok := song.InstrumentNoteMap[id.InstID]
	if !ok {
		inm = make(map[note.Semitone]layout.NoteInstrument[TPeriod])
		song.InstrumentNoteMap[id.InstID] = inm
	}
	for _, st := range sts {
		inm[st.Orig] = layout.NoteInstrument[TPeriod]{
			NoteRemap: st.Remap,
			Inst:      sample,
		}
	}
}

func readIT(r io.Reader, features []feature.Feature) (song.Data, error) {
	f, err := itfile.Read(r)
	if err != nil {
		return nil, err
	}

	return convertItFileToSong(f, features)
}
