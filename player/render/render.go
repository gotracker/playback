package render

import (
	"strings"

	"github.com/gotracker/playback/song"
)

// RowViewModel holds the channel data to be rendered.
type RowViewModel[TChannelData song.ChannelDataIntf] struct {
	Channels    []TChannelData
	MaxChannels int // <=0 means no limit
}

// NewRowViewModel creates an empty row view model with the requested channel capacity.
func NewRowViewModel[TChannelData song.ChannelDataIntf](channels int) RowViewModel[TChannelData] {
	return RowViewModel[TChannelData]{
		Channels:    make([]TChannelData, channels),
		MaxChannels: -1,
	}
}

// RowText formats a row view model as a string.
type RowText[TChannelData song.ChannelDataIntf] struct {
	ViewModel  RowViewModel[TChannelData]
	longFormat bool
}

// FormatRowText builds a stringer from a populated view model.
func FormatRowText[TChannelData song.ChannelDataIntf](vm RowViewModel[TChannelData], longFormat bool) RowText[TChannelData] {
	return RowText[TChannelData]{
		ViewModel:  vm,
		longFormat: longFormat,
	}
}

func (rt RowText[TChannelData]) String() string {
	vm := rt.ViewModel
	maxChannels := vm.MaxChannels
	if maxChannels <= 0 {
		maxChannels = len(vm.Channels)
	}
	items := make([]string, 0, len(vm.Channels))
	for i, c := range vm.Channels {
		if maxChannels >= 0 && i >= maxChannels {
			break
		}
		if rt.longFormat {
			items = append(items, c.String())
		} else {
			items = append(items, c.ShortString())
		}
	}
	return "|" + strings.Join(items, "|") + "|"
}

type RowStringer = song.RowStringer

// RowRender is the final output of a single row's data
type RowRender struct {
	Order   int
	Row     int
	Tick    int
	RowText RowStringer
}
