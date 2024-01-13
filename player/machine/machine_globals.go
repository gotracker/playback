package machine

import (
	"errors"
	"fmt"

	"github.com/gotracker/gomixing/volume"
	"github.com/gotracker/playback/index"
)

type globals[TGlobalVolume Volume] struct {
	bpm   int
	tempo int

	gv    TGlobalVolume // global volume
	mv    volume.Volume // mixing volume
	synv  volume.Volume // synth volume
	sampv volume.Volume // sample volume

	patternLoopStart index.Row
	patternLoopCount int
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetTempo(tempo int) error {
	if tempo == 0 {
		return errors.New("tempo cannot be 0")
	}

	traceValueChangeWithComment(m, "tempo", m.tempo, tempo, "SetTempo")
	m.tempo = tempo

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetBPM(bpm int) error {
	if bpm == 0 {
		return errors.New("bpm cannot be 0")
	}

	traceValueChangeWithComment(m, "bpm", m.bpm, bpm, "SetBPM")
	m.bpm = bpm

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SlideBPM(add int) error {
	if add == 0 {
		return nil
	}

	bpm := m.bpm + add
	if bpm <= 0 || bpm > 255 {
		return fmt.Errorf("resulting bpm would be invalid: %d", bpm)
	}

	traceValueChangeWithComment(m, "bpm", m.bpm, bpm, "SlideBPM")
	m.bpm = bpm

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetGlobalVolume(v TGlobalVolume) error {
	if v.IsInvalid() {
		return fmt.Errorf("global volume out of range: %v", v)
	}

	traceValueChangeWithComment(m, "gv", m.gv, v, "SetGlobalVolume")
	m.gv = v

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SlideGlobalVolume(multiplier, add float32) error {
	fma, ok := any(m.gv).(VolumeFMA[TGlobalVolume])
	if !ok {
		return errors.New("could not determine FMA interface for global volume")
	}
	v := fma.FMA(multiplier, add)

	if v.IsInvalid() {
		return fmt.Errorf("global volume out of range: %v", v)
	}

	traceValueChangeWithComment(m, "gv", m.gv, v, "SlideGlobalVolume")
	m.gv = v
	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetMixingVolume(v volume.Volume) error {
	if v < 0 || v > 1 {
		return fmt.Errorf("mixing volume out of range: %v", v)
	}

	traceValueChangeWithComment(m, "mv", m.mv, v, "SetMixingVolume")
	m.mv = v

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetSynthVolume(v volume.Volume) error {
	if v < 0 || v > 1 {
		return fmt.Errorf("synth volume out of range: %v", v)
	}

	traceValueChangeWithComment(m, "synv", m.synv, v, "SetSynthVolume")
	m.synv = v

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetSampleVolume(v volume.Volume) error {
	if v < 0 || v > 1 {
		return fmt.Errorf("sample volume out of range: %v", v)
	}

	traceValueChangeWithComment(m, "sampv", m.sampv, v, "SetSampleVolume")
	m.sampv = v

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetOrder(o index.Order) error {
	if int(o) >= len(m.songData.GetOrderList()) {
		return fmt.Errorf("order index out of range: %d", o)
	}

	traceOptionalValueChangeWithComment(m, "nextOrder", m.ticker.next.order, o, "SetOrder")
	m.ticker.next.order.Set(o)

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetRow(r index.Row, breakOrder bool) error {
	traceOptionalValueChangeWithComment(m, "nextRow", m.ticker.next.row, r, "SetRow")
	m.ticker.next.row.Set(r)

	return nil
}

func (m *machine[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) SetFilterOnAllChannelsByFilterName(name string, enabled bool) error {
	filtFactory, err := m.ms.GetFilterFactory(name)
	if err != nil {
		return err
	}

	cr := m.songData.GetSystem().GetCommonRate()

	for i := range m.channels {
		c := &m.channels[i]
		if enabled {
			filt, err := filtFactory(cr)
			if err != nil {
				return err
			}
			c.filter = filt
		} else {
			c.filter = nil
		}
	}

	return nil
}
