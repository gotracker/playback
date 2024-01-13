package machine

import (
	"errors"
	"fmt"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/render"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
	"github.com/gotracker/playback/voice"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Tick(s *sampler.Sampler) (*output.PremixData, error) {
	m.getSampleRate = func() period.Frequency {
		return period.Frequency(s.SampleRate)
	}

	for i := range m.channels {
		if err := m.channels[i].DoNoteAction(index.Channel(i), m); err != nil {
			return nil, err
		}
	}

	tickDuration := m.songData.GetTickDuration(m.bpm)
	if tickDuration <= 0 {
		return nil, fmt.Errorf("unexpected tick duration: %v", tickDuration)
	}

	renderRow := render.RowRender{
		Order: int(m.ticker.current.order),
		Row:   int(m.ticker.current.row),
		Tick:  m.ticker.current.tick,
	}

	premix := output.PremixData{
		SamplesLen:  int(float64(s.SampleRate) * tickDuration.Seconds()),
		MixerVolume: m.gv.ToVolume() * m.mv,
		Userdata:    &renderRow,
	}

	if m.ticker.current.tick == 0 {
		// make a copy so it doesn't get stomped
		renderRow.RowText = m.rowStringer
	}

	tickErr := runTick(&m.ticker, m)
	if tickErr != nil {
		if !errors.Is(tickErr, song.ErrStopSong) {
			return nil, tickErr
		}
	}

	sys := m.songData.GetSystem()

	details := render.Details{
		Mix:          s.Mixer(),
		Panmixer:     s.GetPanMixer(),
		SamplerSpeed: sys.GetSamplerSpeed(period.Frequency(s.SampleRate)),
		Samples:      premix.SamplesLen,
		Duration:     tickDuration,
	}

	centerAheadPan := details.Panmixer.GetMixingMatrix(panning.CenterAhead)

	for i := range m.channels {
		ch := index.Channel(i)
		c := &m.channels[ch]
		var mixData []mixing.Data

		if pos, set := c.target.Pos.Get(); set {
			if samp, ok := c.cv.(voice.Sampler); ok {
				samp.SetPos(pos)
			}
			c.target.Pos.Reset()
		}

		rc := &m.outputChannels[ch]

		c.cv.DumpState(ch, m.us.Tracer)
		data, err := c.cv.Render(centerAheadPan, details, rc)
		if err != nil {
			return nil, errors.Join(tickErr, err)
		}
		c.cv.Advance()
		if data != nil {
			mixData = append(mixData, *data)
		}

		pnData, err := c.pn.Render(centerAheadPan, details, rc)
		if err != nil {
			return nil, errors.Join(tickErr, err)
		}
		if len(pnData) > 0 {
			mixData = append(mixData, pnData...)
		}

		if len(mixData) > 0 {
			premix.Data = append(premix.Data, mixData)
		}
	}

	if m.opl2 != nil {
		rr := [1]mixing.Data{}
		if err := m.renderOPL2Tick(&rr[0], s.Mixer(), premix.SamplesLen); err != nil {
			return nil, errors.Join(tickErr, err)
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
				Pan:        panning.CenterAhead,
				Volume:     volume.Volume(0),
				Pos:        0,
				SamplesLen: details.Samples,
			},
		})
	}

	m.age++
	return &premix, tickErr
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onTick() error {
	for i := range m.channels {
		c := &m.channels[i]
		if err := c.Tick(index.Channel(i), m); err != nil {
			return err
		}

		c.pn.UpdatePastNotes()
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onOrderStart() error {
	for ch := range m.channels {
		if err := m.channels[ch].OrderStart(index.Channel(ch), m); err != nil {
			return err
		}
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onRowStart() error {
	rowData, err := m.getRowData()
	if err != nil {
		return err
	}

	m.rowStringer = m.songData.GetRowRenderStringer(rowData, len(m.channels), m.us.LongChannelOutput)

	trace(m, m.rowStringer.String())

	if err := m.singleRowRowStart(); err != nil {
		return err
	}

	if err := m.updateInstructions(rowData); err != nil {
		return err
	}

	for ch := range m.channels {
		if err := m.channels[ch].RowStart(index.Channel(ch), m); err != nil {
			return err
		}
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onRowEnd() error {
	for ch := range m.channels {
		if err := m.channels[ch].RowEnd(index.Channel(ch), m); err != nil {
			return err
		}
	}

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onOrderEnd() error {
	for ch := range m.channels {
		if err := m.channels[ch].OrderEnd(index.Channel(ch), m); err != nil {
			return err
		}
	}

	return nil
}
