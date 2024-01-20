package oscillator

import (
	"fmt"

	oscillatorImpl "github.com/gotracker/playback/oscillator"
	"github.com/gotracker/playback/voice/oscillator"
)

func VibratoFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(4), nil
}

func TremoloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(4), nil
}

func PanbrelloFactory() (oscillator.Oscillator, error) {
	return oscillatorImpl.NewImpulseTrackerOscillator(1), nil
}

func OscillatorFactory(name string) (oscillator.Oscillator, error) {
	switch name {
	case "":
		return nil, nil
	case "vibrato":
		return VibratoFactory()
	case "autovibrato":
		return oscillatorImpl.NewImpulseTrackerOscillator(1), nil
	case "tremolo":
		return TremoloFactory()
	case "panbrello":
		return PanbrelloFactory()
	default:
		return nil, fmt.Errorf("unsupported oscillator: %q", name)
	}
}
