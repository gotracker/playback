package effect

import (
	"fmt"

	"github.com/gotracker/playback/format/s3m/channel"
	effectIntf "github.com/gotracker/playback/format/s3m/effect/intf"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand struct {
	Command uint8
	Info    channel.DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand) PreStart(cs *channel.State, m effectIntf.S3M) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledCommand) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command+'@', e.Info)
}
