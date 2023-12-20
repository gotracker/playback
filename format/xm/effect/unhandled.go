package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/effect/intf"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand[TPeriod period.Period] struct {
	Command channel.Command
	Info    channel.DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], m effectIntf.XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledCommand[TPeriod]) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command.ToRune(), e.Info)
}

// UnhandledVolCommand is an unhandled volume command
type UnhandledVolCommand[TPeriod period.Period] struct {
	Vol xmVolume.VolEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledVolCommand[TPeriod]) PreStart(cs playback.Channel[TPeriod, channel.Memory], m effectIntf.XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledVolCommand[TPeriod]) String() string {
	return fmt.Sprintf("v%0.2x", e.Vol)
}
