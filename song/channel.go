package song

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
)

// ChannelData is an interface for channel data
type ChannelData interface {
	HasNote() bool
	GetNote() note.Note

	HasInstrument() bool
	GetInstrument(note.Semitone) instrument.ID

	HasVolume() bool
	GetVolume() volume.Volume

	HasCommand() bool

	Channel() uint8

	fmt.Stringer
	ShortString() string
}
