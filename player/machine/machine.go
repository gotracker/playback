package machine

import (
	"errors"
	"fmt"

	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/instruction"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice/opl2"
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/gotracker/playback/voice/types"
)

type (
	Period  = settings.Period
	Volume  = settings.Volume
	Panning = settings.Panning
)

type VolumeFMA[TVolume Volume] interface {
	FMA(multiplier, add float32) TVolume
}

type PanningFMA[TPanning Panning] interface {
	FMA(multiplier, add float32) TPanning
}

type MachineInfo interface {
	GetNumOrders() int
	CanOrderLoop() bool
	GetName() string
}

type MachineTicker interface {
	MachineInfo

	Generate(s *sampler.Sampler) (*output.PremixData, error)
	Tick() error
}

type Machine[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	MachineTicker

	ConvertToPeriod(n note.Note) TPeriod
	IgnoreUnknownEffect() bool
	GetQuirks() *settings.MachineQuirks

	// Globals
	SetTempo(tempo int) error
	SetBPM(bpm int) error
	SlideBPM(add int) error
	SetGlobalVolume(v TGlobalVolume) error
	SlideGlobalVolume(multiplier, add float32) error
	SetMixingVolume(v volume.Volume) error
	SetSynthVolume(v volume.Volume) error
	SetSampleVolume(v volume.Volume) error
	SetOrder(o index.Order) error
	SetRow(r index.Row, breakOrder bool) error
	SetFilterOnAllChannelsByFilterName(name string, enabled bool) error
	GetPosition() Position

	// Single Row
	AddExtraTicks(ticks int) error
	RowRepeat(times int) error

	// Channel
	GetChannelMemory(ch index.Channel) (song.ChannelMemory, error)
	IsChannelMuted(ch index.Channel) (bool, error)
	SetChannelMute(ch index.Channel, muted bool) error
	SetChannelMixingVolume(ch index.Channel, v TMixingVolume) error
	GetChannelPeriod(ch index.Channel) (TPeriod, error)
	SetChannelPeriod(ch index.Channel, p TPeriod) error
	SetChannelPeriodDelta(ch index.Channel, d period.Delta) error
	GetChannelInstrument(ch index.Channel) (*instrument.Instrument[TPeriod, TMixingVolume, TVolume, TPanning], error)
	SetChannelInstrumentByID(ch index.Channel, i instrument.ID) error
	GetChannelVolume(ch index.Channel) (TVolume, error)
	SetChannelVolume(ch index.Channel, v TVolume) error
	SetChannelVolumeDelta(ch index.Channel, d types.VolumeDelta) error
	GetChannelPan(ch index.Channel) (TPanning, error)
	SetChannelPan(ch index.Channel, pan TPanning) error
	SetChannelPanningDelta(ch index.Channel, d types.PanDelta) error
	SetChannelSurround(ch index.Channel, enabled bool) error
	SetChannelFilter(ch index.Channel, f filter.Filter) error
	ChannelStopOrRelease(ch index.Channel) error
	ChannelStop(ch index.Channel) error
	ChannelRelease(ch index.Channel) error
	ChannelFadeout(ch index.Channel) error
	GetNextChannelWavetableValue(ch index.Channel, speed int, depth float32, oscSelect Oscillator) (float32, error)
	SetChannelNoteAction(ch index.Channel, na note.Action, tick int) error
	SetPatternLoopStart(ch index.Channel) error
	SetPatternLoops(ch index.Channel, count int) error
	StartChannelPortaToNote(ch index.Channel) error
	DoChannelPortaToNote(ch index.Channel, delta period.Delta) error
	DoChannelPortaDown(ch index.Channel, delta period.Delta) error
	DoChannelPortaUp(ch index.Channel, delta period.Delta) error
	DoChannelArpeggio(ch index.Channel, delta int8) error
	SlideChannelVolume(ch index.Channel, multiplier, add float32) error
	SlideChannelMixingVolume(ch index.Channel, multiplier, add float32) error
	SetChannelPos(ch index.Channel, pos sampling.Pos) error
	SetChannelEnvelopePositions(ch index.Channel, pos int) error
	SlideChannelPan(ch index.Channel, multiplier, add float32) error
	SetChannelVolumeActive(ch index.Channel, on bool) error
	SetChannelOscillatorWaveform(ch index.Channel, osc Oscillator, wave oscillator.WaveTableSelect) error
	DoChannelPastNoteEffect(ch index.Channel, na note.Action) error
	SetChannelNewNoteAction(ch index.Channel, na note.Action) error
	SetChannelVolumeEnvelopeEnable(ch index.Channel, enabled bool) error
	SetChannelPanningEnvelopeEnable(ch index.Channel, enabled bool) error
	SetChannelPitchEnvelopeEnable(ch index.Channel, enabled bool) error

	// Instructions
	DoInstructionOrderStart(ch index.Channel, i instruction.Instruction) error
	DoInstructionRowStart(ch index.Channel, i instruction.Instruction) error
	DoInstructionTick(ch index.Channel, i instruction.Instruction) error
	DoInstructionRowEnd(ch index.Channel, i instruction.Instruction) error
	DoInstructionOrderEnd(ch index.Channel, i instruction.Instruction) error
}

