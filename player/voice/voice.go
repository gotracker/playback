package voice

import (
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/player/render"
)

// New returns a new Voice from the instrument and output channel provided
func New[TPeriod period.Period](periodConverter period.PeriodConverter[TPeriod], inst *instrument.Instrument, output *render.Channel) voice.Voice {
	switch data := inst.GetData().(type) {
	case *instrument.PCM:
		var (
			voiceFilter  filter.Filter
			pluginFilter filter.Filter
		)
		if factory := inst.GetFilterFactory(); factory != nil {
			voiceFilter = factory(inst.SampleRate, output.GetSampleRate())
		}
		if factory := inst.GetPluginFilterFactory(); factory != nil {
			pluginFilter = factory(inst.SampleRate, output.GetSampleRate())
		}
		return NewPCM[TPeriod](periodConverter, PCMConfiguration[TPeriod]{
			SampleRate:    inst.GetSampleRate(),
			InitialVolume: inst.GetDefaultVolume(),
			AutoVibrato:   inst.GetAutoVibrato(),
			Data:          data,
			OutputFilter:  output,
			VoiceFilter:   voiceFilter,
			PluginFilter:  pluginFilter,
		})
	case *instrument.OPL2:
		return NewOPL2[TPeriod](OPLConfiguration[TPeriod]{
			Chip:          output.GetOPL2Chip(),
			Channel:       output.ChannelNum,
			SampleRate:    inst.GetSampleRate(),
			InitialVolume: inst.GetDefaultVolume(),
			AutoVibrato:   inst.GetAutoVibrato(),
			Data:          data,
		})
	}
	return nil
}
