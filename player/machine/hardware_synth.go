package machine

import (
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/voice/mixer"
)

// mixerVolumeAdjuster updates mixer volume to accommodate hardware synth output.
type mixerVolumeAdjuster func(volume.Volume) volume.Volume

// hardwareSynth abstracts rendering of hardware synth frames (e.g., OPL2).
type hardwareSynth interface {
	RenderTick(centerAheadPan panning.PanMixer, details mixer.Details) (mixing.Data, mixerVolumeAdjuster, error)
}
