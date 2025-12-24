package playback

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
	"github.com/gotracker/playback/song"
)

// Effect is an interface to command/effect
type Effect interface {
	//fmt.Stringer
	TraceData() string
}

type Effecter[TMemory song.ChannelMemory] interface {
	GetEffects(TMemory) []Effect
}

func GetEffects[TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TVolume song.Volume](mem TMemory, d TChannelData) []Effect {
	var e []Effect
	if eff, ok := any(d).(Effecter[TMemory]); ok {
		e = eff.GetEffects(mem)
	}
	return e
}

type EffectNamer interface {
	Names() []string
}

func GetEffectNames(e Effect) []string {
	if namer, ok := e.(EffectNamer); ok {
		return namer.Names()
	}
	if s, ok := e.(fmt.Stringer); ok {
		name := s.String()
		if name != "" {
			return []string{name}
		}
	}
	return nil
}

// CombinedEffect specifies multiple simultaneous effects into one
type CombinedEffect[TPeriod period.Period, TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] struct {
	Effects []Effect
}

// String returns the string for the effect list
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) String() string {
	for _, eff := range e.Effects {
		s := fmt.Sprint(eff)
		if s != "" {
			return s
		}
	}
	return ""
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Names() []string {
	var names []string
	for _, eff := range e.Effects {
		names = append(names, GetEffectNames(eff)...)
	}
	return names
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) OrderStart(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionOrderStart(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RowStart(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionRowStart(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Tick(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionTick(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionRowEnd(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) OrderEnd(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionOrderEnd(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) TraceData() string {
	return e.String()
}
