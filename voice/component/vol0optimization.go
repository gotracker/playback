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

	unkeyed struct {
		enabled bool
	}
	keyed struct {
		rowsAt0 int
	}
}

func (c Vol0Optimization) Clone() Vol0Optimization {
	m := c
	return m
}

func (c *Vol0Optimization) Setup(settings vol0optimization.Vol0OptimizationSettings) {
	c.settings = settings
	c.unkeyed.enabled = settings.Enabled
	c.Reset()
}

func (c *Vol0Optimization) SetEnabled(enabled bool) {
	c.unkeyed.enabled = enabled
}

func (c *Vol0Optimization) Reset() {
	c.keyed.rowsAt0 = 0
}

func (c *Vol0Optimization) ObserveVolume(v volume.Volume) {
	if c.unkeyed.enabled {
		if v == 0 {
			c.keyed.rowsAt0++
		} else {
			c.keyed.rowsAt0 = 0
		}
	}
}

func (c Vol0Optimization) IsDone() bool {
	return c.unkeyed.enabled && c.keyed.rowsAt0 >= c.settings.MaxRowsAt0
}

func (c Vol0Optimization) DumpState(ch index.Channel, t tracing.Tracer, comment string) {
	t.TraceChannelWithComment(ch, fmt.Sprintf("enabled{%v} rowsAt0{%v}",
		c.unkeyed.enabled,
		c.keyed.rowsAt0,
	), comment)
}
