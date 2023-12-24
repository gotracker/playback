package component

import "github.com/gotracker/gomixing/volume"

type Vol0Optimization struct {
	enabled     bool
	ticksAt0    int
	maxTicksAt0 int
}

func (c Vol0Optimization) Clone() Vol0Optimization {
	return Vol0Optimization{
		enabled:     c.enabled,
		ticksAt0:    0,
		maxTicksAt0: c.maxTicksAt0,
	}
}

func (c *Vol0Optimization) Init(maxTicksAt0 int) {
	c.enabled = true
	c.maxTicksAt0 = maxTicksAt0
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
	return c.enabled && c.ticksAt0 >= c.maxTicksAt0
}
