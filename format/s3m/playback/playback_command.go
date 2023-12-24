package playback

import (
	"github.com/gotracker/playback"
	"github.com/gotracker/playback/filter"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/state"
)

type doNoteCalc struct {
	Semitone   note.Semitone
	UpdateFunc state.PeriodUpdateFunc[period.Amiga]
}

func (o doNoteCalc) Process(p playback.Playback, cs *channelState) error {
	if o.UpdateFunc == nil {
		return nil
	}

	if inst := cs.GetTargetState().Instrument; inst != nil {
		cs.Semitone = note.Semitone(int(o.Semitone) + int(inst.GetSemitoneShift()))
		period := s3mPeriod.CalcSemitonePeriod(cs.Semitone, inst.GetFinetune(), inst.GetSampleRate())
		o.UpdateFunc(period)
	}
	return nil
}

func (m *manager) processEffect(ch int, cs *channelState, currentTick int, lastTick bool) error {
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

	if err := m.processRowNote(ch, cs, currentTick, lastTick); err != nil {
		return err
	}

	if err := m.processVoiceUpdates(ch, cs, currentTick, lastTick); err != nil {
		return err
	}

	return nil
}

func (m *manager) processRowNote(ch int, cs *channelState, currentTick int, lastTick bool) error {
	triggerTick, noteAction := cs.WillTriggerOn(currentTick)
	if !triggerTick {
		return nil
	}
	n := cs.GetChannelData().GetNote()
	keyOn := false
	keyOff := false
	stop := false

	active := cs.GetActiveState()
	target := cs.GetTargetState()

	if targetInst := target.Instrument; targetInst != nil {
		cs.SetInstrument(targetInst)
		keyOn = true
	} else {
		cs.SetInstrument(nil)
	}

	if cs.UseTargetPeriod {
		if nc := cs.GetVoice(); nc != nil {
			nc.Release()
			nc.Fadeout()
		}
		targetPeriod := target.Period
		active.Period = targetPeriod
		cs.SetPortaTargetPeriod(targetPeriod)
	}
	active.Pos = target.Pos

	if active.Instrument != nil {
		keyOff = active.Instrument.IsReleaseNote(n)
		stop = active.Instrument.IsStopNote(n)
	}

	if nc := cs.GetVoice(); nc != nil {
		if keyOn && noteAction == note.ActionRetrigger {
			// S3M is weird and only sets the global volume on the channel when a KeyOn happens
			cs.SetGlobalVolume(m.GetGlobalVolume())
			nc.Attack()
			mem := cs.GetMemory()
			mem.Retrigger()
		} else if keyOff {
			nc.Release()
			nc.Fadeout()
		} else if stop {
			cs.SetInstrument(nil)
			active.NoteCut()
		}
	}
	return nil
}

func (m *manager) processVoiceUpdates(ch int, cs *channelState, currentTick int, lastTick bool) error {
	if cs.UsePeriodOverride {
		cs.UsePeriodOverride = false
		arpeggioPeriod := cs.GetPeriodOverride()
		cs.GetActiveState().Period = arpeggioPeriod
	}
	return nil
}

// SetFilterEnable activates or deactivates the amiga low-pass filter on the instruments
func (m *manager) SetFilterEnable(on bool) {
	for i := range m.song.ChannelSettings {
		c := m.GetChannel(i)
		if o := c.GetRenderChannel(); o != nil {
			if on {
				if o.Filter == nil {
					o.Filter = filter.NewAmigaLPF(s3mPeriod.DefaultC4SampleRate, m.GetSampleRate())
				}
			} else {
				o.Filter = nil
			}
		}
	}
}

// SetTicks sets the number of ticks the row expects to play for
func (m *manager) SetTicks(ticks int) error {
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
func (m *manager) AddRowTicks(ticks int) error {
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
func (m *manager) SetPatternDelay(rept int) error {
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