type machine[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	globals[TGlobalVolume]
	singleRow
	channels  []channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	pastNotes []pastNote[TPeriod]

	ticker ticker
	age    int

	songData    song.Data
	ms          *settings.MachineSettings[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	us          settings.UserSettings
	opl2        opl2.Chip
	opl2Enabled bool

	rowStringer render.RowStringer
	// 1:1 with channels
	actualOutputs []render.Channel[TPeriod]
	// extra channels for pastNotes playback
	virtualOutputs []render.Channel[TPeriod]
}

func NewMachine(songData song.Data, us settings.UserSettings) (MachineTicker, error) {
	if songData == nil {
		return nil, errors.New("songData is nil")
	}

	tl := typeLookup{
		p:   songData.GetPeriodType(),
		gv:  songData.GetGlobalVolumeType(),
		cmv: songData.GetChannelMixingVolumeType(),
		cv:  songData.GetChannelVolumeType(),
		cp:  songData.GetChannelPanningType(),
	}

	factory, found := factoryRegistry[tl]
	str := tl.String()
	_ = str
	if !found {
		return nil, errors.New("could not identify machine type from song parameters")
	}

	return factory(songData, us)
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetPosition() Position {
	return m.ticker.current
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetQuirks() *settings.MachineQuirks {
	return &m.ms.Quirks
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) ConvertToPeriod(n note.Note) TPeriod {
	return m.ms.PeriodConverter.GetPeriod(n)
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) IgnoreUnknownEffect() bool {
	return m.us.IgnoreUnknownEffect
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetNumOrders() int {
	return len(m.songData.GetOrderList())
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) CanOrderLoop() bool {
	return m.us.SongLoopCount != 0
}

func (m machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetName() string {
	return m.songData.GetName()
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) getChannel(ch index.Channel) (*channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], error) {
	if int(ch) >= len(m.channels) {
		return nil, fmt.Errorf("invalid channel index: %d", ch)
	}

	return &m.channels[ch], nil
}

type dataInstructionGenerator[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	ToInstructions(m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel, songData song.Data) ([]instruction.Instruction, error)
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) getRowData() (song.Row, error) {
	pat, err := m.songData.GetPatternByOrder(m.ticker.current.Order)
	if err != nil {
		return nil, err
	}
	if pat == nil || int(m.ticker.current.Row) >= pat.NumRows() {
		return nil, song.ErrStopSong
	}

	row := pat.GetRow(m.ticker.current.Row)
	if row == nil {
		return nil, song.ErrStopSong
	}

	return row, nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) updateInstructions(rowData song.Row) error {
	for i := range m.channels {
		m.channels[i].instructions = nil
	}

	numRowChannels := song.GetRowNumChannels[TVolume](rowData)
	rowChannels := min(m.songData.GetNumChannels(), numRowChannels)
	return song.ForEachRowChannel(rowData, func(ch index.Channel, d song.ChannelData[TVolume]) (bool, error) {
		if int(ch) >= rowChannels {
			return true, nil
		}

		c := &m.channels[ch]
		if !c.enabled || d == nil || (m.ms.Quirks.DoNotProcessEffectsOnMutedChannels && c.cv.IsMuted()) {
			return true, nil
		}

		if err := c.decodeNote(m, d); err != nil {
			return false, fmt.Errorf("channel[%d] decode error: %w", ch, err)
		}

		if gen, ok := d.(dataInstructionGenerator[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]); ok {
			insts, err := gen.ToInstructions(m, ch, m.songData)
			if err != nil {
				return false, fmt.Errorf("channel[%d] instruction error: %w", ch, err)
			}

			c.instructions = insts
		}
		return true, nil
	})
}

func GetPeriodCalculator[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning](m Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) song.PeriodCalculator[TPeriod] {
	mach, _ := m.(*machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning])
	if mach == nil {
		return nil
	}

	return mach.ms.PeriodConverter
}
