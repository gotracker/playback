package system

import (
	"github.com/gotracker/playback/note"
)

type ClockedSystem struct {
	BaseClock          Frequency
	BaseFinetunes      note.Finetune
	FinetunesPerOctave note.Finetune
}

func (s ClockedSystem) GetBaseClock() Frequency {
	return s.BaseClock
}

func (s ClockedSystem) GetBaseFinetunes() note.Finetune {
	return s.BaseFinetunes
}

func (s ClockedSystem) GetFinetunesPerOctave() note.Finetune {
	return s.FinetunesPerOctave
}
