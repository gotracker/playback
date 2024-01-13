package settings

import (
	"fmt"

	"github.com/gotracker/playback/filter"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	voiceOscillator "github.com/gotracker/playback/voice/oscillator"
)

func GetMachineSettings() *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	return &amigaSettings
}

var (
	amigaSettings = settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PeriodConverter:     s3mPeriod.AmigaConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         true,
	}
)

func filterFactory(name string) (settings.FilterFactoryFunc, error) {
	switch name {
	case "amigalpf":
		return func(instrument, playback period.Frequency) (filter.Filter, error) {
			lpf := filter.NewAmigaLPF(instrument, playback)
			return lpf, nil
		}, nil

	default:
		return nil, fmt.Errorf("unsupported filter: %q", name)
	}
}

func vibratoFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewProtrackerOscillator(), nil
}

func tremoloFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewProtrackerOscillator(), nil
}

func panbrelloFactory() (voiceOscillator.Oscillator, error) {
	return oscillator.NewProtrackerOscillator(), nil
}
