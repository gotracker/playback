package system

import "github.com/gotracker/playback/note"

type System interface {
	GetBaseClock() Frequency
	GetBaseFinetunes() note.Finetune
	GetFinetunesPerOctave() note.Finetune
}
