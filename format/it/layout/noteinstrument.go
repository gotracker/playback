package layout

import (
	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// NoteInstrument is the note remapping and instrument pair
type NoteInstrument[TPeriod period.Period] struct {
	NoteRemap note.Semitone
	Inst      *instrument.Instrument[TPeriod, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]
}
