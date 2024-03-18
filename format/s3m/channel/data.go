package channel

import (
	"fmt"
	"strings"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/mixing/volume"

	"github.com/gotracker/playback"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/player/machine/instruction"
	"github.com/gotracker/playback/song"
)

// DataEffect is the type of a channel's EffectParameter value
type DataEffect uint8

// Data is the data for the channel
type Data struct {
	What       s3mfile.PatternFlags
	Note       s3mfile.Note
	Instrument uint8
	Volume     s3mVolume.Volume
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
func (d Data) GetInstrument() int {
	return int(d.Instrument)
}

// HasVolume returns true if there exists a volume on the channel
func (d Data) HasVolume() bool {
	return d.What.HasVolume()
}

func (d Data) GetVolumeGeneric() volume.Volume {
	return d.Volume.ToVolume()
}

// GetVolume returns the volume for the channel
func (d Data) GetVolume() s3mVolume.Volume {
	return d.Volume
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

func (d Data) ToInstructions(m machine.Machine[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning], ch index.Channel, songData song.Data) ([]instruction.Instruction, error) {
	var instructions []instruction.Instruction

	mem, err := machine.GetChannelMemory[*Memory](m, ch)
	if err != nil {
		return nil, err
	}

	if e := EffectFactory(mem, d); e != nil {
		instructions = append(instructions, e)
	}

	return instructions, nil
}
