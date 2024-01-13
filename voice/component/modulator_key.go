package component

import (
	"fmt"

	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
)

type KeyModulator struct {
	settings KeyModulatorSettings
	slimKeyModulator
	fadeout bool
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

func (k *KeyModulator) DeferredUpdate() {
	if k.keyOn == k.prevKeyOn {
		return
	}

	if k.keyOn {
		if k.settings.DeferredAttack != nil {
			k.settings.DeferredAttack()
		}
	} else {
		if k.settings.DeferredRelease != nil {
			k.settings.DeferredRelease()
		}
	}
}

func (k *KeyModulator) Attack() {
	k.slimKeyModulator.Attack()
	k.fadeout = false

	if k.settings.Attack != nil && k.keyOn != k.prevKeyOn {
		k.settings.Attack()
	}
}

func (k *KeyModulator) Release() {
	k.slimKeyModulator.Release()

	if k.settings.Release != nil && k.keyOn != k.prevKeyOn {
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
