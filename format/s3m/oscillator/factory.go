package oscillator

import (
	"fmt"

	oscillatorImpl "github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/voice/oscillator"
)

func VibratoFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}

func TremoloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}

func PanbrelloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewProtrackerOscillator(), nil
}

func OscillatorFactory(name string) (oscillator.Oscillator, error) {
	switch name {
	case "":
		return nil, nil
	case "vibrato":
		return VibratoFactory()
	case "tremolo":
		return TremoloFactory()
	case "panbrello":
		return PanbrelloFactory()
	default:
		return nil, fmt.Errorf("unsupported oscillator: %q", name)
	}
}
