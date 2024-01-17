package layout

import (
	"reflect"
	"time"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/system"
)

// Song is the full definition of the song data of an Song file
type Song struct {
	System          system.System
	MS              *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]
	Head            Header
	Instruments     []*instrument.Instrument[period.Amiga, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]
	Patterns        []song.Pattern
	ChannelSettings []ChannelSetting
	ChannelOrders   []index.Channel
	NumChannels     int
	OrderList       []index.Pattern
}

func (Song) GetPeriodType() reflect.Type {
	var p period.Amiga
	return reflect.TypeOf(p)
}

func (s Song) GetGlobalVolumeType() reflect.Type {
	return reflect.TypeOf(s.Head.GlobalVolume)
}

func (s Song) GetChannelMixingVolumeType() reflect.Type {
	return reflect.TypeOf(s.Head.MixingVolume)
}

func (s Song) GetChannelVolumeType() reflect.Type {
	var cs ChannelSetting
	return reflect.TypeOf(cs.InitialVolume)
}

func (s Song) GetChannelPanningType() reflect.Type {
	var cs ChannelSetting
	return reflect.TypeOf(cs.InitialPanning)
}

// GetOrderList returns the list of all pattern orders for the song
func (s Song) GetOrderList() []index.Pattern {
	return s.OrderList
}

// GetInitialBPM returns the initial "tempo" (number of beats per minute) for the song
func (s Song) GetInitialBPM() int {
	return s.Head.InitialTempo
}

// GetInitialTempo returns the initial "speed" (number of ticks per row) for the song
func (s Song) GetInitialTempo() int {
	return s.Head.InitialSpeed
}

// GetGlobalVolumeGeneric returns the initial global volume for the song
func (s Song) GetGlobalVolumeGeneric() volume.Volume {
	return s.Head.GlobalVolume.ToVolume()
}

// GetGlobalVolume returns the initial global volume for the song
func (s Song) GetGlobalVolume() s3mVolume.Volume {
	return s.Head.GlobalVolume
}

// GetMixingVolumeGeneric returns the initial mixing volume for the song
func (s Song) GetMixingVolumeGeneric() volume.Volume {
	return s.Head.MixingVolume.ToVolume()
}

// GetMixingVolume returns the initial mixing volume for the song
func (s Song) GetMixingVolume() s3mVolume.FineVolume {
	return s.Head.MixingVolume
}

const durationPerBpm = time.Duration(2500) * time.Millisecond

// GetTickDuration calculates the duration of a tick at a particular BPM
func (s Song) GetTickDuration(bpm int) time.Duration {
	if bpm == 0 {
		return 0
	}

	return durationPerBpm / time.Duration(bpm)
}

// GetPattern returns a specific pattern indexed by `patNum`
func (s Song) GetPattern(patNum index.Pattern) (song.Pattern, error) {
	if int(patNum) >= len(s.Patterns) {
		return nil, song.ErrStopSong
	}
	return s.Patterns[patNum], nil
}

// GetPatternByOrder returns the pattern specified by the order index provided
func (s Song) GetPatternByOrder(o index.Order) (song.Pattern, error) {
	if int(o) >= len(s.OrderList) {
		return nil, song.ErrStopSong
	}

	pat := s.OrderList[o]
	switch pat {
	case index.InvalidPattern:
		return nil, song.ErrStopSong
	case index.NextPattern:
		return nil, index.ErrNextPattern
	}

	return s.GetPattern(pat)
}

// GetNumChannels returns the number of channels the song has
func (s Song) GetNumChannels() int {
	return s.NumChannels
}

// GetChannelSettings returns the channel settings at index `channelNum`
func (s Song) GetChannelSettings(channelNum index.Channel) song.ChannelSettings {
	return s.ChannelSettings[channelNum]
}

// NumInstruments returns the number of instruments in the song
func (s Song) NumInstruments() int {
	return len(s.Instruments)
}

// IsValidInstrumentID returns true if the instrument exists
func (s Song) IsValidInstrumentID(instNum instrument.ID) bool {
	if instNum.IsEmpty() {
		return false
	}
	switch id := instNum.(type) {
	case channel.InstID:
		iid := int(id)
		return iid > 0 && iid <= len(s.Instruments)
	}
	return false
}

// GetInstrument returns the instrument interface indexed by `instNum` (0-based)
func (s Song) GetInstrument(instID instrument.ID) (instrument.InstrumentIntf, note.Semitone) {
	if instID.IsEmpty() {
		return nil, note.UnchangedSemitone
	}
	switch id := instID.(type) {
	case channel.InstID:
		return s.Instruments[int(id)-1], note.UnchangedSemitone
	}

	return nil, note.UnchangedSemitone
}

// GetName returns the name of the song
func (s Song) GetName() string {
	return s.Head.Name
}

func (s Song) GetPeriodCalculator() song.PeriodCalculatorIntf {
	return s3mPeriod.S3MAmigaConverter
}

func (s Song) GetInitialOrder() index.Order {
	return s.Head.InitialOrder
}

func (s Song) GetRowRenderStringer(row song.Row, channels int, longFormat bool) render.RowStringer {
	nch := min(s.NumChannels, channels)
	rt := render.NewRowText[channel.Data](nch, longFormat)
	rowData := make([]channel.Data, 0, nch)
	pr := row.(Row)
	nprch := min(len(pr), nch)
	for i := 0; i < nprch; i++ {
		if !s.ChannelSettings[i].Enabled || s.ChannelSettings[i].Muted {
			continue
		}
		rowData = append(rowData, pr[i])
	}
	for len(rowData) < nch {
		rowData = append(rowData, channel.Data{})
	}
	rt.Channels = rowData
	return rt
}

func (s Song) GetSystem() system.System {
	return s.System
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

func (s Song) GetMachineSettings() any {
	return s.MS
}
