package period

// Period is an interface that defines a sampler period
type Period interface {
	IsInvalid() bool
}
