package player

// PostTickable is an interface which exposes the OnPostTick call
type PostTickable interface {
	OnPostTick() error
}

// DoPostTick calls the OnPostTick() function on the interface, if possible
func DoPostTick(t PostTickable) error {
	if t != nil {
		return t.OnPostTick()
	}
	return nil
}
