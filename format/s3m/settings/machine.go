package settings

import (
	s3mFilter "github.com/gotracker/playback/format/s3m/filter"
	s3mOscillator "github.com/gotracker/playback/format/s3m/oscillator"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
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
		GetFilterFactory:    s3mFilter.Factory,
		GetVibratoFactory:   s3mOscillator.VibratoFactory,
		GetTremoloFactory:   s3mOscillator.TremoloFactory,
		GetPanbrelloFactory: s3mOscillator.PanbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         true,
		Quirks:              s3mQuirks,
	}

	amigaS3MSettings = settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PeriodConverter:     s3mPeriod.S3MAmigaConverter,
		GetFilterFactory:    s3mFilter.Factory,
		GetVibratoFactory:   s3mOscillator.VibratoFactory,
		GetTremoloFactory:   s3mOscillator.TremoloFactory,
		GetPanbrelloFactory: s3mOscillator.PanbrelloFactory,
		VoiceFactory:        amigaVoiceFactory,
		OPL2Enabled:         true,
		Quirks:              s3mQuirks,
	}
)
