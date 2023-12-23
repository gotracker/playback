package playback

import (
	"fmt"
	"reflect"

	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/song"
)

// Effect is an interface to command/effect
type Effect interface {
	//fmt.Stringer
}

type Effecter[TMemory any] interface {
	GetEffects(*TMemory, period.Period) []Effect
}

func GetEffects[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](mem *TMemory, d TChannelData) []Effect {
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

type effectPreStartIntf[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	PreStart(Channel[TPeriod, TMemory, TChannelData], Playback) error
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
func EffectPreStart[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](e Effect, cs Channel[TPeriod, TMemory, TChannelData], p Playback) error {
	if eff, ok := e.(effectPreStartIntf[TPeriod, TMemory, TChannelData]); ok {
		if err := eff.PreStart(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectStartIntf[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	Start(Channel[TPeriod, TMemory, TChannelData], Playback) error
}

// EffectStart triggers on the first tick, but before the Tick() function is called
func EffectStart[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](e Effect, cs Channel[TPeriod, TMemory, TChannelData], p Playback) error {
	if eff, ok := e.(effectStartIntf[TPeriod, TMemory, TChannelData]); ok {
		if err := eff.Start(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectTickIntf[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	Tick(Channel[TPeriod, TMemory, TChannelData], Playback, int) error
}

// EffectTick is called on every tick
func EffectTick[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](e Effect, cs Channel[TPeriod, TMemory, TChannelData], p Playback, currentTick int) error {
	if eff, ok := e.(effectTickIntf[TPeriod, TMemory, TChannelData]); ok {
		if err := eff.Tick(cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

type effectStopIntf[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] interface {
	Stop(Channel[TPeriod, TMemory, TChannelData], Playback, int) error
}

// EffectStop is called on the last tick of the row, but after the Tick() function is called
func EffectStop[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](e Effect, cs Channel[TPeriod, TMemory, TChannelData], p Playback, lastTick int) error {
	if eff, ok := e.(effectStopIntf[TPeriod, TMemory, TChannelData]); ok {
		if err := eff.Stop(cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// CombinedEffect specifies multiple simultaneous effects into one
type CombinedEffect[TPeriod period.Period, TMemory any, TChannelData song.ChannelData] struct {
	Effects []Effect
}

// PreStart triggers when the effect enters onto the channel state
func (e CombinedEffect[TPeriod, TMemory, TChannelData]) PreStart(cs Channel[TPeriod, TMemory, TChannelData], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectPreStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e CombinedEffect[TPeriod, TMemory, TChannelData]) Start(cs Channel[TPeriod, TMemory, TChannelData], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Tick is called on every tick
func (e CombinedEffect[TPeriod, TMemory, TChannelData]) Tick(cs Channel[TPeriod, TMemory, TChannelData], p Playback, currentTick int) error {
	for _, effect := range e.Effects {
		if err := EffectTick(effect, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e CombinedEffect[TPeriod, TMemory, TChannelData]) Stop(cs Channel[TPeriod, TMemory, TChannelData], p Playback, lastTick int) error {
	for _, effect := range e.Effects {
		if err := EffectStop(effect, cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// String returns the string for the effect list
func (e CombinedEffect[TPeriod, TMemory, TChannelData]) String() string {
	for _, eff := range e.Effects {
		s := fmt.Sprint(eff)
		if s != "" {
			return s
		}
	}
	return ""
}

func (e CombinedEffect[TPeriod, TMemory, TChannelData]) Names() []string {
	var names []string
	for _, eff := range e.Effects {
		names = append(names, GetEffectNames(eff)...)
	}
	return names
}

// DoEffect runs the standard tick lifetime of an effect
func DoEffect[TPeriod period.Period, TMemory any, TChannelData song.ChannelData](e Effect, cs Channel[TPeriod, TMemory, TChannelData], p Playback, currentTick int, lastTick bool) error {
	if e == nil {
		return nil
	}

	if currentTick == 0 {
		if err := EffectStart(e, cs, p); err != nil {
			return err
		}
	}
	if err := EffectTick(e, cs, p, currentTick); err != nil {
		return err
	}
	if lastTick {
		if err := EffectStop(e, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}
