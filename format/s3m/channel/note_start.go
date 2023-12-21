package channel

import (
	"fmt"

	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/playback"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

type Start struct {
	semitone note.Semitone
	instId   InstID
}

func start(s note.Semitone, i InstID) EffectS3M {
	if s == 0 && i == 0 {
		return nil
	}

	return Start{
		semitone: 0,
		instId:   i,
	}
}

// Start triggers on the first tick, but before the Tick() function is called
func (e Start) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	mem := cs.GetMemory()
	st := mem.Semitone(e.semitone)
	instID := mem.Inst(e.instId)

	m := p.(S3M)

	inst := cs.GetInstrument()
	prevInst := inst

	var (
		wantRetrigger    bool
		wantRetriggerVol bool
	)
	if instID.IsEmpty() {
		// use current
		inst = prevInst
		wantRetrigger = true
	} else if !m.IsValidInstrumentID(instID) {
		cs.SetTargetInst(nil)
	} else {
		var str note.Semitone
		inst, str = m.GetInstrument(instID)
		if str != note.UnchangedSemitone {
			st = str
		}
		wantRetrigger = true
		wantRetriggerVol = true
	}

	if wantRetrigger {
		var (
			c2Spd period.Frequency
			ft    note.Finetune
		)
		if inst != nil {
			c2Spd = inst.GetC2Spd()
			ft = inst.GetFinetune()
		}
		p := s3mPeriod.CalcSemitonePeriod(st, ft, c2Spd)
		cs.SetTargetPos(sampling.Pos{})
		cs.SetTargetSemitone(st)
		cs.SetTargetPeriod(p)
		cs.SetTargetInst(inst)
	}

	if wantRetriggerVol && inst != nil {
		cs.SetActiveVolume(inst.GetDefaultVolume())
	}

	return nil
}

func (e Start) String() string {
	return fmt.Sprintf("%v %02d", e.semitone, e.instId)
}
