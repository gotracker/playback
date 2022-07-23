package effect

import (
	"fmt"

	"github.com/gotracker/playback/format/xm/channel"
	effectIntf "github.com/gotracker/playback/format/xm/effect/intf"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand struct {
	Command channel.Command
	Info    channel.DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand) PreStart(cs *channel.State, m effectIntf.XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledCommand) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command.ToRune(), e.Info)
}

// UnhandledVolCommand is an unhandled volume command
type UnhandledVolCommand struct {
	Vol xmVolume.VolEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledVolCommand) PreStart(cs *channel.State, m effectIntf.XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledVolCommand) String() string {
	return fmt.Sprintf("v%0.2x", e.Vol)
}
