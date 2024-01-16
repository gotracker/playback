package machine

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/memory"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/machine/instruction"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/heucuva/optional"
)

type channel[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	enabled     bool
	memory      song.ChannelMemory
	osc         [NumOscillators]oscillator.Oscillator
	patternLoop struct {
		Start index.Row
		End   index.Row
		Total int
		Count int
	}

	prev struct {
		Period   TPeriod
		Inst     *instrument.Instrument[TMixingVolume, TVolume, TPanning]
		Semitone memory.Value[note.Semitone]
	}
	target struct {
		Inst       *instrument.Instrument[TMixingVolume, TVolume, TPanning]
		Pos        optional.Value[sampling.Pos]
		ActionTick optional.Value[ActionTick]
	}
	portaPeriod TPeriod
	newNote     NewNoteInfo[TPeriod, TMixingVolume, TVolume, TPanning]

	mute          bool
	surround      bool
	filter        filter.Filter
	filterEnabled bool
	nna           note.Action

	cv        voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	pastNotes []index.Channel

	instructions []instruction.Instruction
}

type channelInfo[TPeriod Period, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	Period TPeriod
	Inst   *instrument.Instrument[TMixingVolume, TVolume, TPanning]
}

type channelTargets[TPeriod Period, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	channelInfo[TPeriod, TMixingVolume, TVolume, TPanning]

	Pos        optional.Value[sampling.Pos]
	Action     note.Action
	ActionTick int
}
