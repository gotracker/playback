package output

import "testing"

func TestPremixDataZeroValue(t *testing.T) {
	var p PremixData
	if p.SamplesLen != 0 {
		t.Fatalf("expected SamplesLen=0, got %d", p.SamplesLen)
	}
	if p.Data != nil {
		t.Fatalf("expected nil Data slice, got %#v", p.Data)
	}
	if p.MixerVolume != 0 {
		t.Fatalf("expected zero MixerVolume, got %v", p.MixerVolume)
	}
	if p.Userdata != nil {
		t.Fatalf("expected nil Userdata, got %#v", p.Userdata)
	}
}
