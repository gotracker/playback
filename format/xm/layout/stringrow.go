package layout

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	xmfile "github.com/gotracker/goaudiofile/music/tracked/xm"
	"github.com/gotracker/playback/format/xm/channel"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type StringRow[TPeriod period.Period] string

func (r StringRow[TPeriod]) Len() int {
	return len(strings.SplitAfter(string(r), "|")) - 1
}

func (r StringRow[TPeriod]) ForEach(fn func(ch index.Channel, d song.ChannelData[xmVolume.XmVolume]) (bool, error)) error {
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
		d.What |= xmfile.ChannelFlagHasNote
		d.Note = 97
	} else if note != "..." {
		key := note[0:2]
		oct, err := strconv.Atoi(note[2:])
		if err != nil {
			return d, err
		}

		switch key {
		case "C-":
			d.Note = uint8(oct*12 + 1)
		case "C#":
			d.Note = uint8(oct*12 + 2)
		case "D-":
			d.Note = uint8(oct*12 + 3)
		case "D#":
			d.Note = uint8(oct*12 + 4)
		case "E-":
			d.Note = uint8(oct*12 + 5)
		case "F-":
			d.Note = uint8(oct*12 + 6)
		case "F#":
			d.Note = uint8(oct*12 + 7)
		case "G-":
			d.Note = uint8(oct*12 + 8)
		case "G#":
			d.Note = uint8(oct*12 + 9)
		case "A-":
			d.Note = uint8(oct*12 + 10)
		case "A#":
			d.Note = uint8(oct*12 + 11)
		case "B-":
			d.Note = uint8(oct*12 + 12)
		default:
			return d, fmt.Errorf("invalid key in note: %q", note)
		}
		d.What |= xmfile.ChannelFlagHasNote
	}

	// instrument
	if instrument != "  " {
		i, err := strconv.ParseUint(strings.TrimSpace(instrument), 16, 8)
		if err != nil {
			return d, err
		}

		if i > 0 {
			d.What |= xmfile.ChannelFlagHasInstrument
			d.Instrument = uint8(i)
		}
	}

	// vol
	if vol != ".." {
		v, err := strconv.ParseUint(vol, 16, 8)
		if err != nil {
			return d, err
		}

		d.What |= xmfile.ChannelFlagHasVolume
		d.Volume = xmVolume.VolEffect(v)
	}

	// cmd
	if cmd != "..." {
		c := cmd[0]
		i, err := strconv.ParseUint(cmd[1:], 16, 8)
		if err != nil {
			return d, err
		}

		d.What |= xmfile.ChannelFlagHasEffect | xmfile.ChannelFlagHasEffectParameter
		if e := c - '0'; e <= 9 {
			d.Effect = channel.Command(e)
		} else if e := c - 'A' + 10; e < 36 {
			d.Effect = channel.Command(e)
		}
		d.EffectParameter = channel.DataEffect(i)
	}

	return d, nil
}
