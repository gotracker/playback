package load

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/format/s3m/layout"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mSystem "github.com/gotracker/playback/format/s3m/system"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"
)

func moduleHeaderToHeader(fh *s3mfile.ModuleHeader) (*layout.Header, error) {
	if fh == nil {
		return nil, errors.New("file header is nil")
	}
	head := layout.Header{
		Name:         fh.GetName(),
		InitialSpeed: int(fh.InitialSpeed),
		InitialTempo: int(fh.InitialTempo),
		GlobalVolume: s3mVolume.Volume(fh.GlobalVolume),
		MixingVolume: s3mVolume.FineVolume(fh.MixingVolume &^ 0x80),
		InitialOrder: 0,
	}

	z := uint32(fh.MixingVolume & 0x7f)
	if z < 0x10 {
		z = 0x10
	}
	head.MixingVolume = s3mVolume.FineVolume(z)

	return &head, nil
}

func scrsNoneToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSNoneHeader) (*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], error) {
	sample := instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		Static: instrument.StaticValues[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.Volume(si.Volume),
		},
		SampleRate: frequency.Frequency(si.C2Spd.Lo),
	}
	return &sample, nil
}

func scrsDp30ToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSDigiplayerHeader, signedSamples bool, features []feature.Feature) (*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], error) {
	sample := instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		Static: instrument.StaticValues[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.Volume(si.Volume),
		},
		SampleRate: frequency.Frequency(si.C2Spd.Lo),
	}
	if sample.SampleRate == 0 {
		sample.SampleRate = frequency.Frequency(s3mfile.DefaultC2Spd)
	}

	instLen := int(si.Length.Lo)
	numChannels := 1
	format := pcm.SampleDataFormat8BitUnsigned

	sustainMode := loop.ModeDisabled
	sustainSettings := loop.Settings{
		Begin: int(si.LoopBegin.Lo),
		End:   int(si.LoopEnd.Lo),
	}

	idata := instrument.PCM[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		Loop: &loop.Disabled{},
		FadeOut: fadeout.Settings{
			Mode:   fadeout.ModeDisabled,
			Amount: volume.Volume(0),
		},
	}
	if signedSamples {
		format = pcm.SampleDataFormat8BitSigned
	}
	if si.Flags.IsLooped() {
		sustainMode = loop.ModeNormal
	}
	if si.Flags.IsStereo() {
		numChannels = 2
	}
	if si.Flags.Is16BitSample() {
		if signedSamples {
			format = pcm.SampleDataFormat16BitLESigned
		} else {
			format = pcm.SampleDataFormat16BitLEUnsigned
		}
	}

	idata.SustainLoop = loop.NewLoop(sustainMode, sustainSettings)

	samp, err := instrument.NewSample(scrs.Sample, instLen, numChannels, format, features)
	if err != nil {
		return nil, err
	}
	idata.Sample = samp

	sample.Inst = &idata
	return &sample, nil
}

func scrsOpl2ToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSAdlibHeader) (*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], error) {
	inst := instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		Static: instrument.StaticValues[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.Volume(si.Volume),
		},
		SampleRate: frequency.Frequency(si.C2Spd.Lo),
	}

	idata := instrument.OPL2{
		Modulator: instrument.OPL2OperatorData{
			KeyScaleRateSelect:  si.OPL2.ModulatorKeyScaleRateSelect(),
			Sustain:             si.OPL2.ModulatorSustain(),
			Vibrato:             si.OPL2.ModulatorVibrato(),
			Tremolo:             si.OPL2.ModulatorTremolo(),
			FrequencyMultiplier: uint8(si.OPL2.ModulatorFrequencyMultiplier()),
			KeyScaleLevel:       uint8(si.OPL2.ModulatorKeyScaleLevel()),
			Volume:              uint8(si.OPL2.ModulatorVolume()),
			AttackRate:          si.OPL2.ModulatorAttackRate(),
			DecayRate:           si.OPL2.ModulatorDecayRate(),
			SustainLevel:        si.OPL2.ModulatorSustainLevel(),
			ReleaseRate:         si.OPL2.ModulatorReleaseRate(),
			WaveformSelection:   uint8(si.OPL2.ModulatorWaveformSelection()),
		},
		Carrier: instrument.OPL2OperatorData{
			KeyScaleRateSelect:  si.OPL2.CarrierKeyScaleRateSelect(),
			Sustain:             si.OPL2.CarrierSustain(),
			Vibrato:             si.OPL2.CarrierVibrato(),
			Tremolo:             si.OPL2.CarrierTremolo(),
			FrequencyMultiplier: uint8(si.OPL2.CarrierFrequencyMultiplier()),
			KeyScaleLevel:       uint8(si.OPL2.CarrierKeyScaleLevel()),
			Volume:              uint8(si.OPL2.CarrierVolume()),
			AttackRate:          si.OPL2.CarrierAttackRate(),
			DecayRate:           si.OPL2.CarrierDecayRate(),
			SustainLevel:        si.OPL2.CarrierSustainLevel(),
			ReleaseRate:         si.OPL2.CarrierReleaseRate(),
			WaveformSelection:   uint8(si.OPL2.CarrierWaveformSelection()),
		},
		ModulationFeedback: uint8(si.OPL2.ModulationFeedback()),
		AdditiveSynthesis:  si.OPL2.AdditiveSynthesis(),
	}

	inst.Inst = &idata
	return &inst, nil
}

