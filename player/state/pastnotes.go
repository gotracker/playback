package state

import (
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
	"github.com/heucuva/optional"
)

type pastNote[TPeriod period.Period] struct {
	ch          int
	activeState *Active[TPeriod]
}

func (pn *pastNote[TPeriod]) IsValid() bool {
	return pn.activeState.Voice != nil && !pn.activeState.Voice.IsDone()
}

type PastNotesProcessor[TPeriod period.Period] struct {
	order    []pastNote[TPeriod]
	max      optional.Value[int]
	maxPerCh optional.Value[int]
}

func (p *PastNotesProcessor[TPeriod]) Add(ch int, data *Active[TPeriod]) {
	if data == nil {
		return
	}

	if max, ok := p.max.Get(); ok {
		if c := len(p.order) - max; c > 0 {
			o := p.order[0:c]
			p.order = p.order[c:]

			for _, pn := range o {
				pn.activeState.Reset()
			}
		}
	}

	cl := pastNote[TPeriod]{
		ch:          ch,
		activeState: data,
	}

	p.order = append(p.order, cl)
}

func (p *PastNotesProcessor[TPeriod]) Do(ch int, action note.Action) {
	if action == note.ActionContinue {
		return
	}

	for _, pn := range p.order {
		if pn.ch != ch {
			continue
		}

		if !pn.IsValid() {
			continue
		}

		switch action {
		case note.ActionCut:
			pn.activeState.Reset()
		case note.ActionRelease:
			pn.activeState.Voice.Release()
		case note.ActionFadeout:
			pn.activeState.Voice.Release()
			pn.activeState.Voice.Fadeout()
		}
	}
}

func (p *PastNotesProcessor[TPeriod]) Update() {
	var nl []pastNote[TPeriod]
	for _, o := range p.order {
		if !o.IsValid() {
			o.activeState.Reset()
			continue
		}

		nl = append(nl, o)
	}
	p.order = nl
}

func (p *PastNotesProcessor[TPeriod]) GetNotesForChannel(ch int) []*Active[TPeriod] {
	var pastNotes []*Active[TPeriod]
	for _, pn := range p.order {
		if pn.ch != ch {
			continue
		}

		if !pn.IsValid() {
			continue
		}

		pastNotes = append(pastNotes, pn.activeState)
		if max, ok := p.maxPerCh.Get(); ok {
			if c := len(pastNotes) - max; c > 0 {
				o := pastNotes[0:c]
				pastNotes = pastNotes[c:]

				for _, pn := range o {
					pn.Reset()
				}
			}
		}
	}
	return pastNotes
}

func (p *PastNotesProcessor[TPeriod]) SetMax(max int) {
	p.max.Set(max)
}

func (p *PastNotesProcessor[TPeriod]) ClearMax() {
	p.max.Reset()
}

func (p *PastNotesProcessor[TPeriod]) SetMaxPerChannel(max int) {
	p.maxPerCh.Set(max)
}

func (p *PastNotesProcessor[TPeriod]) ClearMaxPerChannel() {
	p.maxPerCh.Reset()
}
