package machine

type Oscillator int

const (
	OscillatorVibrato = Oscillator(iota)
	OscillatorTremolo
	OscillatorPanbrello

	//====
	cNumOscillators
)

const NumOscillators = int(cNumOscillators)
