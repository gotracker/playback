package component

import (
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/tracing"
	"github.com/gotracker/playback/voice/vol0optimization"
)

type Vol0Optimization struct {
	settings vol0optimization.Vol0OptimizationSettings

	enabled  bool
	ticksAt0 int
}

func (c Vol0Optimization) Clone() Vol0Optimization {
	return Vol0Optimization{
		settings: c.settings,
		enabled:  c.enabled,
		ticksAt0: 0,
	}
}

func (c *Vol0Optimization) Setup(settings vol0optimization.Vol0OptimizationSettings) {
	c.settings = settings
	c.enabled = settings.Enabled
	c.Reset()
}

func (c *Vol0Optimization) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func (c *Vol0Optimization) Reset() {
	c.ticksAt0 = 0
}

func (c *Vol0Optimization) ObserveVolume(v volume.Volume) {
	if c.enabled {
		if v == 0 {
			c.ticksAt0++
		} else {
			c.ticksAt0 = 0
		}
	}
}

func (c Vol0Optimization) IsDone() bool {
	return c.enabled && c.ticksAt0 >= c.settings.MaxTicksAt0
}

func (c Vol0Optimization) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} ticksAt0{%v}",
		c.enabled,
		c.ticksAt0,
	), comment)
}
