package types

import (
	"testing"

	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
)

type testVol volume.Volume

func (testVol) IsInvalid() bool                { return false }
func (testVol) IsUseInstrumentVol() bool       { return false }
func (v testVol) ToVolume() volume.Volume      { return volume.Volume(v) }
func (testVol) AddDelta(d VolumeDelta) testVol { return testVol(volume.Volume(1) + volume.Volume(d)) }

type maxVol volume.Volume

func (maxVol) IsInvalid() bool           { return false }
func (maxVol) IsUseInstrumentVol() bool  { return false }
func (v maxVol) ToVolume() volume.Volume { return volume.Volume(v) }
func (maxVol) GetMax() maxVol            { return maxVol(2) }

type testPan float32

func (testPan) IsInvalid() bool              { return false }
func (testPan) ToPosition() panning.Position { return panning.CenterAhead }
func (testPan) AddDelta(d PanDelta) testPan  { return testPan(float32(0.5) + float32(d)) }

type panInfo float32

func (panInfo) IsInvalid() bool              { return false }
func (panInfo) ToPosition() panning.Position { return panning.CenterAhead }
func (panInfo) GetDefault() panInfo          { return panInfo(0.25) }
func (panInfo) GetMax() panInfo              { return panInfo(1) }
func (panInfo) AddDelta(d PanDelta) panInfo  { return panInfo(float32(0.25) + float32(d)) }

func TestAddVolumeDelta(t *testing.T) {
	v := testVol(1)
	res := AddVolumeDelta(v, VolumeDelta(2))
	if res != 3 {
		t.Fatalf("expected 3, got %v", res)
	}
}

func TestGetMaxVolume(t *testing.T) {
	if got := GetMaxVolume[maxVol](); got != 2 {
		t.Fatalf("expected max volume 2, got %v", got)
	}
}

func TestPanningHelpers(t *testing.T) {
	p := testPan(0)
	if res := AddPanningDelta(p, PanDelta(0.1)); res != 0.6 {
		t.Fatalf("expected pan 0.6, got %v", res)
	}
	if def := GetPanDefault[panInfo](); def != 0.25 {
		t.Fatalf("expected default 0.25, got %v", def)
	}
	if mx := GetPanMax[panInfo](); mx != 1 {
		t.Fatalf("expected max 1, got %v", mx)
	}
}
