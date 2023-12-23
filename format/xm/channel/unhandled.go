package channel

import (
	"fmt"

	"github.com/gotracker/playback"
	xmVolume "github.com/gotracker/playback/format/xm/volume"
	"github.com/gotracker/playback/period"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand[TPeriod period.Period] struct {
	Command Command
	Info    DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand[TPeriod]) PreStart(cs playback.Channel[TPeriod, Memory, Data], m XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledCommand[TPeriod]) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command.ToRune(), e.Info)
}

func (e UnhandledCommand[TPeriod]) Names() []string {
	return []string{
		fmt.Sprintf("UnhandledCommand(%s)", e.String()),
	}
}

// UnhandledVolCommand is an unhandled volume command
type UnhandledVolCommand[TPeriod period.Period] struct {
	Vol xmVolume.VolEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledVolCommand[TPeriod]) PreStart(cs playback.Channel[TPeriod, Memory, Data], m XM) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledVolCommand[TPeriod]) String() string {
	return fmt.Sprintf("v%0.2x", e.Vol)
}

func (e UnhandledVolCommand[TPeriod]) Names() []string {
	return []string{
		fmt.Sprintf("UnhandledVolCommand(%s)", e.String()),
	}
}
