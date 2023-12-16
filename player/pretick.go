package player

// PreTickable is an interface which exposes the OnPreTick call
type PreTickable interface {
	OnPreTick() error
}

// DoPreTick calls the OnPreTick() function on the interface, if possible
func DoPreTick(t PreTickable) error {
	if t != nil {
		return t.OnPreTick()
	}
	return nil
}
