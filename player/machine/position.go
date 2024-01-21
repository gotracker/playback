package machine

import "github.com/gotracker/playback/index"

type Position struct {
	Order index.Order
	Row   index.Row
	Tick  int
}
