package quirks

import (
	"github.com/gotracker/playback/filter"
	s3mFilter "github.com/gotracker/playback/format/s3m/filter"
	s3mOscillator "github.com/gotracker/playback/format/s3m/oscillator"
	s3mPanning "github.com/gotracker/playback/format/s3m/panning"
	s3mPeriod "github.com/gotracker/playback/format/s3m/period"
	s3mVolume "github.com/gotracker/playback/format/s3m/volume"
	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/machine/settings"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/oscillator"
)

type S3MMachineDefaults struct {
	AmigaPeriod      period.PeriodConverter[period.Amiga]
	FilterFactory    any
	VibratoFactory   any
	TremoloFactory   any
	PanbrelloFactory any
	ModLimits        bool
}

func GetS3MMachineSettings(profile Profile, vf voice.VoiceFactory[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]) *settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning] {
	defs := getS3MDefaults(profile)

	return &settings.MachineSettings[period.Amiga, s3mVolume.Volume, s3mVolume.FineVolume, s3mVolume.Volume, s3mPanning.Panning]{
		PeriodConverter:     defs.AmigaPeriod.(song.PeriodCalculator[period.Amiga]),
		GetFilterFactory:    defs.FilterFactory.(func(string, frequency.Frequency, any) (filter.Filter, error)),
		GetVibratoFactory:   defs.VibratoFactory.(func() (oscillator.Oscillator, error)),
		GetTremoloFactory:   defs.TremoloFactory.(func() (oscillator.Oscillator, error)),
		GetPanbrelloFactory: defs.PanbrelloFactory.(func() (oscillator.Oscillator, error)),
		VoiceFactory:        vf,
		OPL2Enabled:         true,
		ModLimits:           defs.ModLimits,
		Quirks:              Resolve(profile),
	}
}

func getS3MDefaults(profile Profile) S3MMachineDefaults {
	if def, ok := Get(profile); ok {
		if md, ok := def.MachineDefaults.(S3MMachineDefaults); ok {
			return md
		}
	}

	return S3MMachineDefaults{
		AmigaPeriod:      s3mPeriod.S3MAmigaConverter,
		FilterFactory:    s3mFilter.Factory,
		VibratoFactory:   s3mOscillator.VibratoFactory,
		TremoloFactory:   s3mOscillator.TremoloFactory,
		PanbrelloFactory: s3mOscillator.PanbrelloFactory,
		ModLimits:        false,
	}
}
