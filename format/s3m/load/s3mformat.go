package load

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/format/s3m/channel"
	"github.com/gotracker/playback/format/s3m/layout"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	"github.com/gotracker/playback/format/s3m/pattern"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
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
		GlobalVolume: s3mVolume.VolumeFromS3M(fh.GlobalVolume),
		Stereo:       (fh.MixingVolume & 0x80) != 0,
	}

	z := uint32(fh.MixingVolume & 0x7f)
	if z < 0x10 {
		z = 0x10
	}
	head.MixingVolume = volume.Volume(z) / volume.Volume(0x80)

	return &head, nil
}

func scrsNoneToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSNoneHeader) (*instrument.Instrument, error) {
	sample := instrument.Instrument{
		Static: instrument.StaticValues{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.VolumeFromS3M(si.Volume),
		},
		SampleRate: period.Frequency(si.C2Spd.Lo),
	}
	return &sample, nil
}

func scrsDp30ToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSDigiplayerHeader, signedSamples bool, features []feature.Feature) (*instrument.Instrument, error) {
	sample := instrument.Instrument{
		Static: instrument.StaticValues{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.VolumeFromS3M(si.Volume),
		},
		SampleRate: period.Frequency(si.C2Spd.Lo),
	}
	if sample.SampleRate == 0 {
		sample.SampleRate = period.Frequency(s3mfile.DefaultC2Spd)
	}

	instLen := int(si.Length.Lo)
	numChannels := 1
	format := pcm.SampleDataFormat8BitUnsigned

	sustainMode := loop.ModeDisabled
	sustainSettings := loop.Settings{
		Begin: int(si.LoopBegin.Lo),
		End:   int(si.LoopEnd.Lo),
	}

	idata := instrument.PCM{
		Loop:         &loop.Disabled{},
		Panning:      panning.CenterAhead,
		MixingVolume: volume.Volume(1),
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

func scrsOpl2ToInstrument(scrs *s3mfile.SCRSFull, si *s3mfile.SCRSAdlibHeader) (*instrument.Instrument, error) {
	inst := instrument.Instrument{
		Static: instrument.StaticValues{
			Filename: scrs.Head.GetFilename(),
			Name:     si.GetSampleName(),
			Volume:   s3mVolume.VolumeFromS3M(si.Volume),
		},
		SampleRate: period.Frequency(si.C2Spd.Lo),
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

func convertSCRSFullToInstrument(scrs *s3mfile.SCRSFull, signedSamples bool, features []feature.Feature) (*instrument.Instrument, error) {
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

func convertS3MPackedPattern(pkt s3mfile.PackedPattern, numRows uint8) (pattern.Pattern, int) {
	pat := make(pattern.Pattern, numRows)

	buffer := bytes.NewBuffer(pkt.Data)

	maxCh := uint8(0)
	for rowNum := uint8(0); rowNum < numRows; rowNum++ {
		row := make(pattern.Row, 0)
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
			temp.Volume = s3mfile.EmptyVolume
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

	song := layout.Song{
		Head:        *h,
		Instruments: make([]*instrument.Instrument, len(f.InstrumentPointers)),
		OrderList:   make([]index.Pattern, len(f.OrderList)),
	}

	signedSamples := false
	if f.Head.FileFormatInformation == 1 {
		signedSamples = true
	}

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
		song.OrderList[i] = index.Pattern(o)
	}

	song.Instruments = make([]*instrument.Instrument, len(f.Instruments))
	for instNum, scrs := range f.Instruments {
		sample, err := convertSCRSFullToInstrument(&scrs, signedSamples, features)
		if err != nil {
			return nil, err
		}
		if sample == nil {
			continue
		}
		sample.Static.ID = channel.InstID(uint8(instNum + 1))
		song.Instruments[instNum] = sample
	}

	lastEnabledChannel := 0
	song.Patterns = make([]pattern.Pattern, len(f.Patterns))
	for patNum, pkt := range f.Patterns {
		pattern, maxCh := convertS3MPackedPattern(pkt, getPatternLen(patNum))
		if pattern == nil {
			continue
		}
		if lastEnabledChannel < maxCh {
			lastEnabledChannel = maxCh
		}
		song.Patterns[patNum] = pattern
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

	channels := make([]layout.ChannelSetting, 0)
	for chNum, ch := range f.ChannelSettings {
		chn := ch.GetChannel()
		cs := layout.ChannelSetting{
			Enabled:          ch.IsEnabled(),
			Category:         chn.GetChannelCategory(),
			OutputChannelNum: int(ch.GetChannel() & 0x07),
			InitialVolume:    s3mVolume.DefaultVolume,
			InitialPanning:   s3mPanning.DefaultPanning,
			Memory: channel.Memory{
				Shared: &sharedMem,
			},
		}

		cs.Memory.ResetOscillators()

		pf := f.Panning[chNum]
		if pf.IsValid() {
			cs.InitialPanning = s3mPanning.PanningFromS3M(pf.Value())
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

	song.ChannelSettings = channels[:lastEnabledChannel+1]

	return &song, nil
}

func readS3M(r io.Reader, features []feature.Feature) (*layout.Song, error) {
	f, err := s3mfile.Read(r)
	if err != nil {
		return nil, err
	}

	return convertS3MFileToSong(f, func(patNum int) uint8 {
		return 64
	}, features, false)
}
