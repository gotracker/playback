package load

import (
	"errors"
	"io"
	"math"

	xmfile "github.com/gotracker/goaudiofile/music/tracked/xm"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/envelope"
	"github.com/gotracker/playback/voice/fadeout"
	"github.com/gotracker/playback/voice/loop"
	"github.com/gotracker/playback/voice/pcm"

	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/format/xm/layout"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmPeriod "github.com/gotracker/playback/format/xm/period"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/pattern"
)

func moduleHeaderToHeader(fh *xmfile.ModuleHeader) (*layout.Header, error) {
	if fh == nil {
		return nil, errors.New("file header is nil")
	}
	head := layout.Header{
		Name:         fh.GetName(),
		InitialSpeed: int(fh.DefaultSpeed),
		InitialTempo: int(fh.DefaultTempo),
		GlobalVolume: xmVolume.DefaultVolume,
		MixingVolume: xmVolume.DefaultMixingVolume,
	}
	return &head, nil
}

func xmAutoVibratoWSToProtrackerWS(vibtype uint8) uint8 {
	switch vibtype {
	case 0:
		return uint8(oscillator.WaveTableSelectSineRetrigger)
	case 1:
		return uint8(oscillator.WaveTableSelectSquareRetrigger)
	case 2:
		return uint8(oscillator.WaveTableSelectInverseSawtoothRetrigger)
	case 3:
		return uint8(oscillator.WaveTableSelectSawtoothRetrigger)
	case 4:
		return uint8(oscillator.WaveTableSelectRandomRetrigger)
	default:
		return uint8(oscillator.WaveTableSelectSineRetrigger)
	}
}

