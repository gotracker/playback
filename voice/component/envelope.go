package component

// Envelope is an envelope component interface
type Envelope interface {
	//Init(env *envelope.Envelope)
	SetEnabled(enabled bool)
	IsEnabled() bool
	Reset()
	Advance(keyOn bool, prevKeyOn bool)
}
