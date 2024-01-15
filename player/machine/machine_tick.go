package machine

import (
	"errors"

	"github.com/gotracker/playback/frequency"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/sampler"
	"github.com/gotracker/playback/song"
)

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) Tick(s *sampler.Sampler) (*output.PremixData, error) {
	for i := range m.channels {
		if err := m.channels[i].DoNoteAction(index.Channel(i), m, frequency.Frequency(s.SampleRate)); err != nil {
			return nil, err
		}
	}

	premix, err := m.render(s)
	if err != nil {
		return premix, err
	}

	tickErr := runTick(&m.ticker, m)
	if tickErr != nil {
		if !errors.Is(tickErr, song.ErrStopSong) {
			return nil, tickErr
		}
	}

	m.age++
	return premix, errors.Join(tickErr, err)
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) onTick() error {
	for i := range m.channels {
		c := &m.channels[i]
		if err := c.Tick(index.Channel(i), m); err != nil {
			return err
		}

		c.updatePastNotes()
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
