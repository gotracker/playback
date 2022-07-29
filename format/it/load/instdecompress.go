package load

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type itBitReader struct {
	*bytes.Reader
	bitNum  uint
	bitBuf  uint32
	bufPos  int
	bufSize int
	buffer  [8]byte
}

func (r *itBitReader) ReadBits(n uint) (uint, error) {
	for r.bitNum < n {
		// Fetch more bits
		if r.bufPos >= r.bufSize {
			var err error
			r.bufSize, err = r.Reader.Read(r.buffer[:])
			if err != nil {
				return 0, err
			}
			r.bufPos = 0
			if r.bufSize == 0 {
				return 0, io.ErrUnexpectedEOF
			}
		}
		val := r.buffer[r.bufPos]
		r.bufPos++
		r.bitBuf |= uint32(val) << r.bitNum
		r.bitNum += 8
	}

	v := r.bitBuf & ((1 << n) - 1)
	r.bitBuf >>= n
	r.bitNum -= n
	return uint(v), nil
}

type itSampleDecompress struct {
	file   io.Reader
	sample *bytes.Buffer

	writtenSamples int

	isIT215 bool

	expectedLength int
	innerDecoder   func(channels int) error
}

type itSampleBlockDecoder struct {
	blockSize int
	width     int
	fetchA    int
	lowerB    int
	upperB    int
}

type itSampleBlockDecoderWriteFunc func(value int, w io.Writer) error

func (d itSampleBlockDecoder) writeSample(v, topBit int, mem1, mem2 *int, out *itSampleDecompress, fn itSampleBlockDecoderWriteFunc) error {
	if (v & topBit) != 0 {
		v -= topBit << 1
	}

	*mem1 += v
	*mem2 += *mem1

	val := mem1
	if out.isIT215 {
		val = mem2
	}

	err := fn(*val, out.sample)
	if err != nil {
		return err
	}

	out.writtenSamples++
	return nil
}

var ErrBitWidthOutOfRange = errors.New("bit width out of range")

func (d itSampleBlockDecoder) Decode(r *bytes.Reader, out *itSampleDecompress, fn itSampleBlockDecoderWriteFunc) error {
	br := itBitReader{
		Reader: r,
	}

	// Initialise bit reader
	var (
		mem1      int
		mem2      int
		curLength = out.expectedLength - out.writtenSamples
	)

	if curLength > d.blockSize {
		curLength = d.blockSize
	}

	width := d.width
	for curLength > 0 {
		if width > d.width {
			// error
			return ErrBitWidthOutOfRange
		}

		bv, err := br.ReadBits(uint(width))
		if err != nil {
			return err
		}

		v := int(bv)

		topBit := 1 << (width - 1)
		if width <= 6 {
			// mode A :: 1..6 bits
			if v == topBit {
				newWidth := d.fetchA + 1
				if newWidth >= width {
					newWidth++
				}
				width = newWidth
			} else {
				if err := d.writeSample(v, topBit, &mem1, &mem2, out, fn); err != nil {
					return err
				}
				curLength--
			}
		} else if width < d.width {
			// mode B :: 7..8 bits [16bit = 16]
			if v >= topBit+d.lowerB && v < topBit+d.upperB {
				newWidth := v - (topBit + d.lowerB) + 1
				if newWidth >= width {
					newWidth++
				}
				width = newWidth
			} else {
				if err := d.writeSample(v, topBit, &mem1, &mem2, out, fn); err != nil {
					return err
				}
				curLength--
			}
		} else {
			// mode C :: 9 bits [16bit = 17]
			if (v & topBit) != 0 {
				width = (v &^ topBit) + 1
			} else {
				if err := d.writeSample(v&^topBit, 0, &mem1, &mem2, out, fn); err != nil {
					return err
				}
				curLength--
			}
		}
	}
	return nil
}

func (d *itSampleDecompress) blockDecoder(channels int, blockDecoder itSampleBlockDecoder, write itSampleBlockDecoderWriteFunc) error {
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

			r := bytes.NewReader(block)
			_ = blockDecoder.Decode(r, d, write)
		}
	}
	return nil
}

func (d *itSampleDecompress) blockDecoder8(channels int) error {
	props := itSampleBlockDecoder{
		blockSize: 0x8000,
		width:     9,
		fetchA:    3,
		lowerB:    -4,
		upperB:    3,
	}
	return d.blockDecoder(channels, props, func(value int, w io.Writer) error {
		_, err := w.Write([]byte{uint8(value)})
		return err
	})
}

func (d *itSampleDecompress) blockDecoder16(channels int, order binary.ByteOrder) error {
	props := itSampleBlockDecoder{
		blockSize: 0x4000,
		width:     17,
		fetchA:    4,
		lowerB:    -8,
		upperB:    7,
	}
	return d.blockDecoder(channels, props, func(value int, w io.Writer) error {
		var buf [2]byte
		order.PutUint16(buf[:], uint16(value))
		_, err := w.Write(buf[:])
		return err
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
