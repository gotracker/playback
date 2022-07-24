package instrument

import "github.com/gotracker/playback/note"

type Keyboard[TRemap any] struct {
	remap map[note.Semitone]TRemap
	Inst  *Instrument
}

func (k Keyboard[TRemap]) GetInstrument() *Instrument {
	return k.Inst
}

func (k Keyboard[TRemap]) GetRemap(st note.Semitone) (TRemap, bool) {
	var empty TRemap
	if k.remap == nil {
		return empty, false
	}

	remap, ok := k.remap[st]
	return remap, ok
}

func (k *Keyboard[TRemap]) SetRemap(st note.Semitone, remap TRemap) {
	k.expect()

	k.remap[st] = remap
}

func (k *Keyboard[TRemap]) expect() {
	if k.remap == nil {
		k.remap = make(map[note.Semitone]TRemap)
	}
}
