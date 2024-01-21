package index

const (
	InvalidChannel    = Channel(0xFFFF)
	InvalidOPLChannel = OPLChannel(0xFF)
)

type Channel uint16

func (c Channel) IsValid() bool {
	return c != InvalidChannel
}

type OPLChannel uint8

func (c OPLChannel) IsValid() bool {
	return c != InvalidOPLChannel
}
