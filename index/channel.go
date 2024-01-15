package index

type Channel uint8

const (
	InvalidOPLChannel = OPLChannel(0xFF)
)

type OPLChannel uint8

func (c OPLChannel) IsValid() bool {
	return c != InvalidOPLChannel
}
