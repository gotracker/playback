package channel

import (
	"fmt"
	"strings"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

// DataEffect is the type of a channel's EffectParameter value
type DataEffect uint8

// Data is the data for the channel
type Data struct {
	What       s3mfile.PatternFlags
	Note       s3mfile.Note
	Instrument InstID
	Volume     s3mfile.Volume
	Command    uint8
	Info       DataEffect
}

// HasNote returns true if there exists a note on the channel
func (d Data) HasNote() bool {
	return d.What.HasNote()
}

// GetNote returns the note for the channel
func (d Data) GetNote() note.Note {
	return NoteFromS3MNote(d.Note)
}

// HasInstrument returns true if there exists an instrument on the channel
func (d Data) HasInstrument() bool {
	return d.Instrument != 0
}

// GetInstrument returns the instrument for the channel
func (d Data) GetInstrument(stmem note.Semitone) instrument.ID {
	return d.Instrument
}

// HasVolume returns true if there exists a volume on the channel
func (d Data) HasVolume() bool {
	return d.What.HasVolume()
}

// GetVolume returns the volume for the channel
func (d Data) GetVolume() volume.Volume {
	return s3mVolume.VolumeFromS3M(d.Volume)
}

// HasCommand returns true if there exists a command on the channel
func (d Data) HasCommand() bool {
	return d.What.HasCommand()
}

// Channel returns the channel ID for the channel
func (d Data) Channel() uint8 {
	return d.What.Channel()
}

func (d Data) GetEffects(mem *Memory, p period.Period) []playback.Effect {
	if d.HasCommand() {
		if e := EffectFactory(mem, d); e != nil {
			return []playback.Effect{e}
		}
	}

	return nil
}

func (d Data) String() string {
	pieces := []string{
		"...", // note
		"..",  // inst
		"..",  // vol
		"...", // eff
	}
	if d.HasNote() {
		pieces[0] = d.GetNote().String()
	}
	if d.HasInstrument() {
		pieces[1] = fmt.Sprintf("%02d", d.Instrument)
	}
	if d.HasVolume() {
		pieces[2] = fmt.Sprintf("%02d", d.Volume)
	}
	if d.HasCommand() {
		pieces[3] = fmt.Sprintf("%c%02X", d.Command+'@', d.Info)
	}
	return strings.Join(pieces, " ")
}

func (d Data) ShortString() string {
	if d.HasNote() {
		return d.GetNote().String()
	}
	return "..."
}

// NoteFromS3MNote converts an S3M file note into a player note
func NoteFromS3MNote(sn s3mfile.Note) note.Note {
	switch {
	case sn == s3mfile.EmptyNote:
		return note.EmptyNote{}
	case sn == s3mfile.StopNote:
		return note.StopOrReleaseNote{}
	default:
		k := uint8(sn.Key()) & 0x0f
		o := uint8(sn.Octave()) & 0x0f
		if k < 12 && o < 10 {
			s := note.Semitone(o*12 + k)
			return note.Normal(s)
		}
	}
	return note.InvalidNote{}
}
