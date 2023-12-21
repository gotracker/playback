package channel

import (
	"fmt"
)

// UnhandledCommand is an unhandled command
type UnhandledCommand struct {
	Command uint8
	Info    DataEffect
}

// PreStart triggers when the effect enters onto the channel state
func (e UnhandledCommand) PreStart(cs S3MChannel, m S3M) error {
	if !m.IgnoreUnknownEffect() {
		panic("unhandled command")
	}
	return nil
}

func (e UnhandledCommand) String() string {
	return fmt.Sprintf("%c%0.2x", e.Command+'@', e.Info)
}

func (e UnhandledCommand) Names() []string {
	return []string{
		fmt.Sprintf("UnhandledCommand(%s)", e.String()),
	}
}
