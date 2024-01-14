package machine

import (
	"sort"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/voice"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) canPastNote() bool {
	return m.songData.GetSystem().GetMaxPastNotesPerChannel() > 0
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) addPastNote(ch index.Channel, pn voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) {
	type pastNoteAges struct {
		holder int
		entry  pastNote[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	}

	// first pass, try to add it to the channel listed
	maxNotes := len(m.channels) * m.songData.GetSystem().GetMaxPastNotesPerChannel()
	if len(m.pastNotes) < maxNotes {
		c := &m.channels[ch]

		m.pastNotes = append(m.pastNotes, ch)
		c.pn.AddPastNote(ch, pn, m.age)
		return
	}

	ages := make([]pastNoteAges, 0, len(m.outputChannels))
	// second pass, try to bump the oldest
	for i := range m.channels {
		c := &m.channels[i]

		if oldest := c.pn.GetOldest(); oldest != nil {
			ages = append(ages, pastNoteAges{
				holder: i,
				entry:  *oldest,
			})
		}
	}

	if len(ages) == 0 {
		// impossible, but just in case
		pn.Stop()
		return
	}

	// find the oldest
	sort.Slice(ages, func(i, j int) bool {
		return ages[i].entry.age < ages[j].entry.age
	})

	// jam it in
	oldestHolder := ages[0].holder
	m.channels[oldestHolder].pn.AddPastNote(ch, pn, m.age)
	m.pastNotes = append(m.pastNotes, ch)

	if over := len(m.pastNotes) - maxNotes; over > 0 {
		for _, n := range m.pastNotes[0:over] {
			if oldest := m.channels[n].pn.GetOldest(); oldest != nil {
				oldest.v.Stop()
			}
		}
		m.pastNotes = m.pastNotes[over:]
	}
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPastNoteAction(ch index.Channel, na note.Action) {
	if int(ch) >= len(m.outputChannels) {
		return
	}

	for i := range m.channels {
		m.channels[i].pn.DoPastNoteAction(ch, na)
	}
}
