package player

// Tickable is an interface which exposes the OnTick call
type Tickable interface {
	OnTick() error
}

// DoTick calls the OnTick() function on the interface, if possible
func DoTick(t Tickable) error {
	if t != nil {
		return t.OnTick()
	}
	return nil
}
