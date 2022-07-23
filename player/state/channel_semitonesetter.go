package state

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

type PeriodUpdateFunc func(period.Period)

type SemitoneSetterFactory[TChannelState any] func(note.Semitone, PeriodUpdateFunc) NoteOp[TChannelState]

type ChannelStateSemitoneSetter[TChannelState any] struct {
	SemitoneSetterFactory SemitoneSetterFactory[TChannelState]
}
