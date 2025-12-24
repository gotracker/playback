package channel

import (
	"testing"

	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

func TestEffectFactoryNoCommandReturnsNil(t *testing.T) {
	var mem Memory
	d := Data[period.Amiga]{}
	if eff := EffectFactory[period.Amiga](&mem, d); eff != nil {
		t.Fatalf("expected nil effect when no command present")
	}

	if eff := EffectFactory[period.Amiga](nil, nil); eff != nil {
		t.Fatalf("expected nil effect when data is nil")
	}

	if got := d.GetVolume(); got != xmVolume.XmVolume(0) {
		t.Fatalf("expected zero volume for empty data, got %v", got)
	}
}
