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
	frame, err := m.prepareRenderFrame(s)
	if err != nil {
		return nil, err
	}

	mixData, err := m.renderVoices(frame)
	if err != nil {
		return nil, err
	}

	if len(mixData) > 0 {
		frame.premix.Data = append(frame.premix.Data, mixData)
	}

	if err := m.mixHardwareSynths(frame, &frame.premix); err != nil {
		return nil, err
	}

	m.normalizePremix(&frame)

	return &frame.premix, nil
}

type renderFrame struct {
	renderRow      render.RowRender
	premix         output.PremixData
	details        mixer.Details
	centerAheadPan panning.PanMixer
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) prepareRenderFrame(s *sampler.Sampler) (renderFrame, error) {
	tickDuration := m.songData.GetTickDuration(m.bpm)
	if tickDuration <= 0 {
		return renderFrame{}, fmt.Errorf("unexpected tick duration: %v", tickDuration)
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

	return renderFrame{
		renderRow:      renderRow,
		premix:         premix,
		details:        details,
		centerAheadPan: centerAheadPan,
	}, nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) renderVoices(frame renderFrame) ([]mixing.Data, error) {
	var mixData []mixing.Data

	for i := range m.actualOutputs {
		rc := &m.actualOutputs[i]

		rc.GlobalVolume = m.gv.ToVolume()

		rc.GetVoice().DumpState(index.Channel(i), m.us.Tracer)
		data, err := rc.RenderAndTick(m.ms.PeriodConverter, frame.centerAheadPan, frame.details)
		if err != nil {
			return nil, err
		}

		if data != nil {
			mixData = append(mixData, *data)
		} else {
			mixData = append(mixData, mixing.Data{
				Data:       frame.details.Mix.NewMixBuffer(frame.details.Samples),
				PanMatrix:  frame.centerAheadPan,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: frame.details.Samples,
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
			data, err = rc.RenderAndTick(m.ms.PeriodConverter, frame.centerAheadPan, frame.details)
			if err != nil {
				return nil, err
			}
		}

		if data != nil {
			mixData = append(mixData, *data)
		} else {
			mixData = append(mixData, mixing.Data{
				Data:       frame.details.Mix.NewMixBuffer(frame.details.Samples),
				PanMatrix:  frame.centerAheadPan,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: frame.details.Samples,
			})
		}
	}

	return mixData, nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) mixHardwareSynths(frame renderFrame, premix *output.PremixData) error {
	if len(m.hardwareSynths) == 0 {
		return nil
	}

	for _, synth := range m.hardwareSynths {
		data, adjust, err := synth.RenderTick(frame.centerAheadPan, frame.details)
		if err != nil {
			return err
		}

		premix.Data = append(premix.Data, mixing.ChannelData{data})

		if adjust != nil {
			premix.MixerVolume = adjust(premix.MixerVolume)
		}
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) normalizePremix(frame *renderFrame) {
	if len(frame.premix.Data) != 0 {
		return
	}

	frame.premix.Data = append(frame.premix.Data, mixing.ChannelData{
		mixing.Data{
			Data:       frame.details.Mix.NewMixBuffer(frame.details.Samples),
			PanMatrix:  frame.centerAheadPan,
			Volume:     volume.Volume(0),
			Pos:        0,
			SamplesLen: frame.details.Samples,
		},
	})
	return
}
