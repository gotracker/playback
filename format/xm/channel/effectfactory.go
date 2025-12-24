package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmPanning "github.com/gotracker/playback/format/xm/panning"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

type EffectXM = playback.Effect

// VolEff is a combined effect that includes a volume effect and a standard effect
type VolEff[TPeriod period.Period] struct {
	playback.CombinedEffect[TPeriod, xmVolume.XmVolume, xmVolume.XmVolume, xmVolume.XmVolume, xmPanning.Panning]
	eff EffectXM
}

func (e VolEff[TPeriod]) String() string {
	if e.eff == nil {
		return "..."
	}
	return fmt.Sprint(e.eff)
}

func (e VolEff[TPeriod]) Names() []string {
	names := e.CombinedEffect.Names()
	if e.eff != nil {
		names = append(names, playback.GetEffectNames(e.eff)...)
	}
	return names
}

func (e VolEff[TPeriod]) TraceData() string {
	return e.String()
}

// Factory produces an effect for the provided channel pattern data
func EffectFactory[TPeriod period.Period](mem *Memory, data song.ChannelData[xmVolume.XmVolume]) EffectXM {
	if data == nil {
		return nil
	}

	d, _ := data.(Data[TPeriod])
	if !d.HasCommand() {
		return nil
	}

	var eff VolEff[TPeriod]
	if d.What.HasVolume() {
		ve := volumeEffectFactory[TPeriod](mem, d.Volume)
		if ve != nil {
			eff.Effects = append(eff.Effects, ve)
		}
	}

	if e := standardEffectFactory[TPeriod](mem, d); e != nil {
		eff.Effects = append(eff.Effects, e)
		eff.eff = e
	}

	switch len(eff.Effects) {
	case 0:
		return nil
	case 1:
		return eff.Effects[0]
	default:
		return &eff
	}
}
