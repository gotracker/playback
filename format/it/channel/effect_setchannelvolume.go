package channel

import (
	"fmt"

	itfile "github.com/gotracker/goaudiofile/music/tracked/it"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback"
	"github.com/gotracker/playback/period"
)

// SetChannelVolume defines a set channel volume effect
type SetChannelVolume[TPeriod period.Period] DataEffect // 'Mxx'

// Start triggers on the first tick, but before the Tick() function is called
func (e SetChannelVolume[TPeriod]) Start(cs playback.Channel[TPeriod, Memory], p playback.Playback) error {
	cs.ResetRetriggerCount()

	xx := DataEffect(e)

	cv := itfile.Volume(xx)

	vol := volume.Volume(cv.Value())
	if vol > 1 {
		vol = 1
	}

	cs.SetChannelVolume(vol)
	return nil
}

func (e SetChannelVolume[TPeriod]) String() string {
	return fmt.Sprintf("M%0.2x", DataEffect(e))
}
