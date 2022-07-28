package load

import (
	"bytes"
	"encoding/binary"
	"io"
)

type itBitReader struct {
	*bytes.Reader
	bitbuf uint
	bitnum uint32
}

func (r *itBitReader) ReadBits(n uint) (uint, error) {
	var value uint

	// this could be better
	for n > 0 {
		n--
		if r.bitnum == 0 {
			b, err := r.ReadByte()
			if err != nil {
				return value >> (32 - n), err
			}
			r.bitbuf = uint(b)
			r.bitnum = 8
		}
		value >>= 1
		value |= r.bitbuf << 31
		r.bitbuf >>= 1
		r.bitnum--
	}
	return value >> (32 - n), nil
}

type itSampleDecompress struct {
	file    io.Reader
	bitFile *itBitReader
	sample  *bytes.Buffer

	writtenSamples int
	curLength      int

	mem1 int
	mem2 int

	isIT215 bool

	expectedLength int
	innerDecoder   func(channels int) error
}

const (
	itSampleDecompress8BlockSize = 0x8000
	itSampleDecompress8Width     = 9
	itSampleDecompress8FetchA    = 0
)

type itSampleDecompressProps struct {
	blockSize int
	width     int
	fetchA    int
	lowerB    int
	upperB    int
}

func (d *itSampleDecompress) blockDecoder(channels int, props itSampleDecompressProps, write func(v, topBit int, target io.Writer)) error {
	for chn := 0; chn < channels; chn++ {
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

			d.bitFile = &itBitReader{
				Reader: bytes.NewReader(block),
			}

			// Initialise bit reader
			d.mem1 = 0
			d.mem2 = 0

			d.curLength = d.expectedLength - d.writtenSamples
			if d.curLength > itSampleDecompress8BlockSize {
				d.curLength = itSampleDecompress8BlockSize
			}

			width := props.width
			for d.curLength > 0 {
				if width > props.width {
					// error
					continue blockReadLoop
				}

				bv, err := d.bitFile.ReadBits(uint(width))
				if err != nil {
					continue blockReadLoop
				}

				v := int(bv)

				topBit := 1 << (width - 1)
				if width <= 6 {
					// mode A :: 1..6 bits
					if v == topBit {
						newWidth := props.fetchA + 1
						if newWidth >= width {
							newWidth++
						}
						width = newWidth
					} else {
						write(v, topBit, d.sample)
					}
				} else if width < props.width {
					// mode B :: 7..8 bits [16bit = 16]
					if v >= topBit+props.lowerB && v < topBit+props.upperB {
						newWidth := v - (topBit + props.lowerB) + 1
						if newWidth >= width {
							newWidth++
						}
						width = newWidth
					} else {
						write(v, topBit, d.sample)
					}
				} else {
					// mode C :: 9 bits [16bit = 17]
					if (v & topBit) != 0 {
						width = (v &^ topBit) + 1
					} else {
						write(v&^topBit, 0, d.sample)
					}
				}
			}
		}
	}
	return nil
}

func (d *itSampleDecompress) blockDecoder8(channels int) error {
	props := itSampleDecompressProps{
		blockSize: 0x8000,
		width:     9,
		fetchA:    3,
		lowerB:    -4,
		upperB:    3,
	}
	return d.blockDecoder(channels, props, func(v, topBit int, target io.Writer) {
		if (v & topBit) != 0 {
			v -= topBit << 1
		}
		d.mem1 += v
		d.mem2 += d.mem1
		if d.isIT215 {
			value := int16(d.mem2)
			_ = binary.Write(target, binary.LittleEndian, value)
		} else {
			value := int16(d.mem1)
			_ = binary.Write(target, binary.LittleEndian, value)
		}
		d.writtenSamples++
		d.curLength--
	})
}

func (d *itSampleDecompress) blockDecoder16(channels int, order binary.ByteOrder) error {
	props := itSampleDecompressProps{
		blockSize: 0x4000,
		width:     17,
		fetchA:    4,
		lowerB:    -8,
		upperB:    7,
	}
	return d.blockDecoder(channels, props, func(v, topBit int, target io.Writer) {
		if (v & topBit) != 0 {
			v -= topBit << 1
		}
		d.mem1 += v
		d.mem2 += d.mem1
		if d.isIT215 {
			value := int16(d.mem2)
			_ = binary.Write(target, order, value)
		} else {
			value := int16(d.mem1)
			_ = binary.Write(target, order, value)
		}
		d.writtenSamples++
		d.curLength--
	})
}

func (d *itSampleDecompress) blockDecoder16BE(channels int) error {
	return d.blockDecoder16(channels, binary.BigEndian)
}

func (d *itSampleDecompress) blockDecoder16LE(channels int) error {
	return d.blockDecoder16(channels, binary.LittleEndian)
}

func (d *itSampleDecompress) Decompress(channels int) ([]byte, error) {
	d.sample = &bytes.Buffer{}

	if err := d.innerDecoder(channels); err != nil {
		return nil, err
	}

	return d.sample.Bytes(), nil
}

func newITSampleDecompress(data []byte, sampleLen int, is16Bit, isBigEndian bool, compatVer uint16) *itSampleDecompress {
	decompress := itSampleDecompress{
		file:           bytes.NewReader(data),
		expectedLength: sampleLen,
		isIT215:        compatVer >= 0x0215,
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
