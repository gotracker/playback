package playback

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/format/it/channel"
	itPeriod "github.com/gotracker/playback/format/it/period"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/state"
)

type doNoteCalc[TPeriod period.Period] struct {
	Semitone   note.Semitone
	UpdateFunc state.PeriodUpdateFunc[TPeriod]
}

func (o doNoteCalc[TPeriod]) Process(p playback.Playback, cs *state.ChannelState[TPeriod, channel.Memory]) error {
	if o.UpdateFunc == nil {
		return nil
	}

	if inst := cs.GetTargetInst(); inst != nil {
		cs.Semitone = note.Semitone(int(o.Semitone) + int(inst.GetSemitoneShift()))
		ft := inst.GetFinetune()
		period := itPeriod.CalcSemitonePeriod[TPeriod](cs.Semitone, ft, inst.GetC2Spd())
		o.UpdateFunc(period)
	}
	return nil
}

func (m *manager[TPeriod]) processEffect(ch int, cs *state.ChannelState[TPeriod, channel.Memory], currentTick int, lastTick bool) error {
	if txn := cs.GetTxn(); txn != nil {
		if err := txn.CommitPreTick(m, cs, currentTick, lastTick, cs.SemitoneSetterFactory); err != nil {
			return err
		}
		if err := txn.CommitTick(m, cs, currentTick, lastTick, cs.SemitoneSetterFactory); err != nil {
			return err
		}
		if err := txn.CommitPostTick(m, cs, currentTick, lastTick, cs.SemitoneSetterFactory); err != nil {
			return err
		}
	}
	cs.SetGlobalVolume(m.GetGlobalVolume())

	if err := m.processRowNote(ch, cs, currentTick, lastTick); err != nil {
		return err
	}

	if err := m.processVoiceUpdates(ch, cs, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (m *manager[TPeriod]) processRowNote(ch int, cs *state.ChannelState[TPeriod, channel.Memory], currentTick int, lastTick bool) error {
	targetTick, noteAction := cs.WillTriggerOn(currentTick)
	if !targetTick {
		return nil
	}

	keyOn := false
	if nc := cs.GetVoice(); nc != nil {
		keyOn = nc.IsKeyOn()
	}

	if noteAction == note.ActionRetrigger {
		cs.TransitionActiveToPastState()
	}

	wantAttack := false
	targetPeriod := cs.GetTargetPeriod()
	if !targetPeriod.IsInvalid() {
		targetInst := cs.GetTargetInst()
		if targetInst != nil {
			keyOn = true
			wantAttack = noteAction == note.ActionRetrigger
		}

		if cs.UseTargetPeriod {
			cs.SetPeriod(targetPeriod)
			cs.SetPortaTargetPeriod(targetPeriod)
		}

		cs.SetInstrument(targetInst)
		cs.SetPos(cs.GetTargetPos())
	}

	var invalidPeriod TPeriod
	if nc := cs.GetVoice(); nc != nil {
		switch noteAction {
		case note.ActionRetrigger:
			if keyOn && wantAttack {
				nc.Attack()
				mem := cs.GetMemory()
				mem.Retrigger()
			}
		case note.ActionRelease:
			nc.Release()
		case note.ActionCut:
			cs.SetInstrument(nil)
			cs.SetPeriod(invalidPeriod)
		}
	}

	return nil
}

func (m *manager[TPeriod]) processVoiceUpdates(ch int, cs *state.ChannelState[TPeriod, channel.Memory], currentTick int, lastTick bool) error {
	if cs.UsePeriodOverride {
		cs.UsePeriodOverride = false
		arpeggioPeriod := cs.GetPeriodOverride()
		cs.SetPeriod(arpeggioPeriod)
	}
	return nil
}

// SetFilterEnable activates or deactivates the amiga low-pass filter on the instruments
func (m *manager[TPeriod]) SetFilterEnable(on bool) {
	for i := range m.song.ChannelSettings {
		c := m.GetChannel(i)
		if o := c.GetRenderChannel(); o != nil {
			if on {
				if o.Filter == nil {
					o.Filter = filter.NewAmigaLPF(period.Frequency(itPeriod.DefaultC2Spd), m.GetSampleRate())
				}
			} else {
				o.Filter = nil
			}
		}
	}
}

// SetTicks sets the number of ticks the row expects to play for
func (m *manager[TPeriod]) SetTicks(ticks int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.Ticks.Set(ticks)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.Ticks.Set(ticks)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// AddRowTicks increases the number of ticks the row expects to play for
func (m *manager[TPeriod]) AddRowTicks(ticks int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.FinePatternDelay.Set(ticks)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.FinePatternDelay.Set(ticks)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// SetPatternDelay sets the repeat number for the row to `rept`
// NOTE: this may be set 1 time (first in wins) and will be reset only by the next row being read in
func (m *manager[TPeriod]) SetPatternDelay(rept int) error {
	if m.preMixRowTxn != nil {
		m.preMixRowTxn.SetPatternDelay(rept)
	} else {
		rowTxn := m.pattern.StartTransaction()
		defer rowTxn.Cancel()

		rowTxn.SetPatternDelay(rept)
		if err := rowTxn.Commit(); err != nil {
			return err
		}
	}

	return nil
}
