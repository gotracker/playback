package effect

import (
	"fmt"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/format/it/channel"
	effectIntf "github.com/gotracker/playback/format/it/effect/intf"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand struct {
	Command channel.Command
	Info    channel.DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand) PreStart(cs playback.Channel[channel.Memory, channel.Data], m effectIntf.IT) error {
	if !m.IgnoreUnknownEffect() {
		panic(fmt.Sprintf("unhandled command: ce:%0.2X cp:%0.2X", e.Command, e.Info))
	}
	return nil
}

func (e UnhandledCommand) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command.ToRune(), e.Info)
}

// UnhandledVolCommand is an unhandled volume command
type UnhandledVolCommand struct {
	Vol uint8
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledVolCommand) PreStart(cs playback.Channel[channel.Memory, channel.Data], m effectIntf.IT) error {
	if !m.IgnoreUnknownEffect() {
		panic(fmt.Sprintf("unhandled command: volCmd:%0.2X", e.Vol))
	}
	return nil
}

func (e UnhandledVolCommand) String() string {
	return fmt.Sprintf("v%0.2x", e.Vol)
}