func convertSCRSFullToInstrument(scrs *s3mfile.SCRSFull, signedSamples bool, features []feature.Feature) (*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], error) {
	if scrs == nil {
		return nil, errors.New("scrs is nil")
	}

	switch si := scrs.Ancillary.(type) {
	case nil:
		return nil, errors.New("scrs ancillary is nil")
	case *s3mfile.SCRSNoneHeader:
		return scrsNoneToInstrument(scrs, si)
	case *s3mfile.SCRSDigiplayerHeader:
		return scrsDp30ToInstrument(scrs, si, signedSamples, features)
	case *s3mfile.SCRSAdlibHeader:
		return scrsOpl2ToInstrument(scrs, si)
	default:
	}

	return nil, errors.New("unhandled scrs ancillary type")
}

func convertS3MPackedPattern(pkt s3mfile.PackedPattern, numRows uint8) (song.Pattern, int) {
	pat := make(song.Pattern, numRows)

	buffer := bytes.NewBuffer(pkt.Data)

	maxCh := uint8(0)
	for rowNum := uint8(0); rowNum < numRows; rowNum++ {
		row := make(layout.Row, 0)
	channelLoop:
		for {
			var what s3mfile.PatternFlags
			if err := binary.Read(buffer, binary.LittleEndian, &what); err != nil {
				panic(err)
			}

			if what == 0 {
				break channelLoop
			}

			channelNum := what.Channel()
			for len(row) <= int(channelNum) {
				row = append(row, channel.Data{})
			}
			temp := &row[channelNum]
			if maxCh < channelNum {
				maxCh = channelNum
			}

			temp.What = what
			temp.Note = 0
			temp.Instrument = 0
			temp.Volume = s3mVolume.Volume(s3mfile.EmptyVolume)
			temp.Command = 0
			temp.Info = 0

			if temp.What.HasNote() {
				if err := binary.Read(buffer, binary.LittleEndian, &temp.Note); err != nil {
					panic(err)
				}
				if err := binary.Read(buffer, binary.LittleEndian, &temp.Instrument); err != nil {
					panic(err)
				}
			}

			if temp.What.HasVolume() {
				if err := binary.Read(buffer, binary.LittleEndian, &temp.Volume); err != nil {
					panic(err)
				}
			}

			if temp.What.HasCommand() {
				if err := binary.Read(buffer, binary.LittleEndian, &temp.Command); err != nil {
					panic(err)
				}
				if err := binary.Read(buffer, binary.LittleEndian, &temp.Info); err != nil {
					panic(err)
				}
			}
		}
		pat[rowNum] = row
	}

	return pat, int(maxCh)
}

