package crypto

import (
	"crypto/rand"
	"errors"
	"io"
)

var (
	InvalidSecret    = errors.New("invalid secret key")
	InvalidPublicKey = errors.New("invalid public key")
)

// SecretFrom seed reduces a seed to a secret key.
func SecretFromSeed(secret, seed *[32]byte) { reduce32(secret, seed) }

// PublicFromSecret generates a public key from a secret key.
func PublicFromSecret(public, secret *[32]byte) {
	var point geP3
	geScalarMultBase(&point, secret)
	geP3ToBytes(public, &point)
}

func CheckSecret(secret *[32]byte) bool { return scCheck(secret) }

func checkKey(key []byte) bool {
	var point geP3
	return geFromBytesVarTime(&point, key)
}

// newECScalar generates a new random ECScalar
func newECScalar() *[32]byte {
	tmp := make([]byte, 64)
	rand.Read(tmp)
	s := new([32]byte)
	scReduce(s[:], tmp)
	return s
}

func generateKeyDerivation(pub, sec *[32]byte) (*[32]byte, error) {
	var (
		point  geP3
		point2 geP2
		point3 geP1P1
	)

	if !scCheck(sec) {
		return nil, InvalidSecret
	}
	if !geFromBytesVarTime(&point, pub[:]) {
		return nil, InvalidPublicKey
	}

	geScalarMult(&point2, sec, &point)
	geMul8(&point3, &point2)
	geP1P1ToP2(&point2, &point3)

	d := new([32]byte)
	geToBytes(d, &point2)
	return d, nil
}

func hashToPoint(h []byte) *[32]byte {
	point := geFromFeFromBytesVarTime(h)
	b := new([32]byte)
	geToBytes(b, point)
	return b
}

// ViewFromSpend generates a view secret key from a spend secret key.
func ViewFromSpend(view, spend *[32]byte) {
	h := NewHash()
	h.Write(spend[:])
	h.Sum(view[:0])
	SecretFromSeed(view, view)
}

func GenerateSecret(random io.Reader) (sec [32]byte, err error) {
	tmp := make([]byte, 64)
	_, err = random.Read(tmp)
	scReduce(sec[:], tmp)
	return
}
