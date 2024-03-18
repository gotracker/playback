package machine

import (
	"fmt"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/mixing"
	"github.com/gotracker/playback/mixing/panning"
	"github.com/gotracker/playback/mixing/volume"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/voice/mixer"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) render(s *sampler.Sampler) (*output.PremixData, error) {
	tickDuration := m.songData.GetTickDuration(m.bpm)
	if tickDuration <= 0 {
		return nil, fmt.Errorf("unexpected tick duration: %v", tickDuration)
	}

	renderRow := render.RowRender{
		Order: int(m.ticker.current.Order),
		Row:   int(m.ticker.current.Row),
		Tick:  m.ticker.current.Tick,
	}

	premix := output.PremixData{
		SamplesLen:  int(float64(s.SampleRate) * tickDuration.Seconds()),
		MixerVolume: m.gv.ToVolume() * m.mv,
		Userdata:    &renderRow,
	}

	if m.ticker.current.Tick == 0 {
		// make a copy so it doesn't get stomped
		renderRow.RowText = m.rowStringer
	}

	details := mixer.Details{
		Mix:              s.Mixer(),
		Panmixer:         s.GetPanMixer(),
		SampleRate:       frequency.Frequency(s.SampleRate),
		StereoSeparation: s.StereoSeparation,
		Samples:          premix.SamplesLen,
		Duration:         tickDuration,
	}

	centerAheadPan := details.Panmixer.GetMixingMatrix(panning.CenterAhead, s.StereoSeparation)

	var mixData []mixing.Data
	for i := range m.actualOutputs {
		rc := &m.actualOutputs[i]

		rc.GlobalVolume = m.gv.ToVolume()

		rc.GetVoice().DumpState(index.Channel(i), m.us.Tracer)
		data, err := rc.RenderAndTick(m.ms.PeriodConverter, centerAheadPan, details)
		if err != nil {
			return nil, err
		}

		if data != nil {
			mixData = append(mixData, *data)
		} else {
			mixData = append(mixData, mixing.Data{
				Data:       details.Mix.NewMixBuffer(details.Samples),
				PanMatrix:  centerAheadPan,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: details.Samples,
			})
		}
	}

	for i := range m.virtualOutputs {
		rc := &m.virtualOutputs[i]

		var data *mixing.Data
		if rc.GetVoice() != nil {
			rc.GlobalVolume = m.gv.ToVolume()

			//rc.GetVoice().DumpState(index.Channel(i), m.us.Tracer)
			var err error
			data, err = rc.RenderAndTick(m.ms.PeriodConverter, centerAheadPan, details)
			if err != nil {
				return nil, err
			}
		}

		if data != nil {
			mixData = append(mixData, *data)
		} else {
			mixData = append(mixData, mixing.Data{
				Data:       details.Mix.NewMixBuffer(details.Samples),
				PanMatrix:  centerAheadPan,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: details.Samples,
			})
		}
	}

	if len(mixData) > 0 {
		premix.Data = append(premix.Data, mixData)
	}

	if m.opl2 != nil {
		rr := [1]mixing.Data{}
		if err := m.renderOPL2Tick(centerAheadPan, &rr[0], s.Mixer(), premix.SamplesLen); err != nil {
			return nil, err
		}
		premix.Data = append(premix.Data, rr[:])

		// make room in the mixer for the OPL2 data
		// effectively, we can do this by calculating the new number (+1) of channels from the mixer volume (channels = reciprocal of mixer volume):
		//   numChannels = (1/mv) + 1
		// then by taking the reciprocal of it:
		//   1 / numChannels
		// but that ends up being simplified to:
		//   mv / (mv + 1)
		// and we get protection from div/0 in the process - provided, of course, that the mixerVolume is not exactly -1...
		mv := premix.MixerVolume
		premix.MixerVolume /= (mv + 1)
	}

	if len(premix.Data) == 0 {
		premix.Data = append(premix.Data, mixing.ChannelData{
			mixing.Data{
				Data:       details.Mix.NewMixBuffer(details.Samples),
				PanMatrix:  centerAheadPan,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: details.Samples,
			},
		})
	}

	return &premix, nil
}
