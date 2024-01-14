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
		i   int
		age int
	}

	// first pass, try to add it to the channel listed
	{
		c := &m.channels[ch]

		if c.pn.CanAddPastNote() {
			c.pn.AddPastNote(ch, pn, m.age)
			return
		}
	}

	ages := make([]pastNoteAges, 0, len(m.outputChannels))
	// second pass, try to bump the oldest listed for the channel
	for i := range m.channels {
		c := &m.channels[i]

		if c.pn.HasPastNoteForChannel(ch) {
			ages = append(ages, pastNoteAges{
				i:   i,
				age: c.pn.GetAge(),
			})
		}
	}

	// optional third pass, if no entries with existing past notes, then find any possible
	if len(ages) == 0 {
		for i := range m.channels {
			c := &m.channels[i]

			if c.pn.CanAddPastNote() {
				c.pn.AddPastNote(ch, pn, m.age)
				return
			}

			ages = append(ages, pastNoteAges{
				i:   i,
				age: c.pn.GetAge(),
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
		return ages[i].age < ages[j].age
	})

	// jam it in
	oldest := ages[0].i
	m.channels[oldest].pn.AddPastNote(ch, pn, m.age)
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPastNoteAction(ch index.Channel, na note.Action) {
	if int(ch) >= len(m.outputChannels) {
		return
	}

	for i := range m.channels {
		m.channels[i].pn.DoPastNoteAction(ch, na)
	}
}
