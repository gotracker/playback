package machine

import (
	"slices"
	"sort"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) addPastNote(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], pn voice.Voice) {
	// try to find an empty spot in output channels
	type pnVolChan struct {
		vol volume.Volume
		ch  index.Channel
	}
	pnvc := make([]pnVolChan, len(m.virtualOutputs))
	for i := range m.virtualOutputs {
		ch := index.Channel(i)
		rc := &m.virtualOutputs[ch]

		v := rc.GetVoice()
		if v == nil {
			rc.StartVoice(pn, func() {
				c.removePastNote(m, ch)
			})
			c.pastNotes = append(c.pastNotes, ch)
			return
		}

		vc := pnVolChan{
			vol: volume.Volume(1),
			ch:  ch,
		}
		if ampMod, ok := v.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok {
			vc.vol = ampMod.GetFinalVolume()
		}

		pnvc = append(pnvc, vc)
	}

	// we failed to find a spot?
	if len(pnvc) == 0 {
		// no room at all? strange
		pn.Stop()
		return
	}

	// look for lowest volume
	sort.SliceStable(pnvc, func(i, j int) bool {
		return pnvc[i].vol < pnvc[j].vol
	})

	lowest := pnvc[0].ch
	rc := &m.virtualOutputs[lowest]

	rc.StartVoice(pn, func() {
		c.removePastNote(m, lowest)
	})
	c.pastNotes = append(c.pastNotes, lowest)
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) removePastNote(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], ch index.Channel) {
	c.pastNotes = slices.DeleteFunc(c.pastNotes, func(i index.Channel) bool {
		return i == ch
	})
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPastNoteAction(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], na note.Action) {
	for _, ch := range c.pastNotes {
		rc := &m.virtualOutputs[ch]
		v := rc.GetVoice()
		if v == nil {
			continue
		}

		switch na {
		case note.ActionCut:
			rc.StopVoice()
		case note.ActionRelease:
			v.Release()
		case note.ActionFadeout:
			v.Release()
			v.Fadeout()
		case note.ActionRetrigger:
			v.Release()
			v.Attack()

		case note.ActionContinue:
			fallthrough
		default:
			// nothing
		}
	}
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) updatePastNotes(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) {
	var updated []index.Channel
	for _, ch := range c.pastNotes {
		rc := &m.virtualOutputs[ch]
		v := rc.GetVoice()
		if v == nil {
			continue
		}

		if v.IsDone() {
			rc.StopVoice()
			continue
		}

		updated = append(updated, ch)
	}
	c.pastNotes = updated
}