func convertS3MFileToSong(f *s3mfile.File, getPatternLen func(patNum int) uint8, features []feature.Feature, wasModFile bool) (*layout.Song, error) {
	h, err := moduleHeaderToHeader(&f.Head)
	if err != nil {
		return nil, err
	}

	s := layout.Song{
		System:      s3mSystem.S3MSystem,
		Head:        *h,
		Instruments: make([]*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], len(f.InstrumentPointers)),
		OrderList:   make([]index.Pattern, len(f.OrderList)),
	}

	signedSamples := false
	if f.Head.FileFormatInformation == 1 {
		signedSamples = true
	}

	stereoMode := (f.Head.MixingVolume & 0x80) != 0
	st2Vibrato := (f.Head.Flags & 0x0001) != 0
	st2Tempo := (f.Head.Flags & 0x0002) != 0
	amigaSlides := (f.Head.Flags & 0x0004) != 0
	zeroVolOpt := (f.Head.Flags & 0x0008) != 0
	amigaLimits := (f.Head.Flags & 0x0010) != 0
	sbFilterEnable := (f.Head.Flags & 0x0020) != 0
	st300volSlides := (f.Head.Flags & 0x0040) != 0
	if f.Head.TrackerVersion == 0x1300 {
		st300volSlides = true
	}
	//ptrSpecialIsValid := (f.Head.Flags & 0x0080) != 0

	for i, o := range f.OrderList {
		s.OrderList[i] = index.Pattern(o)
	}

	s.Instruments = make([]*instrument.Instrument[s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], len(f.Instruments))
	for instNum, scrs := range f.Instruments {
		sample, err := convertSCRSFullToInstrument(&scrs, signedSamples, features)
		if err != nil {
			return nil, err
		}
		if sample == nil {
			continue
		}
		sample.Static.ID = channel.InstID(uint8(instNum + 1))
		s.Instruments[instNum] = sample
	}

	maxPatternChannel := 0
	s.Patterns = make([]song.Pattern, len(f.Patterns))
	for patNum, pkt := range f.Patterns {
		pattern, maxCh := convertS3MPackedPattern(pkt, getPatternLen(patNum))
		if pattern == nil {
			continue
		}
		if maxPatternChannel < maxCh {
			maxPatternChannel = maxCh
		}
		s.Patterns[patNum] = pattern
	}

	sharedMem := channel.SharedMemory{
		VolSlideEveryFrame:         st300volSlides,
		LowPassFilterEnable:        sbFilterEnable,
		ResetMemoryAtStartOfOrder0: true,
		ST2Vibrato:                 st2Vibrato,
		ST2Tempo:                   st2Tempo,
		AmigaSlides:                amigaSlides,
		ZeroVolOptimization:        zeroVolOpt,
		AmigaLimits:                amigaLimits,
		ModCompatibility:           wasModFile,
	}

	channels := make([]layout.ChannelSetting, 0, maxPatternChannel+1)
	lastEnabledChannel := 0
	for chNum, ch := range f.ChannelSettings {
		chn := ch.GetChannel()
		cs := layout.ChannelSetting{
			Enabled:          ch.IsEnabled(),
			Category:         chn.GetChannelCategory(),
			OutputChannelNum: int(ch.GetChannel() & 0x07),
			InitialVolume:    s3mVolume.Volume(s3mfile.DefaultVolume),
			PanEnabled:       stereoMode,
			InitialPanning:   s3mPanning.DefaultPanning,
			Memory: channel.Memory{
				Shared: &sharedMem,
			},
			DefaultFilterName: "",
		}

		if sbFilterEnable {
			cs.DefaultFilterName = "amigalpf"
		}

		pf := f.Panning[chNum]
		if pf.IsValid() {
			cs.InitialPanning = s3mPanning.Panning(pf.Value())
		} else {
			switch cs.Category {
			case s3mfile.ChannelCategoryPCMLeft:
				cs.InitialPanning = s3mPanning.DefaultPanningLeft
				cs.OutputChannelNum = int(chn - s3mfile.ChannelIDL1)
			case s3mfile.ChannelCategoryPCMRight:
				cs.InitialPanning = s3mPanning.DefaultPanningRight
				cs.OutputChannelNum = int(chn - s3mfile.ChannelIDR1)
			}
		}

		channels = append(channels, cs)
		if cs.Enabled && lastEnabledChannel < chNum {
			lastEnabledChannel = chNum
		}
	}

	s.NumChannels = lastEnabledChannel + 1
	s.ChannelSettings = channels[:maxPatternChannel+1]

	return &s, nil
}

func readS3M(r io.Reader, features []feature.Feature) (song.Data, error) {
	f, err := s3mfile.Read(r)
	if err != nil {
		return nil, err
	}

	return convertS3MFileToSong(f, func(patNum int) uint8 {
		return 64
	}, features, false)
}
