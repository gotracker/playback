package layout

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/playback/format/it/channel"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type StringRow[TPeriod period.Period] string

func (r StringRow[TPeriod]) Len() int {
	return len(strings.SplitAfter(string(r), "|")) - 1
}

func (r StringRow[TPeriod]) ForEach(fn func(ch index.Channel, d song.ChannelData[itVolume.Volume]) (bool, error)) error {
	cstrPieces := strings.SplitAfter(string(r), "|")
	cstrPieces = slices.DeleteFunc(cstrPieces, func(s string) bool {
		return len(s) == 0 || s == "|"
	})

	row := make(Row[TPeriod], len(cstrPieces))
	for ch, cstr := range cstrPieces {
		d, err := r.decodeChannel(strings.TrimSuffix(cstr, "|"))
		if err != nil {
			return err
		}
		row[ch] = d
	}

	return row.ForEach(fn)
}

var channelRegex = regexp.MustCompile(`^(...) +(..) +(..) +(...)$`)

func (StringRow[TPeriod]) decodeChannel(cstr string) (channel.Data[TPeriod], error) {
	var d channel.Data[TPeriod]

	pieces := channelRegex.FindStringSubmatch(cstr)
	if len(pieces) != 5 {
		return d, fmt.Errorf("could not parse channel: %q", cstr)
	}
	note, instrument, vol, cmd := pieces[1], pieces[2], pieces[3], pieces[4]

	d.Note = 0

	// note
	if note == "===" || note == "== " {
		d.What |= itfile.ChannelDataFlagNote
		d.Note = 255
	} else if note == "^^^" || note == "^^ " {
		d.What |= itfile.ChannelDataFlagNote
		d.Note = 254
	} else if note == "vvv" || note == "vv " {
		d.What |= itfile.ChannelDataFlagNote
		d.Note = 234
	} else if note != "..." {
		key := note[0:2]
		oct, err := strconv.Atoi(note[2:])
		if err != nil {
			return d, err
		}

		switch key {
		case "C-":
			d.Note = itfile.Note(oct*12 + 0)
		case "C#":
			d.Note = itfile.Note(oct*12 + 1)
		case "D-":
			d.Note = itfile.Note(oct*12 + 2)
		case "D#":
			d.Note = itfile.Note(oct*12 + 3)
		case "E-":
			d.Note = itfile.Note(oct*12 + 4)
		case "F-":
			d.Note = itfile.Note(oct*12 + 5)
		case "F#":
			d.Note = itfile.Note(oct*12 + 6)
		case "G-":
			d.Note = itfile.Note(oct*12 + 7)
		case "G#":
			d.Note = itfile.Note(oct*12 + 8)
		case "A-":
			d.Note = itfile.Note(oct*12 + 9)
		case "A#":
			d.Note = itfile.Note(oct*12 + 10)
		case "B-":
			d.Note = itfile.Note(oct*12 + 11)
		default:
			return d, fmt.Errorf("invalid key in note: %q", note)
		}
		d.What |= itfile.ChannelDataFlagNote
	}

	// instrument
	if instrument != ".." {
		i, err := strconv.ParseUint(strings.TrimSpace(instrument), 16, 8)
		if err != nil {
			return d, err
		}

		if i > 0 {
			d.What |= itfile.ChannelDataFlagInstrument
			d.Instrument = uint8(i)
		}
	}

	// vol
	if vol != ".." {
		v, err := strconv.ParseUint(vol, 16, 8)
		if err != nil {
			return d, err
		}

		d.What |= itfile.ChannelDataFlagVolPan
		d.VolPan = uint8(v)
	}

	// cmd
	if cmd != "..." {
		c := cmd[0]
		i, err := strconv.ParseUint(cmd[1:], 16, 8)
		if err != nil {
			return d, err
		}

		d.What |= itfile.ChannelDataFlagCommand
		if e := c - '@'; e < 26 {
			d.Effect = channel.Command(e)
		}
		d.EffectParameter = channel.DataEffect(i)
	}

	return d, nil
}
