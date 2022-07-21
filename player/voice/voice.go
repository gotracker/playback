package voice

import (
	"github.com/gotracker/playback/voice"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/instrument"
	"github.com/gotracker/playback/player/output"
)

// New returns a new Voice from the instrument and output channel provided
func New(inst *instrument.Instrument, output *output.Channel) voice.Voice {
	switch data := inst.GetData().(type) {
	case *instrument.PCM:
		var (
			voiceFilter  filter.Filter
			pluginFilter filter.Filter
		)
		if factory := inst.GetFilterFactory(); factory != nil {
			voiceFilter = factory(inst.C2Spd, output.GetSampleRate())
		}
		if factory := inst.GetPluginFilterFactory(); factory != nil {
			pluginFilter = factory(inst.C2Spd, output.GetSampleRate())
		}
		return NewPCM(PCMConfiguration{
			C2SPD:         inst.GetC2Spd(),
			InitialVolume: inst.GetDefaultVolume(),
			AutoVibrato:   inst.GetAutoVibrato(),
			DataIntf:      data,
			OutputFilter:  output,
			VoiceFilter:   voiceFilter,
			PluginFilter:  pluginFilter,
		})
	case *instrument.OPL2:
		return NewOPL2(OPLConfiguration{
			Chip:          output.GetOPL2Chip(),
			Channel:       output.ChannelNum,
			C2SPD:         inst.GetC2Spd(),
			InitialVolume: inst.GetDefaultVolume(),
			AutoVibrato:   inst.GetAutoVibrato(),
			DataIntf:      data,
		})
	}
	return nil
}
