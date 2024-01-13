package playback

import (
	"fmt"
	"reflect"

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
	GetEffects(TMemory, period.Period) []Effect
}

func GetEffects[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](mem TMemory, d TChannelData) []Effect {
	var e []Effect
	if eff, ok := any(d).(Effecter[TMemory]); ok {
		var p TPeriod
		e = eff.GetEffects(mem, p)
	}
	return e
}

type EffectNamer interface {
	Names() []string
}

type effectPreStartIntf[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] interface {
	PreStart(Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], Playback) error
}

func GetEffectNames(e Effect) []string {
	if namer, ok := e.(EffectNamer); ok {
		return namer.Names()
	} else {
		typ := reflect.TypeOf(e)
		return []string{typ.Name()}
	}
}

// EffectPreStart triggers when the effect enters onto the channel state
func EffectPreStart[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](e Effect, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback) error {
	if eff, ok := e.(effectPreStartIntf[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning]); ok {
		if err := eff.PreStart(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectStartIntf[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] interface {
	Start(Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], Playback) error
}

// EffectStart triggers on the first tick, but before the Tick() function is called
func EffectStart[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](e Effect, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback) error {
	if eff, ok := e.(effectStartIntf[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning]); ok {
		if err := eff.Start(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectTickIntf[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] interface {
	OldTick(Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], Playback, int) error
}

// EffectTick is called on every tick
func EffectTick[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](e Effect, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback, currentTick int) error {
	if eff, ok := e.(effectTickIntf[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning]); ok {
		if err := eff.OldTick(cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

type effectStopIntf[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning] interface {
	Stop(Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], Playback, int) error
}

// EffectStop is called on the last tick of the row, but after the Tick() function is called
func EffectStop[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](e Effect, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback, lastTick int) error {
	if eff, ok := e.(effectStopIntf[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning]); ok {
		if err := eff.Stop(cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// CombinedEffect specifies multiple simultaneous effects into one
type CombinedEffect[TPeriod period.Period, TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume]] struct {
	Effects []Effect
}

// PreStart triggers when the effect enters onto the channel state
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) PreStart(cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectPreStart[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) Start(cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectStart[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Tick is called on every tick
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) OldTick(cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback, currentTick int) error {
	for _, effect := range e.Effects {
		if err := EffectTick[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](effect, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) Stop(cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback, lastTick int) error {
	for _, effect := range e.Effects {
		if err := EffectStop[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](effect, cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// String returns the string for the effect list
func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) String() string {
	for _, eff := range e.Effects {
		s := fmt.Sprint(eff)
		if s != "" {
			return s
		}
	}
	return ""
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) Names() []string {
	var names []string
	for _, eff := range e.Effects {
		names = append(names, GetEffectNames(eff)...)
	}
	return names
}

// DoEffect runs the standard tick lifetime of an effect
func DoEffect[TPeriod period.Period, TMemory song.ChannelMemory, TChannelData song.ChannelData[TVolume], TGlobalVolume, TMixingVolume, TVolume song.Volume, TPanning song.Panning](e Effect, cs Channel[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume, TPanning], p Playback, currentTick int, lastTick bool) error {
	if e == nil {
		return nil
	}

	if currentTick == 0 {
		if err := EffectStart[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](e, cs, p); err != nil {
			return err
		}
	}
	if err := EffectTick[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](e, cs, p, currentTick); err != nil {
		return err
	}
	if lastTick {
		if err := EffectStop[TPeriod, TMemory, TChannelData, TGlobalVolume, TMixingVolume, TVolume](e, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

////////

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) RowStart(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionRowStart(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) PreTick(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionPreTick(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) Tick(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionTick(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) PostTick(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning], tick int) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionPostTick(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) RowEnd(ch index.Channel, m machine.Machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) error {
	for _, effect := range e.Effects {
		if err := m.DoInstructionRowEnd(ch, effect); err != nil {
			return err
		}
	}
	return nil
}

func (e CombinedEffect[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning, TMemory, TChannelData]) TraceData() string {
	return e.String()
}
