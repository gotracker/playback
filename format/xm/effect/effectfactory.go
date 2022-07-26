package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
)

type EffectXM interface {
	playback.Effect
}

// VolEff is a combined effect that includes a volume effect and a standard effect
type VolEff struct {
	playback.CombinedEffect[channel.Memory, channel.Data]
	eff EffectXM
}

func (e VolEff) String() string {
	if e.eff == nil {
		return "..."
	}
	return fmt.Sprint(e.eff)
}

// Factory produces an effect for the provided channel pattern data
func Factory(mem *channel.Memory, data *channel.Data) EffectXM {
	if data == nil {
		return nil
	}

	if !data.HasCommand() {
		return nil
	}

	eff := VolEff{}
	if data.What.HasVolume() {
		ve := volumeEffectFactory(mem, data.Volume)
		if ve != nil {
			eff.Effects = append(eff.Effects, ve)
		}
	}

	if e := standardEffectFactory(mem, data); e != nil {
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
