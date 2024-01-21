package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

type KeyModulator struct {
	settings KeyModulatorSettings
	slimKeyModulator
	attackTriggersRelease bool
	fadeout               bool
	deferredUpdates       []bool
}

type KeyModulatorSettings struct {
	Attack          func()
	Release         func()
	Fadeout         func()
	DeferredAttack  func()
	DeferredRelease func()
}

func (k *KeyModulator) Setup(settings KeyModulatorSettings) {
	k.settings = settings
}

func (k KeyModulator) GetAttackTriggersRelease() bool {
	return k.attackTriggersRelease
}

func (k *KeyModulator) SetAttackTriggersRelease(enabled bool) error {
	k.attackTriggersRelease = enabled
	return nil
}

func (k *KeyModulator) DeferredUpdate() {
	var deferredUpdates []bool
	deferredUpdates, k.deferredUpdates = k.deferredUpdates, nil
	for _, keyOn := range deferredUpdates {
		if keyOn {
			if k.settings.DeferredAttack != nil {
				k.settings.DeferredAttack()
			}
		} else {
			if k.settings.DeferredRelease != nil {
				k.settings.DeferredRelease()
			}
		}
	}
}

func (k *KeyModulator) Attack() {
	if k.attackTriggersRelease && k.prevKeyOn && k.keyOn {
		k.Release()
	}

	k.slimKeyModulator.Attack()
	k.fadeout = false

	if k.settings.DeferredAttack != nil {
		k.deferredUpdates = append(k.deferredUpdates, true)
	}

	if k.settings.Attack != nil {
		k.settings.Attack()
	}
}

func (k *KeyModulator) Release() {
	k.slimKeyModulator.Release()

	if k.settings.DeferredRelease != nil {
		k.deferredUpdates = append(k.deferredUpdates, false)
	}

	if k.settings.Release != nil {
		k.settings.Release()
	}
}

func (k KeyModulator) IsKeyFadeout() bool {
	return k.fadeout
}

func (k *KeyModulator) Fadeout() {
	k.fadeout = true
	if k.settings.Fadeout != nil {
		k.settings.Fadeout()
	}
}

func (k *KeyModulator) Advance() {
	k.prevKeyOn = k.keyOn
}

func (k KeyModulator) Clone(settings KeyModulatorSettings) KeyModulator {
	m := k
	m.settings = settings
	return m
}

func (k KeyModulator) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("keyOn{%v} prevKeyOn{%v}",
		k.keyOn,
		k.prevKeyOn,
	), comment)
}

func (k KeyModulator) String() string {
	switch {
	case k.fadeout:
		return "Fadeout"
	case k.keyOn:
		return "Attack"
	default:
		return "Release"
	}
}
