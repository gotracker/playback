package render

import (
	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/voice"
	"github.com/gotracker/playback/voice/mixer"
)

type ChannelIntf interface {
	ApplyFilter(dry volume.Matrix) volume.Matrix
	GetPremixVolume() volume.Volume
}

// Channel is the important bits to make output to a particular downmixing channel work
type Channel[TPeriod period.Period] struct {
	PluginFilter filter.Filter
	OutputFilter filter.Filter
	GetOPL2Chip  func() OPL2Chip
	GlobalVolume volume.Volume // this is the channel's version of the GlobalVolume
	Voice        voice.Voice
}

func (c *Channel[TPeriod]) RenderAndAdvance(pc period.PeriodConverter[TPeriod], centerAheadPan volume.Matrix, details mixer.Details) (*mixing.Data, error) {
	if filt := c.PluginFilter; filt != nil {
		filt.SetPlaybackRate(details.SampleRate)
	}

	if filt := c.OutputFilter; filt != nil {
		filt.SetPlaybackRate(details.SampleRate)
	}

	data, err := voice.RenderAndAdvance(c.Voice, pc, centerAheadPan, details, c)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}
	return data, nil
}

func (c *Channel[TPeriod]) StopVoice() {
	if c.Voice == nil {
		return
	}

	c.Voice.Stop()
	c.Voice = nil
}

// ApplyFilter will apply the channel filter, if there is one.
func (c *Channel[TPeriod]) ApplyFilter(dry volume.Matrix) volume.Matrix {
	if dry.Channels == 0 {
		return dry
	}
	wet := dry
	if c.PluginFilter != nil {
		wet = c.PluginFilter.Filter(wet)
	}
	wet = wet.Apply(c.GlobalVolume)
	if c.OutputFilter != nil {
		wet = c.OutputFilter.Filter(wet)
	}
	return wet
}

func (Channel[TPeriod]) SetFilterEnvelopeValue(envVal uint8) {
}
