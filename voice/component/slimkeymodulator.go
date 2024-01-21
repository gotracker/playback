package component

type slimKeyModulator struct {
	keyOn     bool
	prevKeyOn bool
}

func (k slimKeyModulator) IsKeyOn() bool {
	return k.keyOn
}

func (k slimKeyModulator) WasKeyOn() bool {
	return k.prevKeyOn
}

func (k *slimKeyModulator) Attack() {
	k.keyOn, k.prevKeyOn = true, k.keyOn
}

func (k *slimKeyModulator) Release() {
	k.keyOn, k.prevKeyOn = false, k.keyOn
}
