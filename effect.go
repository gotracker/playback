package playback

import "fmt"

// Effect is an interface to command/effect
type Effect interface{}

// Effect is an interface to command/effect
type Effecter[TChannelState ChannelState] interface {
	Effect
	//fmt.Stringer
}

type EffectFactory[TChannelData any, TChannelState ChannelState] func(cs *TChannelState, cd *TChannelData) Effecter[TChannelState]

type effectPreStartIntf[TChannelState ChannelState] interface {
	PreStart(*TChannelState, Playback) error
}

// EffectPreStart triggers when the effect enters onto the channel state
func EffectPreStart[TChannelState ChannelState](e Effecter[TChannelState], cs *TChannelState, p Playback) error {
	if eff, ok := e.(effectPreStartIntf[TChannelState]); ok {
		if err := eff.PreStart(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectStartIntf[TChannelState ChannelState] interface {
	Start(*TChannelState, Playback) error
}

// EffectStart triggers on the first tick, but before the Tick() function is called
func EffectStart[TChannelState ChannelState](e Effecter[TChannelState], cs *TChannelState, p Playback) error {
	if eff, ok := e.(effectStartIntf[TChannelState]); ok {
		if err := eff.Start(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectTickIntf[TChannelState ChannelState] interface {
	Tick(*TChannelState, Playback, int) error
}

// EffectTick is called on every tick
func EffectTick[TChannelState ChannelState](e Effecter[TChannelState], cs *TChannelState, p Playback, currentTick int) error {
	if eff, ok := e.(effectTickIntf[TChannelState]); ok {
		if err := eff.Tick(cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

type effectStopIntf[TChannelState ChannelState] interface {
	Stop(*TChannelState, Playback, int) error
}

// EffectStop is called on the last tick of the row, but after the Tick() function is called
func EffectStop[TChannelState ChannelState](e Effecter[TChannelState], cs *TChannelState, p Playback, lastTick int) error {
	if eff, ok := e.(effectStopIntf[TChannelState]); ok {
		if err := eff.Stop(cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// CombinedEffect[TChannelState] specifies multiple simultaneous effects into one
type CombinedEffect[TChannelState ChannelState] struct {
	Effects []Effecter[TChannelState]
}

// PreStart triggers when the effect enters onto the channel state
func (e CombinedEffect[TChannelState]) PreStart(cs *TChannelState, p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectPreStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e CombinedEffect[TChannelState]) Start(cs *TChannelState, p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Tick is called on every tick
func (e CombinedEffect[TChannelState]) Tick(cs *TChannelState, p Playback, currentTick int) error {
	for _, effect := range e.Effects {
		if err := EffectTick(effect, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e CombinedEffect[TChannelState]) Stop(cs *TChannelState, p Playback, lastTick int) error {
	for _, effect := range e.Effects {
		if err := EffectStop(effect, cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// String returns the string for the effect list
func (e CombinedEffect[TChannelState]) String() string {
	for _, eff := range e.Effects {
		s := fmt.Sprint(eff)
		if s != "" {
			return s
		}
	}
	return ""
}

// DoEffect runs the standard tick lifetime of an effect
func DoEffect[TChannelState ChannelState](e Effecter[TChannelState], cs *TChannelState, p Playback, currentTick int, lastTick bool) error {
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
