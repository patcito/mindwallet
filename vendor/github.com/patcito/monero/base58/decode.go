package base58

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/patcito/monero/crypto"
)

func decodeBlock(dst, src []byte) (int, error) {
	answer := big.NewInt(0)
	j := big.NewInt(1)
	for i := len(src) - 1; i >= 0; i-- {
		tmp := bytes.IndexByte(alphabet, src[i])
		if tmp == -1 {
			if src[i] == 0x00 {
				continue
			}
			return 0, fmt.Errorf("invalid character '%s' found", src[i])
		}
		idx := big.NewInt(int64(tmp))
		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		answer.Add(answer, tmp1)
		j.Mul(j, big.NewInt(alphabetSize))
	}

	l := blockSizes[len(src)]
	tmp := answer.Bytes()
	copy(dst[l-len(tmp):], tmp)
	return l, nil
}

func decodeAddr(s string) (tag uint64, data []byte) {
	b, err := DecodeString(s)
	if err != nil {
		return
	}

	if len(b) < checksumSize {
		return
	}

	checksum := b[len(b)-checksumSize:]
	b = b[:len(b)-checksumSize]
	hash := crypto.NewHash()
	hash.Write(b)
	digest := hash.Sum(nil)
	if !bytes.Equal(checksum, digest[:checksumSize]) {
		return
	}

	var n int
	tag, n = binary.Uvarint(b)
	data = b[n:]

	return
}

type decoder struct {
	err  error
	r    io.Reader
	buf  [fullEncodedBlockSize]byte
	nbuf int
	out  []byte //leftover decoded output
}

func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r, out: make([]byte, 0, fullBlockSize)}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Use leftover decode ouput from last read.
	if len(d.out) != 0 {
		n = copy(p, d.out)
		d.out = d.out[n:]
		return n, nil
	}
	if d.err != nil {
		return 0, d.err
	}

	var nn int
	// Read a block
	for len(p) >= fullBlockSize {
		nn, d.err = d.r.Read(d.buf[d.nbuf:fullEncodedBlockSize])
		if d.nbuf+nn != fullEncodedBlockSize {
			d.nbuf += nn
			break
		}
		_, d.err = decodeBlock(p, d.buf[:fullEncodedBlockSize])
		if d.err != nil {
			return n, d.err
		}
		d.nbuf = 0
		p = p[fullBlockSize:]
		n += fullBlockSize
	}
	if d.err == io.EOF {
		if d.nbuf != 0 {
			if len(p) >= fullEncodedBlockSize-d.nbuf {
				nn, d.err = decodeBlock(p, d.buf[:d.nbuf])
				d.nbuf = 0
				n += nn
				return n, d.err
			} else {
				nn, d.err = decodeBlock(d.out, d.buf[:d.nbuf])
			}
		}
	}
	return n, d.err
}

func DecodedLen(n int) int {
	return ((n / fullEncodedBlockSize) * fullBlockSize) + blockSizes[n%fullEncodedBlockSize]
}

func Decode(dst, src []byte) (n int, err error) {
	var nn int

	for len(src) >= fullEncodedBlockSize {
		nn, err = decodeBlock(dst, src[:fullEncodedBlockSize])
		n += nn
		if err != nil {
			return
		}
		dst = dst[nn:]
		src = src[fullEncodedBlockSize:]
	}

	if len(src) != 0 {
		nn, err = decodeBlock(dst, src)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func DecodeString(s string) ([]byte, error) {
	b := make([]byte, DecodedLen(len(s)))
	_, err := Decode(b, []byte(s))
	return b, err
}