func xmInstrumentToInstrument(inst *xmfile.InstrumentHeader, linearFrequencySlides bool, features []feature.Feature) ([]*instrument.Instrument, map[int][]note.Semitone, error) {
	noteMap := make(map[int][]note.Semitone)

	var instruments []*instrument.Instrument

	for _, si := range inst.Samples {
		v := xmVolume.XmVolume(si.Volume)
		if v >= 0x40 {
			v = 0x40
		}
		sample := instrument.Instrument{
			Static: instrument.StaticValues{
				Filename:           si.GetName(),
				Name:               inst.GetName(),
				Volume:             v.Volume(),
				RelativeNoteNumber: si.RelativeNoteNumber,
				AutoVibrato: voice.AutoVibrato{
					Enabled:           (inst.VibratoDepth != 0 && inst.VibratoRate != 0),
					Sweep:             int(inst.VibratoSweep),
					WaveformSelection: xmAutoVibratoWSToProtrackerWS(inst.VibratoType),
					Depth:             float32(inst.VibratoDepth),
					Rate:              int(inst.VibratoRate),
					Factory:           oscillator.NewProtrackerOscillator,
				},
			},
			C2Spd: period.Frequency(0), // uses si.Finetune, below
		}

		if !linearFrequencySlides {
			sample.Static.AutoVibrato.Depth /= 64.0
		}

		instLen := int(si.Length)
		numChannels := 1
		format := pcm.SampleDataFormat8BitSigned

		sustainMode := xmLoopModeToLoopMode(si.Flags.LoopMode())
		sustainSettings := loop.Settings{
			Begin: int(si.LoopStart),
			End:   int(si.LoopStart + si.LoopLength),
		}

		volEnvLoopMode := loop.ModeDisabled
		volEnvLoopSettings := loop.Settings{
			Begin: int(inst.VolLoopStartPoint),
			End:   int(inst.VolLoopEndPoint),
		}
		volEnvSustainMode := loop.ModeDisabled
		volEnvSustainSettings := loop.Settings{
			Begin: int(inst.VolSustainPoint),
			End:   int(inst.VolSustainPoint),
		}

		panEnvLoopMode := loop.ModeDisabled
		panEnvLoopSettings := loop.Settings{
			Begin: int(inst.PanLoopStartPoint),
			End:   int(inst.PanLoopEndPoint),
		}
		panEnvSustainMode := loop.ModeDisabled
		panEnvSustainSettings := loop.Settings{
			Begin: int(inst.PanSustainPoint),
			End:   int(inst.PanSustainPoint),
		}

		ii := instrument.PCM{
			Loop:         &loop.Disabled{},
			MixingVolume: volume.Volume(1),
			FadeOut: fadeout.Settings{
				Mode:   fadeout.ModeOnlyIfVolEnvActive,
				Amount: volume.Volume(inst.VolumeFadeout) / 65536,
			},
			Panning: xmPanning.PanningFromXm(si.Panning),
			VolEnv: envelope.Envelope[volume.Volume]{
				Enabled: (inst.VolFlags & xmfile.EnvelopeFlagEnabled) != 0,
			},
			PanEnv: envelope.Envelope[panning.Position]{
				Enabled: (inst.PanFlags & xmfile.EnvelopeFlagEnabled) != 0,
			},
		}

		if ii.VolEnv.Enabled && (volEnvLoopSettings.End-volEnvLoopSettings.Begin) >= 0 {
			if enabled := (inst.VolFlags & xmfile.EnvelopeFlagLoopEnabled) != 0; enabled {
				volEnvLoopMode = loop.ModeNormal
			}
			if enabled := (inst.VolFlags & xmfile.EnvelopeFlagSustainEnabled) != 0; enabled {
				volEnvSustainMode = loop.ModeNormal
			}

			ii.VolEnv.Values = make([]envelope.EnvPoint[volume.Volume], int(inst.VolPoints))
			for i := range ii.VolEnv.Values {
				x1 := int(inst.VolEnv[i].X)
				y1 := uint8(inst.VolEnv[i].Y)
				var x2 int
				if i+1 < len(ii.VolEnv.Values) {
					x2 = int(inst.VolEnv[i+1].X)
				} else {
					x2 = math.MaxInt64
				}
				ii.VolEnv.Values[i].Init(x2-x1, xmVolume.XmVolume(y1).Volume())
			}
		}

		if ii.PanEnv.Enabled && (panEnvLoopSettings.End-panEnvLoopSettings.Begin) >= 0 {
			if enabled := (inst.PanFlags & xmfile.EnvelopeFlagLoopEnabled) != 0; enabled {
				panEnvLoopMode = loop.ModeNormal
			}
			if enabled := (inst.PanFlags & xmfile.EnvelopeFlagSustainEnabled) != 0; enabled {
				panEnvSustainMode = loop.ModeNormal
			}

			ii.PanEnv.Values = make([]envelope.EnvPoint[panning.Position], int(inst.VolPoints))
			for i := range ii.PanEnv.Values {
				x1 := int(inst.PanEnv[i].X)
				// XM stores pan envelope values in 0..64
				// So we have to do some gymnastics to remap the values
				panEnv01 := float64(uint8(inst.PanEnv[i].Y)) / 64
				y1 := uint8(panEnv01 * 255)
				var x2 int
				if i+1 < len(ii.PanEnv.Values) {
					x2 = int(inst.PanEnv[i+1].X)
				} else {
					x2 = math.MaxInt64
				}
				ii.PanEnv.Values[i].Init(x2-x1, xmPanning.PanningFromXm(y1))
			}
		}

		if si.Finetune != 0 {
			sample.C2Spd = xmPeriod.CalcFinetuneC2Spd(xmPeriod.DefaultC2Spd, note.Finetune(si.Finetune), linearFrequencySlides)
		}
		if sample.C2Spd == 0 {
			sample.C2Spd = period.Frequency(xmPeriod.DefaultC2Spd)
		}
		if si.Flags.IsStereo() {
			numChannels = 2
		}
		stride := numChannels
		if si.Flags.Is16Bit() {
			format = pcm.SampleDataFormat16BitLESigned
			stride *= 2
		}
		instLen /= stride
		sustainSettings.Begin /= stride
		sustainSettings.End /= stride

		ii.SustainLoop = loop.NewLoop(sustainMode, sustainSettings)
		ii.VolEnv.Loop = loop.NewLoop(volEnvLoopMode, volEnvLoopSettings)
		ii.VolEnv.Sustain = loop.NewLoop(volEnvSustainMode, volEnvSustainSettings)
		ii.PanEnv.Loop = loop.NewLoop(panEnvLoopMode, panEnvLoopSettings)
		ii.PanEnv.Sustain = loop.NewLoop(panEnvSustainMode, panEnvSustainSettings)

		samp, err := instrument.NewSample(si.SampleData, instLen, numChannels, format, features)
		if err != nil {
			return nil, nil, err
		}
		ii.Sample = samp

		sample.Inst = &ii
		instruments = append(instruments, &sample)
	}

	for st, sn := range inst.SampleNumber {
		i := int(sn)
		if i < len(instruments) {
			noteMap[i] = append(noteMap[i], note.Semitone(st))
		}
	}

	return instruments, noteMap, nil
}

