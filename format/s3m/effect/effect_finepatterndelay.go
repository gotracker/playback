package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/s3m/channel"
	effectIntf "github.com/gotracker/playback/format/s3m/effect/intf"
)

// FinePatternDelay defines an fine pattern delay effect
type FinePatternDelay ChannelCommand // 'S6x'

// Start triggers on the first tick, but before the Tick() function is called
func (e FinePatternDelay) Start(cs S3MChannel, p playback.Playback) error {
	cs.ResetRetriggerCount()

	x := channel.DataEffect(e) & 0xf

	m := p.(effectIntf.S3M)
	return m.AddRowTicks(int(x))
}

func (e FinePatternDelay) String() string {
	return fmt.Sprintf("S%0.2x", channel.DataEffect(e))
}
