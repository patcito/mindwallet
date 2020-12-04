package base58

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/patcito/monero/crypto"
)

func encodeBlock(dst, src []byte) {
	var num uint64
	var i int
	if len(src) > fullBlockSize {
		num = uint8beTo64(src[:fullBlockSize])
		i = fullEncodedBlockSize
	} else {
		num = uint8beTo64(src)
		i = encodedBlockSizes[len(src)]
	}
	i--
	for i >= 0 {
		remainder := num % alphabetSize
		num /= alphabetSize
		dst[i] = alphabet[remainder]
		i--
	}
}

func uint8beTo64(b []byte) (r uint64) {
	i := 0
	switch 9 - len(b) {
	case 1:
		r |= uint64(b[i])
		i++
		fallthrough
	case 2:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 3:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 4:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 5:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 6:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 7:
		r <<= 8
		r |= uint64(b[i])
		i++
		fallthrough
	case 8:
		r <<= 8
		r |= uint64(b[i])
	}
	return r
}

func encodeAddr(tag uint64, data []byte) string {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	hash := crypto.NewHash()
	mw := io.MultiWriter(enc, hash)

	//binary.Write(mw, binary.LittleEndian, tag)
	b := make([]byte, 16)
	n := binary.PutUvarint(b, tag)
	mw.Write(b[:n])
	mw.Write(data)

	digest := hash.Sum(nil)
	enc.Write(digest[:checksumSize])
	enc.Close()
	return buf.String()
}

type encoder struct {
	err  error
	w    io.Writer
	buf  [fullBlockSize]byte
	nbuf int
	out  [fullEncodedBlockSize]byte
}

func NewEncoder(w io.Writer) io.WriteCloser {
	return &encoder{w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.err != nil {
		return 0, e.err
	}

	// Leading fringe.
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < fullBlockSize; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < fullBlockSize {
			return
		}
		encodeBlock(e.out[0:], e.buf[0:])
		if _, e.err = e.w.Write(e.out[0:]); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	// Full interior blocks.
	for len(p) >= fullBlockSize {
		encodeBlock(e.out[0:], p[0:fullBlockSize])
		if _, e.err = e.w.Write(e.out[0:]); err != nil {
			return n, err
		}
		n += fullBlockSize
		p = p[fullBlockSize:]
	}

	// Trailing fringe.
	for i := 0; i < len(p); i++ {
		e.buf[i] = p[i]
	}
	e.nbuf = len(p)
	n += len(p)
	return
}

func (e *encoder) Close() error {
	if e.err == nil && e.nbuf > 0 {
		encodeBlock(e.out[0:], e.buf[0:e.nbuf])
		_, e.err = e.w.Write(e.out[0:encodedBlockSizes[e.nbuf]])
		e.nbuf = 0
	}
	return e.err
}

func EncodedLen(n int) int {
	return ((n / fullBlockSize) * fullEncodedBlockSize) + encodedBlockSizes[n%fullBlockSize]
}

func Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}
	fullBlockCount := len(src) / fullBlockSize
	lastBlockSize := len(src) % fullBlockSize
	//if len(dst) < fullBlockCount {
	//	panic("Encode: dst is too small to hold src")
	//}

	for i := 0; i < fullBlockCount; i++ {
		encodeBlock(dst[i*fullEncodedBlockSize:], src[i*fullBlockSize:])
	}
	if lastBlockSize > 0 {
		encodeBlock(dst[fullBlockCount*fullEncodedBlockSize:], src[fullBlockCount*fullBlockSize:])
	}
}

func EncodeToString(src []byte) string {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.Write(src)
	enc.Close()
	return buf.String()
}
