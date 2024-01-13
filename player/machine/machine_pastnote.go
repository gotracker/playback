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

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) addPastNote(pn voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) {
	type pastNoteAges struct {
		i   int
		age int
	}

	ages := make([]pastNoteAges, 0, len(m.outputChannels))
	// first pass, just try to add it
	for i := range m.channels {
		c := &m.channels[i]

		if c.pn.CanAddPastNote() {
			c.pn.AddPastNote(pn, m.age)
			return
		}

		ages = append(ages, pastNoteAges{
			i:   i,
			age: c.pn.GetAge(),
		})
	}

	if len(ages) == 0 {
		// impossible, but just in case
		pn.Stop()
		return
	}

	// second pass, find the oldest
	sort.Slice(ages, func(i, j int) bool {
		return ages[i].age < ages[j].age
	})

	// jam it in
	oldest := ages[0].i
	m.channels[oldest].pn.AddPastNote(pn, m.age)
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) doPastNoteAction(ch index.Channel, na note.Action) {
	if int(ch) >= len(m.outputChannels) {
		return
	}

	m.channels[ch].pn.DoPastNoteAction(na)
}
