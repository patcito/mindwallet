package base58

import "math/big"

const (
	fullBlockSize        = 8
	fullEncodedBlockSize = 11

	alphabetSize = 58
	checksumSize = 4
)

var (
	blockSizes        = [...]int{0, 0, 1, 2, 3, 3, 4, 5, 6, 6, 7, 8}
	encodedBlockSizes = [...]int{0, 2, 3, 5, 6, 7, 9, 10, 11}
	alphabet          = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	radix             = big.NewInt(alphabetSize)
	bigZero           = big.NewInt(0)
)
