package machine

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/opl2"
)

type typeLookup struct {
	p   reflect.Type
	gv  reflect.Type
	cmv reflect.Type
	cv  reflect.Type
	cp  reflect.Type
}

func (t typeLookup) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v, %v]", t.p, t.gv, t.cmv, t.cv, t.cp)
}

type factory func(songData song.Data, us settings.UserSettings) (MachineTicker, error)

var factoryRegistry = make(map[typeLookup]factory)

func RegisterMachine[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](ms *settings.MachineSettings[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) {
	var (
		p   TPeriod
		gv  TGlobalVolume
		cmv TMixingVolume
		cv  TVolume
		cp  TPanning
	)
	tl := typeLookup{
		p:   reflect.TypeOf(p),
		gv:  reflect.TypeOf(gv),
		cmv: reflect.TypeOf(cmv),
		cv:  reflect.TypeOf(cv),
		cp:  reflect.TypeOf(cp),
	}
	factoryRegistry[tl] = func(songData song.Data, us settings.UserSettings) (MachineTicker, error) {
		var m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
		// we have to use the songData's machine settings
		var ok bool
		m.ms, ok = songData.GetMachineSettings().(*settings.MachineSettings[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
		if !ok {
			return nil, errors.New("invalid machine settings from songdata")
		}

		m.songData = songData
		m.us = us

		order := songData.GetInitialOrder()
		if o, set := us.Start.Order.Get(); set {
			order = index.Order(o)
		}

		var row index.Row
		if r, set := us.Start.Row.Get(); set {
			row = index.Row(r)
		}

		sys := songData.GetSystem()

		bpm := songData.GetInitialBPM()
		if us.Start.BPM != 0 {
			bpm = us.Start.BPM
		}

		tempo := songData.GetInitialTempo()
		if us.Start.Tempo != 0 {
			tempo = us.Start.Tempo
		}

		if err := m.SetBPM(bpm); err != nil {
			return nil, err
		}
		if err := m.SetTempo(tempo); err != nil {
			return nil, err
		}
		gv, err := song.GetGlobalVolume[TGlobalVolume](songData)
		if err != nil {
			return nil, err
		}
		if err := m.SetGlobalVolume(gv); err != nil {
			return nil, err
		}
		mv, err := song.GetMixingVolume[TMixingVolume](songData)
		if err != nil {
			return nil, err
		}
		if err := m.SetMixingVolume(mv.ToVolume()); err != nil {
			return nil, err
		}

		channels := songData.GetNumChannels()

		m.channels = make([]channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], channels)
		m.actualOutputs = make([]render.Channel[TPeriod], channels)

		m.opl2Enabled = songData.IsOPL2Enabled()

		mpnpc := sys.GetMaxPastNotesPerChannel()
		if mpnpc > 0 {
			m.virtualOutputs = make([]render.Channel[TPeriod], channels*mpnpc)
		}

		for i := 0; i < channels; i++ {
			ch := index.Channel(i)
			cs := songData.GetChannelSettings(ch)

			rc := &m.actualOutputs[ch]
			rc.OutputFilter = nil

			if cs.IsDefaultFilterEnabled() {
				name := cs.GetDefaultFilterName()
				filtFactory, err := ms.GetFilterFactory(name)
				if err != nil {
					return nil, err
				}

				filt, err := filtFactory(sys.GetCommonRate())
				if err != nil {
					return nil, err
				}

				rc.OutputFilter = filt
			}
			rc.GetOPL2Chip = func() opl2.Chip {
				return m.opl2
			}

			initialVolume, err := song.GetChannelInitialVolume[TVolume](cs)
			if err != nil {
				return nil, err
			}

			initialMixing, err := song.GetChannelMixingVolume[TMixingVolume](cs)
			if err != nil {
				return nil, err
			}

			initialPan, err := song.GetChannelInitialPanning[TPanning](cs)
			if err != nil {
				return nil, err
			}

			rc.GlobalVolume = volume.Volume(1)

			c := &m.channels[ch]
			c.enabled = cs.IsEnabled()
			c.cv = m.ms.VoiceFactory.NewVoice(voice.VoiceConfig[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]{
				PC:               ms.PeriodConverter,
				OPLChannel:       cs.GetOPLChannel(),
				InitialVolume:    initialVolume,
				InitialMixing:    initialMixing,
				PanEnabled:       cs.IsPanEnabled(),
				InitialPan:       initialPan,
				Vol0Optimization: cs.GetVol0OptimizationSettings(),
			})
			c.memory = cs.GetMemory()
			rc.StartVoice(c.cv, func() {}) // can't remove this channel, as it's hard-wired into actual
			c.target.ActionTick.Reset()

			c.nna = note.ActionCut
			if c.osc[OscillatorVibrato], err = ms.GetVibratoFactory(); err != nil {
				return nil, err
			}
			if c.osc[OscillatorTremolo], err = ms.GetTremoloFactory(); err != nil {
				return nil, err
			}
			if c.osc[OscillatorPanbrello], err = ms.GetPanbrelloFactory(); err != nil {
				return nil, err
			}
			cmv, err := song.GetChannelMixingVolume[TMixingVolume](cs)
			if err != nil {
				return nil, fmt.Errorf("channel[%d]: %w", i, err)
			}
			m.SetChannelMute(ch, cs.IsMuted())
			m.SetChannelMixingVolume(ch, cmv)
			cv, err := song.GetChannelInitialVolume[TVolume](cs)
			if err != nil {
				return nil, fmt.Errorf("channel[%d]: %w", i, err)
			}
			m.SetChannelVolume(ch, cv)
			cp, err := song.GetChannelInitialPanning[TPanning](cs)
			if err != nil {
				return nil, fmt.Errorf("channel[%d]: %w", i, err)
			}
			m.SetChannelPan(ch, cp)
		}

		if err := initTick(&m.ticker, &m, tickerSettings{
			InitialOrder:          order,
			InitialRow:            row,
			SongLoopStartingOrder: 0,
			SongLoopCount:         us.SongLoopCount,
			PlayUntilOrder:        us.PlayUntil.Order,
			PlayUntilRow:          us.PlayUntil.Row,
		}); err != nil {
			return nil, err
		}

		return &m, nil
	}
}
