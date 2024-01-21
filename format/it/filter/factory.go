package filter

import (
	"errors"
	"fmt"

	"github.com/gotracker/playback/filter"
	"github.com/gotracker/playback/frequency"
)

func Factory(name string, instrumentRate frequency.Frequency, params any) (filter.Filter, error) {
	var f filter.Filter
	switch name {
	case "":
		// nothing

	case "amigalpf":
		f = filter.NewAmigaLPF(instrumentRate)

	case "itresonant":
		p, ok := params.(filter.ITResonantFilterParams)
		if !ok {
			return nil, errors.New("could not convert it resonant filter parameters")
		}
		f = filter.NewITResonantFilter(p.Cutoff, p.Resonance, p.ExtendedFilterRange, p.Highpass)

	case "echo":
		p, ok := params.(filter.EchoFilterSettings)
		if !ok {
			return nil, errors.New("could not convert echo filter parameters")
		}
		f = &filter.EchoFilter{EchoFilterSettings: p}

	default:
		return nil, fmt.Errorf("unsupported filter name: %q", name)
	}

	return f, nil
}
