package autovibrato

import (
	"github.com/gotracker/playback/voice/oscillator"
	"github.com/gotracker/playback/voice/types"
)

type AutoVibratoSettings[TPeriod types.Period] struct {
	AutoVibratoConfig[TPeriod]
	Factory func(string) (oscillator.Oscillator, error)
}
