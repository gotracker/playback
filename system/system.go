package system

type System interface {
	GetMaxPastNotesPerChannel() int
	GetCommonRate() Frequency
}
