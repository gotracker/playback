package note

import (
	"testing"

	pnote "github.com/gotracker/playback/note"
)

func TestFromXmNoteSpecials(t *testing.T) {
	cases := []struct {
		name string
		in   uint8
		exp  pnote.Note
	}{
		{"release", 97, pnote.ReleaseNote{}},
		{"empty", 0, pnote.EmptyNote{}},
		{"invalid", 98, pnote.InvalidNote{}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromXmNote(tt.in); got != tt.exp {
				t.Fatalf("FromXmNote(%d) = %#v, want %#v", tt.in, got, tt.exp)
			}
		})
	}
}

func TestFromXmNoteNormal(t *testing.T) {
	in := uint8(1)
	exp := pnote.Normal(pnote.Semitone(0))

	if got := FromXmNote(in); got != exp {
		t.Fatalf("FromXmNote(%d) = %#v, want %#v", in, got, exp)
	}
}
