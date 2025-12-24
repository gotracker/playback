package instrument

import (
	"testing"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/sampling"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/note"
	"github.com/gotracker/playback/period"
)

type stubPeriod struct{}

func (stubPeriod) IsInvalid() bool { return false }

type stubVol float32

func (stubVol) IsInvalid() bool           { return false }
func (stubVol) IsUseInstrumentVol() bool  { return false }
func (v stubVol) ToVolume() volume.Volume { return volume.Volume(v) }

type stubPan float32

func (stubPan) IsInvalid() bool              { return false }
func (stubPan) ToPosition() panning.Position { return panning.CenterAhead }

type stubID struct{ idx int }

func (stubID) IsEmpty() bool                    { return false }
func (id stubID) GetIndexAndSample() (int, int) { return id.idx, 0 }
func (id stubID) String() string                { return "id" }

type stubData struct{ l sampling.Pos }

func (s stubData) GetLength() sampling.Pos { return s.l }

func TestInstrumentFinetuneOverridesStatic(t *testing.T) {
	inst := Instrument[stubPeriod, stubVol, stubVol, stubPan]{
		Static: StaticValues[stubPeriod, stubVol, stubVol, stubPan]{Finetune: 2},
	}

	if inst.GetFinetune() != 2 {
		t.Fatalf("expected default finetune 2")
	}

	inst.SetFinetune(5)
	if inst.GetFinetune() != 5 {
		t.Fatalf("expected override finetune 5")
	}
}

func TestInstrumentReleaseStopNotesForOPL2(t *testing.T) {
	oplInst := Instrument[period.Amiga, stubVol, stubVol, stubPan]{
		Inst: &OPL2{},
	}
	otherInst := Instrument[period.Amiga, stubVol, stubVol, stubPan]{
		Inst: stubData{},
	}

	n := note.StopOrReleaseNote{}
	if !oplInst.IsReleaseNote(n) || !oplInst.IsStopNote(n) {
		t.Fatalf("expected OPL2 stop/release to map to release/stop")
	}
	if otherInst.IsReleaseNote(n) || otherInst.IsStopNote(n) {
		t.Fatalf("expected non-OPL2 stop/release to stay non-release/stop")
	}
}

func TestInstrumentGetters(t *testing.T) {
	inst := Instrument[stubPeriod, stubVol, stubVol, stubPan]{
		Static: StaticValues[stubPeriod, stubVol, stubVol, stubPan]{
			Volume:             stubVol(0.5),
			RelativeNoteNumber: 3,
			NewNoteAction:      note.ActionRelease,
			ID:                 stubID{idx: 4},
			VoiceFilter:        filter.Info{Name: "vf"},
			PluginFilter:       filter.Info{Name: "pf"},
		},
		Inst:       stubData{l: sampling.Pos{Pos: 7}},
		SampleRate: 1234,
	}

	if inst.GetDefaultVolume() != stubVol(0.5) || inst.GetDefaultVolumeGeneric() != volume.Volume(0.5) {
		t.Fatalf("unexpected default volume values")
	}
	if inst.GetLength() != (sampling.Pos{Pos: 7}) {
		t.Fatalf("unexpected length")
	}
	if inst.GetID().(stubID).idx != 4 {
		t.Fatalf("unexpected id")
	}
	if inst.GetSemitoneShift() != 3 {
		t.Fatalf("unexpected semitone shift")
	}
	if inst.GetNewNoteAction() != note.ActionRelease {
		t.Fatalf("unexpected new note action")
	}
	if inst.GetVoiceFilterInfo().Name != "vf" || inst.GetPluginFilterInfo().Name != "pf" {
		t.Fatalf("unexpected filter info")
	}
	if inst.GetSampleRate() != 1234 {
		t.Fatalf("unexpected sample rate")
	}
	inst.SetSampleRate(4321)
	if inst.GetSampleRate() != 4321 {
		t.Fatalf("expected updated sample rate")
	}
}
