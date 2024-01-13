package machine

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/note"
	playerRender "github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/state/render"
	"github.com/gotracker/playback/voice"
)

type pastNote[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	v   voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	age int
}

type pastNotes[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] struct {
	MaxPastNotes int
	pn           []pastNote[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
}

func (o *pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) CanAddPastNote() bool {
	return len(o.pn) < o.MaxPastNotes
}

func (o *pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) AddPastNote(v voice.RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], age int) {
	o.pn = append(o.pn, pastNote[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]{
		v:   v,
		age: age,
	})
	if over := len(o.pn) - o.MaxPastNotes; over > 0 {
		for _, n := range o.pn[0:over] {
			n.v.Stop()
		}
		o.pn = o.pn[over:]
	}
}

func (o pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) GetAge() int {
	if len(o.pn) == 0 {
		return 0
	}

	return o.pn[0].age
}

func (o *pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) UpdatePastNotes() {
	var updated []pastNote[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
	for _, pn := range o.pn {
		if pn.v.IsDone() {
			pn.v.Stop()
			continue
		}
		updated = append(updated, pn)
	}
	o.pn = updated
}

func (o *pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) DoPastNoteAction(na note.Action) {
	for _, pn := range o.pn {
		switch na {
		case note.ActionCut:
			pn.v.Stop()
		case note.ActionRelease:
			pn.v.Release()
		case note.ActionFadeout:
			pn.v.Release()
			pn.v.Fadeout()
		case note.ActionRetrigger:
			pn.v.Release()
			pn.v.Attack()

		case note.ActionContinue:
			fallthrough
		default:
			// nothing
		}
	}
}

func (p *pastNotes[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Render(centerAheadPan volume.Matrix, details render.Details, rc *playerRender.Channel[TGlobalVolume, TMixingVolume, TPanning]) ([]mixing.Data, error) {
	var mixData []mixing.Data

	for _, pn := range p.pn {
		if pn.v.IsDone() {
			continue
		}

		if ampMod, ok := pn.v.(voice.AmpModulator[TGlobalVolume, TMixingVolume, TVolume]); ok && !ampMod.IsActive() {
			continue
		}

		data, err := pn.v.Render(centerAheadPan, details, rc)
		if err != nil {
			return nil, err
		}

		if data != nil {
			mixData = append(mixData, *data)
		}
	}
	return mixData, nil
}