package layout

import (
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

// NoteInstrument is the note remapping and instrument pair
type NoteInstrument struct {
	NoteRemap note.Semitone
	Inst      *instrument.Instrument
}
