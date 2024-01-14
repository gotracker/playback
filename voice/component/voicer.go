package component

import (
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/types"
)

type Voicer[TPeriod types.Period, TMixingVolume, TVolume types.Volume] interface {
	Clone() Voicer[TPeriod, TMixingVolume, TVolume]
	GetDefaultVolume() TVolume
	GetNumChannels() int
	Attack()
	Release()
	Fadeout()
	DeferredAttack()
	DeferredRelease()
	DumpState(ch index.Channel, t tracing.Tracer, comment string)
}
