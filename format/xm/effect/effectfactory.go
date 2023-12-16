package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	"github.com/gotracker/playback/song"
)

type EffectXM interface {
	playback.Effect
}

// VolEff is a combined effect that includes a volume effect and a standard effect
type VolEff struct {
	playback.CombinedEffect[channel.Memory]
	eff EffectXM
}

func (e VolEff) String() string {
	if e.eff == nil {
		return "..."
	}
	return fmt.Sprint(e.eff)
}

// Factory produces an effect for the provided channel pattern data
func Factory(mem *channel.Memory, data song.ChannelData) EffectXM {
	d, _ := data.(*channel.Data)
	if d == nil {
		return nil
	}

	if !d.HasCommand() {
		return nil
	}

	eff := VolEff{}
	if d.What.HasVolume() {
		ve := volumeEffectFactory(mem, d.Volume)
		if ve != nil {
			eff.Effects = append(eff.Effects, ve)
		}
	}

	if e := standardEffectFactory(mem, d); e != nil {
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
