package system

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
)

type ClockableSystem interface {
	System
	GetBaseClock() frequency.Frequency
	GetCommonPeriod() float64
	GetBaseFinetunes() note.Finetune
	GetFinetunesPerOctave() note.Finetune
	GetFinetunesPerSemitone() note.Finetune
	GetSemitonePeriod(note.Key) (float32, bool)
}

type ClockedSystem struct {
	MaxPastNotesPerChannel int

	BaseClock          frequency.Frequency
	BaseFinetunes      note.Finetune
	FinetunesPerOctave note.Finetune
	FinetunesPerNote   note.Finetune
	CommonPeriod       float64
	CommonRate         frequency.Frequency
	SemitonePeriods    [note.NumKeys]float32
}

var _ ClockableSystem = (*ClockedSystem)(nil)

func (s ClockedSystem) GetMaxPastNotesPerChannel() int {
	return s.MaxPastNotesPerChannel
}

func (s ClockedSystem) GetBaseClock() frequency.Frequency {
	return s.BaseClock
}

func (s ClockedSystem) GetBaseFinetunes() note.Finetune {
	return s.BaseFinetunes
}

func (s ClockedSystem) GetFinetunesPerOctave() note.Finetune {
	return s.FinetunesPerOctave
}

func (s ClockedSystem) GetFinetunesPerSemitone() note.Finetune {
	return s.FinetunesPerNote
}

func (s ClockedSystem) GetCommonPeriod() float64 {
	return s.CommonPeriod
}

func (s ClockedSystem) GetCommonRate() frequency.Frequency {
	return s.CommonRate
}

func (s ClockedSystem) GetSemitonePeriod(k note.Key) (float32, bool) {
	if int(k) < note.NumKeys {
		return s.SemitonePeriods[int(k)], true
	}
	return 0, false
}