func xmLoopModeToLoopMode(mode xmfile.SampleLoopMode) loop.Mode {
	switch mode {
	case xmfile.SampleLoopModeDisabled:
		return loop.ModeDisabled
	case xmfile.SampleLoopModeEnabled:
		return loop.ModeNormal
	case xmfile.SampleLoopModePingPong:
		return loop.ModePingPong
	default:
		return loop.ModeDisabled
	}
}

func convertXMInstrumentToInstrument(ih *xmfile.InstrumentHeader, linearFrequencySlides bool, features []feature.Feature) ([]*instrument.Instrument, map[int][]note.Semitone, error) {
	if ih == nil {
		return nil, nil, errors.New("instrument is nil")
	}

	return xmInstrumentToInstrument(ih, linearFrequencySlides, features)
}

func convertXmPattern(pkt xmfile.Pattern) (*pattern.Pattern[channel.Data], int) {
	pat := &pattern.Pattern[channel.Data]{
		Orig: pkt,
	}

	maxCh := uint8(0)
	for rowNum, drow := range pkt.Data {
		pat.Rows = append(pat.Rows, pattern.RowData[channel.Data]{})
		row := &pat.Rows[rowNum]
		row.Channels = make([]channel.Data, len(drow))
		for channelNum, chn := range drow {
			cd := channel.Data{
				What:            chn.Flags,
				Note:            chn.Note,
				Instrument:      chn.Instrument,
				Volume:          xmVolume.VolEffect(chn.Volume),
				Effect:          channel.Command(chn.Effect),
				EffectParameter: channel.DataEffect(chn.EffectParameter),
			}
			row.Channels[channelNum] = cd
			if maxCh < uint8(channelNum) {
				maxCh = uint8(channelNum)
			}
		}
	}

	return pat, int(maxCh)
}

func convertXmFileToSong(f *xmfile.File, features []feature.Feature) (*layout.Song, error) {
	h, err := moduleHeaderToHeader(&f.Head)
	if err != nil {
		return nil, err
	}

	linearFrequencySlides := f.Head.Flags.IsLinearSlides()

	song := layout.Song{
		Head:              *h,
		Instruments:       make(map[uint8]*instrument.Instrument),
		InstrumentNoteMap: make(map[uint8]map[note.Semitone]*instrument.Instrument),
		Patterns:          make([]pattern.Pattern[channel.Data], len(f.Patterns)),
		OrderList:         make([]index.Pattern, int(f.Head.SongLength)),
	}

	for i := 0; i < int(f.Head.SongLength); i++ {
		song.OrderList[i] = index.Pattern(f.Head.OrderTable[i])
	}

	for instNum, ih := range f.Instruments {
		samples, noteMap, err := convertXMInstrumentToInstrument(&ih, linearFrequencySlides, features)
		if err != nil {
			return nil, err
		}
		for _, sample := range samples {
			if sample == nil {
				continue
			}
			id := channel.SampleID{
				InstID: uint8(instNum + 1),
			}
			sample.Static.ID = id
			song.Instruments[id.InstID] = sample
		}
		for i, sts := range noteMap {
			sample := samples[i]
			id, ok := sample.Static.ID.(channel.SampleID)
			if !ok {
				continue
			}
			inm, ok := song.InstrumentNoteMap[id.InstID]
			if !ok {
				inm = make(map[note.Semitone]*instrument.Instrument)
				song.InstrumentNoteMap[id.InstID] = inm
			}
			for _, st := range sts {
				inm[st] = samples[i]
			}
		}
	}

	lastEnabledChannel := 0
	song.Patterns = make([]pattern.Pattern[channel.Data], len(f.Patterns))
	for patNum, pkt := range f.Patterns {
		pattern, maxCh := convertXmPattern(pkt)
		if pattern == nil {
			continue
		}
		if lastEnabledChannel < maxCh {
			lastEnabledChannel = maxCh
		}
		song.Patterns[patNum] = *pattern
	}

	sharedMem := channel.SharedMemory{
		LinearFreqSlides:           linearFrequencySlides,
		ResetMemoryAtStartOfOrder0: true,
	}

	channels := make([]layout.ChannelSetting, lastEnabledChannel+1)
	for chNum := range channels {
		cs := layout.ChannelSetting{
			Enabled:        true,
			InitialVolume:  xmVolume.DefaultVolume,
			InitialPanning: xmPanning.DefaultPanning,
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

func readXM(r io.Reader, features []feature.Feature) (*layout.Song, error) {
	f, err := xmfile.Read(r)
	if err != nil {
		return nil, err
	}

	return convertXmFileToSong(f, features)
}
