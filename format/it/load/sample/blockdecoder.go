package sample

import (
	"bytes"
	"errors"
	"io"
)

type blockDecoder struct {
	blockSize int
	width     int
	fetchA    int
	lowerB    int
	upperB    int
}

type blockDecoderWriteFunc func(value int, w io.Writer) error

func (d blockDecoder) writeSample(v, topBit int, mem1, mem2 *int, out *it214Decompressor, fn blockDecoderWriteFunc) error {
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

func (d blockDecoder) Decode(r *bytes.Reader, out *it214Decompressor, fn blockDecoderWriteFunc) error {
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
