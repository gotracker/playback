package render

import (
	"testing"

	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
)

type stubChannelData struct{}

func (stubChannelData) HasNote() bool                   { return false }
func (stubChannelData) GetNote() note.Note              { return note.EmptyNote{} }
func (stubChannelData) HasInstrument() bool             { return false }
func (stubChannelData) GetInstrument() int              { return 0 }
func (stubChannelData) HasVolume() bool                 { return false }
func (stubChannelData) GetVolumeGeneric() volume.Volume { return 0 }
func (stubChannelData) HasCommand() bool                { return false }
func (stubChannelData) Channel() uint8                  { return 0 }
func (stubChannelData) String() string                  { return "AAA" }
func (stubChannelData) ShortString() string             { return "A" }

func TestRowTextShortAndLong(t *testing.T) {
	vm := NewRowViewModel[stubChannelData](3)
	vm.Channels[0] = stubChannelData{}
	vm.Channels[1] = stubChannelData{}
	vm.Channels[2] = stubChannelData{}

	short := FormatRowText(vm, false).String()
	if short != "|A|A|A|" {
		t.Fatalf("unexpected short row text: %q", short)
	}

	long := FormatRowText(vm, true).String()
	if long != "|AAA|AAA|AAA|" {
		t.Fatalf("unexpected long row text: %q", long)
	}

	vm.MaxChannels = 2
	truncated := FormatRowText(vm, false).String()
	if truncated != "|A|A|" {
		t.Fatalf("expected truncation to 2 channels, got %q", truncated)
	}
}
