package instrument

import "github.com/gotracker/playback/note"

type Keyboard[TRemap any] map[note.Semitone]TRemap

func (k Keyboard[TRemap]) GetRemap(st note.Semitone) (TRemap, bool) {
	var empty TRemap
	if k == nil {
		return empty, false
	}

	remap, ok := k[st]
	return remap, ok
}

func (k *Keyboard[TRemap]) SetRemap(st note.Semitone, remap TRemap) {
	k.expect()

	(*k)[st] = remap
}

func (k *Keyboard[TRemap]) expect() {
	if k == nil {
		*k = make(map[note.Semitone]TRemap)
	}
}
