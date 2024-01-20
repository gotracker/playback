package filter

import (
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

	default:
		return nil, fmt.Errorf("unsupported filter name: %q", name)
	}

	return f, nil
}
