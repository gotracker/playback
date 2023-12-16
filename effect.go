package playback

import "fmt"

// Effect is an interface to command/effect
type Effect interface {
	//fmt.Stringer
}

type effectPreStartIntf[TMemory any] interface {
	PreStart(Channel[TMemory], Playback) error
}

// EffectPreStart triggers when the effect enters onto the channel state
func EffectPreStart[TMemory any](e Effect, cs Channel[TMemory], p Playback) error {
	if eff, ok := e.(effectPreStartIntf[TMemory]); ok {
		if err := eff.PreStart(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectStartIntf[TMemory any] interface {
	Start(Channel[TMemory], Playback) error
}

// EffectStart triggers on the first tick, but before the Tick() function is called
func EffectStart[TMemory any](e Effect, cs Channel[TMemory], p Playback) error {
	if eff, ok := e.(effectStartIntf[TMemory]); ok {
		if err := eff.Start(cs, p); err != nil {
			return err
		}
	}
	return nil
}

type effectTickIntf[TMemory any] interface {
	Tick(Channel[TMemory], Playback, int) error
}

// EffectTick is called on every tick
func EffectTick[TMemory any](e Effect, cs Channel[TMemory], p Playback, currentTick int) error {
	if eff, ok := e.(effectTickIntf[TMemory]); ok {
		if err := eff.Tick(cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

type effectStopIntf[TMemory any] interface {
	Stop(Channel[TMemory], Playback, int) error
}

// EffectStop is called on the last tick of the row, but after the Tick() function is called
func EffectStop[TMemory any](e Effect, cs Channel[TMemory], p Playback, lastTick int) error {
	if eff, ok := e.(effectStopIntf[TMemory]); ok {
		if err := eff.Stop(cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// CombinedEffect specifies multiple simultaneous effects into one
type CombinedEffect[TMemory any] struct {
	Effects []Effect
}

// PreStart triggers when the effect enters onto the channel state
func (e CombinedEffect[TMemory]) PreStart(cs Channel[TMemory], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectPreStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Start triggers on the first tick, but before the Tick() function is called
func (e CombinedEffect[TMemory]) Start(cs Channel[TMemory], p Playback) error {
	for _, effect := range e.Effects {
		if err := EffectStart(effect, cs, p); err != nil {
			return err
		}
	}
	return nil
}

// Tick is called on every tick
func (e CombinedEffect[TMemory]) Tick(cs Channel[TMemory], p Playback, currentTick int) error {
	for _, effect := range e.Effects {
		if err := EffectTick(effect, cs, p, currentTick); err != nil {
			return err
		}
	}
	return nil
}

// Stop is called on the last tick of the row, but after the Tick() function is called
func (e CombinedEffect[TMemory]) Stop(cs Channel[TMemory], p Playback, lastTick int) error {
	for _, effect := range e.Effects {
		if err := EffectStop(effect, cs, p, lastTick); err != nil {
			return err
		}
	}
	return nil
}

// String returns the string for the effect list
func (e CombinedEffect[TMemory]) String() string {
	for _, eff := range e.Effects {
		s := fmt.Sprint(eff)
		if s != "" {
			return s
		}
	}
	return ""
}

// DoEffect runs the standard tick lifetime of an effect
func DoEffect[TMemory any](e Effect, cs Channel[TMemory], p Playback, currentTick int, lastTick bool) error {
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
