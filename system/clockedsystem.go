package system

import (
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/note"
)

type ClockableSystem interface {
	System
	GetBaseClock() frequency.Frequency
	GetCommonPeriod() uint16
	GetBaseFinetunes() note.Finetune
	GetFinetunesPerOctave() note.Finetune
	GetFinetunesPerSemitone() note.Finetune
	GetSemitonePeriod(note.Key) (uint16, bool)
	GetOctaveShift() uint16
}

type ClockedSystem struct {
	MaxPastNotesPerChannel int

	BaseClock          frequency.Frequency
	BaseFinetunes      note.Finetune
	FinetunesPerOctave note.Finetune
	FinetunesPerNote   note.Finetune
	CommonPeriod       uint16
	CommonRate         frequency.Frequency
	SemitonePeriods    [note.NumKeys]uint16
	OctaveShift        uint16
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

func (s ClockedSystem) GetCommonPeriod() uint16 {
	return s.CommonPeriod
}

func (s ClockedSystem) GetCommonRate() frequency.Frequency {
	return s.CommonRate
}

func (s ClockedSystem) GetSemitonePeriod(k note.Key) (uint16, bool) {
	if int(k) < note.NumKeys {
		return s.SemitonePeriods[int(k)] >> s.OctaveShift, true
	}
	return 0, false
}

func (s ClockedSystem) GetOctaveShift() uint16 {
	return s.OctaveShift
}
