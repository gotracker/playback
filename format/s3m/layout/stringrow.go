package layout

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"
	"github.com/gotracker/playback/format/s3m/channel"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/song"
)

type StringRow string

func (r StringRow) Len() int {
	return len(strings.SplitAfter(string(r), "|")) - 1
}

func (r StringRow) ForEach(fn func(ch index.Channel, d song.ChannelData[s3mVolume.Volume]) (bool, error)) error {
	cstrPieces := strings.SplitAfter(string(r), "|")
	cstrPieces = slices.DeleteFunc(cstrPieces, func(s string) bool {
		return len(s) == 0 || s == "|"
	})
	for ch, cstr := range cstrPieces {
		d, err := r.decodeChannel(strings.TrimSuffix(cstr, "|"))
		if err != nil {
			return err
		}
		cont, err := fn(index.Channel(ch), d)
		if err != nil {
			return err
		}

		if !cont {
			break
		}
	}
	return nil
}

var channelRegex = regexp.MustCompile(`^(...) +(..) +(..) +(...)$`)

func (StringRow) decodeChannel(cstr string) (channel.Data, error) {
	var d channel.Data

	pieces := channelRegex.FindStringSubmatch(cstr)
	if len(pieces) != 5 {
		return d, fmt.Errorf("could not parse channel: %q", cstr)
	}
	note, instrument, vol, cmd := pieces[1], pieces[2], pieces[3], pieces[4]

	d.Note = s3mfile.EmptyNote

	// note
	if note == "^^." {
		d.What |= s3mfile.PatternFlagNote
		d.Note = s3mfile.StopNote
	} else if note != "..." {
		key := note[0:2]
		oct, err := strconv.Atoi(note[2:])
		if err != nil {
			return d, err
		}

		switch key {
		case "C-":
			d.Note = s3mfile.Note(oct<<4 | 0)
		case "C#":
			d.Note = s3mfile.Note(oct<<4 | 1)
		case "D-":
			d.Note = s3mfile.Note(oct<<4 | 2)
		case "D#":
			d.Note = s3mfile.Note(oct<<4 | 3)
		case "E-":
			d.Note = s3mfile.Note(oct<<4 | 4)
		case "F-":
			d.Note = s3mfile.Note(oct<<4 | 5)
		case "F#":
			d.Note = s3mfile.Note(oct<<4 | 6)
		case "G-":
			d.Note = s3mfile.Note(oct<<4 | 7)
		case "G#":
			d.Note = s3mfile.Note(oct<<4 | 8)
		case "A-":
			d.Note = s3mfile.Note(oct<<4 | 9)
		case "A#":
			d.Note = s3mfile.Note(oct<<4 | 10)
		case "B-":
			d.Note = s3mfile.Note(oct<<4 | 11)
		default:
			return d, fmt.Errorf("invalid key in note: %q", note)
		}
		d.What |= s3mfile.PatternFlagNote
	}

	// instrument
	if instrument != ".." {
		i, err := strconv.Atoi(instrument)
		if err != nil {
			return d, err
		}

		if i > 0 {
			d.What |= s3mfile.PatternFlagNote
			d.Instrument = channel.InstID(i)
		}
	}

	// vol
	if vol != ".." {
		v, err := strconv.Atoi(vol)
		if err != nil {
			return d, err
		}

		d.What |= s3mfile.PatternFlagVolume
		d.Volume = s3mVolume.Volume(v)
	}

	// cmd
	if cmd != "..." {
		c := cmd[0]
		i, err := strconv.ParseUint(cmd[1:], 16, 8)
		if err != nil {
			return d, err
		}

		d.What |= s3mfile.PatternFlagCommand
		d.Command = c - '@'
		d.Info = channel.DataEffect(i)
	}

	return d, nil
}
