package channel

import (
	"fmt"

	itPanning "github.com/gotracker/playback/format/it/panning"
	itVolume "github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand[TPeriod period.Period] struct {
	Command Command
	Info    DataEffect
}

func (e UnhandledCommand[TPeriod]) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command.ToRune(), e.Info)
}

func (e UnhandledCommand[TPeriod]) Names() []string {
	return []string{
		fmt.Sprintf("UnhandledCommand(%s)", e.String()),
	}
}

func (e UnhandledCommand[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	if !m.IgnoreUnknownEffect() {
		panic(fmt.Sprintf("unhandled command: ce:%0.2X cp:%0.2X", e.Command, e.Info))
	}
	return nil
}

func (e UnhandledCommand[TPeriod]) TraceData() string {
	return e.String()
}

////////

// UnhandledVolCommand is an unhandled volume command
type UnhandledVolCommand[TPeriod period.Period] struct {
	Vol uint8
}

func (e UnhandledVolCommand[TPeriod]) String() string {
	return fmt.Sprintf("v%0.2x", e.Vol)
}

func (e UnhandledVolCommand[TPeriod]) Names() []string {
	return []string{
		fmt.Sprintf("UnhandledVolCommand(%s)", e.String()),
	}
}

func (e UnhandledVolCommand[TPeriod]) RowStart(ch index.Channel, m machine.Machine[TPeriod, itVolume.FineVolume, itVolume.FineVolume, itVolume.Volume, itPanning.Panning]) error {
	if !m.IgnoreUnknownEffect() {
		panic(fmt.Sprintf("unhandled command: volCmd:%0.2X", e.Vol))
	}
	return nil
}

func (e UnhandledVolCommand[TPeriod]) TraceData() string {
	return e.String()
}
