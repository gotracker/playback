package voice

// FilterEnveloper is a filter envelope interface
type FilterEnveloper interface {
	EnableFilterEnvelope(enabled bool)
	IsFilterEnvelopeEnabled() bool
	GetCurrentFilterEnvelope() int8
	SetFilterEnvelopePosition(pos int)
}
