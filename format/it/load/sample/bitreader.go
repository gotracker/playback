package sample

import (
	"bytes"
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
			r.bufSize, _ = r.Reader.Read(r.buffer[:])
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
