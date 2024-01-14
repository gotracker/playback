package machine

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/voice"
)

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) addPastNote(m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], pn voice.Voice) {
	mpnpc := m.songData.GetSystem().GetMaxPastNotesPerChannel()

	// make room
	if over := len(c.pastNotes) + 1 - mpnpc; over > 0 {
		for _, rc := range c.pastNotes[0:over] {
			rc.StopVoice()
		}

		c.pastNotes = c.pastNotes[over:]
	}

	for i := range m.virtualOutputs {
		rc := &m.virtualOutputs[i]

		if rc.Voice == nil {
			rc.Voice = pn
			c.pastNotes = append(c.pastNotes, rc)
			return
		}
	}

	// we failed to find a spot?
	pn.Stop()
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPastNoteAction(na note.Action) {
	for _, pn := range c.pastNotes {
		if pn.Voice == nil {
			continue
		}

		switch na {
		case note.ActionCut:
			pn.StopVoice()
		case note.ActionRelease:
			pn.Voice.Release()
		case note.ActionFadeout:
			pn.Voice.Release()
			pn.Voice.Fadeout()
		case note.ActionRetrigger:
			pn.Voice.Release()
			pn.Voice.Attack()

		case note.ActionContinue:
			fallthrough
		default:
			// nothing
		}
	}
}

func (c *channel[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) updatePastNotes() {
	var updated []*render.Channel[TPeriod]
	for _, pn := range c.pastNotes {
		if pn.Voice == nil {
			continue
		}

		if pn.Voice.IsDone() {
			pn.StopVoice()
			continue
		}

		updated = append(updated, pn)
	}
	c.pastNotes = updated
}
