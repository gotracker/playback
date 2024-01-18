package settings

import (
	"fmt"

	"github.com/gotracker/playback/filter"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
	oscillatorImpl "github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/voice/oscillator"
)

func GetMachineSettings(modLimits bool) *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	if modLimits {
		return &amigaMOD31Settings
	}
	return &amigaS3MSettings
}

var (
	s3mQuirks = settings.MachineQuirks{
		PreviousPeriodUsesModifiedPeriod:   true,
		PortaToNoteUsesModifiedPeriod:      true,
		DoNotProcessEffectsOnMutedChannels: true,
	}

	amigaMOD31Settings = settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PeriodConverter:     s3mPeriod.S3MAmigaConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         true,
		Quirks:              s3mQuirks,
	}

	amigaS3MSettings = settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PeriodConverter:     s3mPeriod.S3MAmigaConverter,
		GetFilterFactory:    filterFactory,
		GetVibratoFactory:   vibratoFactory,
		GetTremoloFactory:   tremoloFactory,
		GetPanbrelloFactory: panbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         true,
		Quirks:              s3mQuirks,
	}
)

func filterFactory(name string) (settings.FilterFactoryFunc, error) {
	switch name {
	case "amigalpf":
		return func(instrument frequency.Frequency) (filter.Filter, error) {
			lpf := filter.NewAmigaLPF(instrument)
			return lpf, nil
		}, nil

	default:
		return nil, fmt.Errorf("unsupported filter: %q", name)
	}
}

func vibratoFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}

func tremoloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}

func panbrelloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}
