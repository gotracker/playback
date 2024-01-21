package song

import "github.com/gotracker/playback/index"

type rowIntf[TVolume Volume] interface {
	Len() int
	ForEach(fn func(ch index.Channel, d ChannelData[TVolume]) (bool, error)) error
}

// Row is a structure containing a single row
type Row any

func GetRowNumChannels[TVolume Volume](r Row) int {
	if row, ok := r.(rowIntf[TVolume]); ok {
		return row.Len()
	}
	return 0
}

func ForEachRowChannel[TVolume Volume](r Row, fn func(ch index.Channel, d ChannelData[TVolume]) (bool, error)) error {
	row, ok := r.(rowIntf[TVolume])
	if !ok {
		return nil
	}

	return row.ForEach(fn)
}

type RowStringer interface {
	String(options ...any) string
}
