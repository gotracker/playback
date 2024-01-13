package system

import (
	"github.com/gotracker/playback/note"
)

type ClockableSystem interface {
	System
	GetBaseClock() Frequency
	GetBaseFinetunes() note.Finetune
	GetFinetunesPerOctave() note.Finetune
	GetFinetunesPerSemitone() note.Finetune
	GetSemitonePeriod(note.Key) (float32, bool)
	GetSamplerSpeed(sampleRate Frequency) float32
}

type ClockedSystem struct {
	MaxPastNotesPerChannel int

	BaseClock          Frequency
	BaseFinetunes      note.Finetune
	FinetunesPerOctave note.Finetune
	FinetunesPerNote   note.Finetune
	CommonRate         Frequency
	SemitonePeriods    [note.NumKeys]float32
}

var _ ClockableSystem = (*ClockedSystem)(nil)

func (s ClockedSystem) GetMaxPastNotesPerChannel() int {
	return s.MaxPastNotesPerChannel
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

func (s ClockedSystem) GetFinetunesPerSemitone() note.Finetune {
	return s.FinetunesPerNote
}

func (s ClockedSystem) GetCommonRate() Frequency {
	return s.CommonRate
}

func (s ClockedSystem) GetSamplerSpeed(sampleRate Frequency) float32 {
	if sampleRate == 0 {
		panic("sampleRate is 0")
	}
	return float32(s.BaseClock) / float32(sampleRate)
}

func (s ClockedSystem) GetSemitonePeriod(k note.Key) (float32, bool) {
	if int(k) < note.NumKeys {
		return s.SemitonePeriods[int(k)], true
	}
	return 0, false
}
