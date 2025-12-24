package note

import (
	"testing"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	pnote "github.com/gotracker/playback/note"
)

func TestFromItNoteSpecials(t *testing.T) {
	cases := []struct {
		name string
		in   itfile.Note
		exp  pnote.Note
	}{
		{"note off", itfile.Note(255), pnote.ReleaseNote{}},
		{"note cut", itfile.Note(254), pnote.StopNote{}},
		{"note fade", itfile.Note(120), pnote.FadeoutNote{}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromItNote(tt.in); got != tt.exp {
				t.Fatalf("FromItNote(%d) = %#v, want %#v", tt.in, got, tt.exp)
			}
		})
	}
}

func TestFromItNoteNormal(t *testing.T) {
	in := itfile.Note(48)
	exp := pnote.Normal(pnote.Semitone(48))

	if got := FromItNote(in); got != exp {
		t.Fatalf("FromItNote(%d) = %#v, want %#v", in, got, exp)
	}
}
