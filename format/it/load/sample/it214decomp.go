package sample

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Decompressor interface {
	Decompress(channels int) ([]byte, error)
}

type it214Decompressor struct {
	file   io.Reader
	sample *bytes.Buffer

	writtenSamples int

	isIT215 bool

	channels       int
	expectedLength int
	innerDecoder   func() error
}

func New(data []byte, sampleLen, channels int, is16Bit, isBigEndian bool, compatVer uint16) *it214Decompressor {
	decompress := it214Decompressor{
		file:           bytes.NewReader(data),
		isIT215:        compatVer >= 0x0215,
		channels:       channels,
		expectedLength: sampleLen,
	}
	if is16Bit {
		if isBigEndian {
			decompress.innerDecoder = decompress.blockDecoder16BE
		} else {
			decompress.innerDecoder = decompress.blockDecoder16LE
		}
	} else {
		decompress.innerDecoder = decompress.blockDecoder8
	}
	return &decompress
}

func (d *it214Decompressor) blockDecoder(blockDecoder blockDecoder, write blockDecoderWriteFunc) error {
	for chn := 0; chn < d.channels; chn++ {
		d.writtenSamples = 0

	blockReadLoop:
		for d.writtenSamples < d.expectedLength {
			var compressedSize uint16
			if err := binary.Read(d.file, binary.LittleEndian, &compressedSize); err != nil {
				break blockReadLoop
			}

			if compressedSize == 0 {
				// malformed sample?
				continue
			}

			block := make([]byte, compressedSize)
			if err := binary.Read(d.file, binary.LittleEndian, &block); err != nil {
				break blockReadLoop
			}

			r := bytes.NewReader(block)
			_ = blockDecoder.Decode(r, d, write)
		}
	}
	return nil
}

func (d *it214Decompressor) blockDecoder8() error {
	props := blockDecoder{
		blockSize: 0x8000,
		width:     9,
		fetchA:    3,
		lowerB:    -4,
		upperB:    3,
	}
	return d.blockDecoder(props, func(value int, w io.Writer) error {
		_, err := w.Write([]byte{uint8(value)})
		return err
	})
}

func (d *it214Decompressor) blockDecoder16(order binary.ByteOrder) error {
	props := blockDecoder{
		blockSize: 0x4000,
		width:     17,
		fetchA:    4,
		lowerB:    -8,
		upperB:    7,
	}
	return d.blockDecoder(props, func(value int, w io.Writer) error {
		var buf [2]byte
		order.PutUint16(buf[:], uint16(value))
		_, err := w.Write(buf[:])
		return err
	})
}

func (d *it214Decompressor) blockDecoder16BE() error {
	return d.blockDecoder16(binary.BigEndian)
}

func (d *it214Decompressor) blockDecoder16LE() error {
	return d.blockDecoder16(binary.LittleEndian)
}

func (d *it214Decompressor) Decompress() ([]byte, error) {
	d.sample = &bytes.Buffer{}

	if err := d.innerDecoder(); err != nil {
		return nil, err
	}

	return d.sample.Bytes(), nil
}
