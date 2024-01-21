package tracing

import (
	"fmt"

	"github.com/gotracker/playback/index"
)

type Tick struct {
	Order index.Order
	Row   index.Row
	Tick  int
}

func (t Tick) Equals(rhs Tick) bool {
	return t.Tick == rhs.Tick && t.Row == rhs.Row && t.Order == rhs.Order
}

func (t Tick) String() string {
	ts := fmt.Sprint(t.Tick)
	if len(ts) < 2 {
		ts = " " + ts
	}
	return fmt.Sprintf("%03d:%03d ", t.Order, t.Row) + ts
}

func (t Tick) GetTick() Tick {
	return t
}

type Ticker interface {
	fmt.Stringer
	GetTick() Tick
}
