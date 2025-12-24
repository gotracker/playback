package machine

type stubRowStringer struct{ s string }

func (s stubRowStringer) String() string { return s.s }
