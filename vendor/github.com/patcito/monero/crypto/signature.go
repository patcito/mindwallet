package crypto

//import "crypto/rand"

type Signature struct {
	c, r [32]byte
}

/*
func generateSignature(prefixHash []byte, public *PublicKey, secret *SecretKey) *Signature {
	var (
		tmp3 geP3
		k    ECScalar
	)
	l := len(prefixHash)
	buf := make([]byte, l+64)

	copy(buf, prefixHash)
	copy(buf[l:], public[:])

	rand.Read(k[:])

	geScalarMultBase(&tmp3, &k)
	copy(buf[l+32:], geP3ToBytes(&tmp3)[:])

	sig := new(Signature)
	hashToScalar(&sig.c, buf)
	scMulSub(&sig.r, &sig.c, secret[:], k[:])
	return sig
}
*/

func checkSignature(prefixHash []byte, pub *[32]byte, sig []byte) bool {
	var (
		tmp2 geP2
		tmp3 geP3
		c    [32]byte
	)

	buf := make([]byte, 96)
	copy(buf[:32], prefixHash)
	copy(buf[32:64], pub[:])

	if !geFromBytesVarTime(&tmp3, pub[:]) {
		return false
	}

	// still need a consistant way to pass arrays around
	var sigC, sigR [32]byte
	copy(sigC[:], sig[:32])
	copy(sigR[:], sig[32:])

	if !scCheck(&sigC) || !scCheck(&sigR) {
		return false
	}
	geDoubleScalarMultBaseVarTime(&tmp2, &sigC, &tmp3, &sigR)
	var b [32]byte
	geToBytes(&b, &tmp2)
	copy(buf[64:], b[:])
	hashToScalar(&c, buf)
	scSub(&c, &c, &sigC)
	return !scIsNonZero(&c)
}

func generateKeyImage(public, secret *[32]byte) *[32]byte {
	var point2 geP2

	point := hashToEC(public)
	geScalarMult(&point2, secret, point)
	b := new([32]byte)
	geToBytes(b, &point2)
	return b
}

type ringSignature struct {
	hash, a, b [32]byte
}

func checkRingSignature(prefixHash, image []byte, pubs []*[32]byte, sig []byte) bool {
	var (
		imageUnp geP3
		imagePre geDsmp
		sum, h   [32]byte
	)

	if !geFromBytesVarTime(&imageUnp, image) {
		return false
	}

	geDsmPrecomp(&imagePre, &imageUnp)

	//if (len(sig) % 64) != 0 {
	//	return true
	//}
	sigs := make([]*Signature, len(sig)/64)
	j := 0
	k := 32
	for i := 0; i < len(sigs); i++ {
		s := new(Signature)
		copy(s.c[:], sig[j:k])
		j += 32
		k += 32
		copy(s.r[:], sig[j:k])
		sigs[i] = s
		j += 32
		k += 32
	}

	buf := make([]byte, 32+len(sig))
	copy(buf, prefixHash)
	j = 32
	k = 64

	var b [32]byte
	for i := 0; i < len(pubs); i++ {
		var (
			tmp2 geP2
			tmp3 geP3
		)
		if !scCheck(&sigs[i].c) || !scCheck(&sigs[i].r) {
			return false
		}

		if !geFromBytesVarTime(&tmp3, pubs[i][:]) {
			panic("abort()")
		}

		geDoubleScalarMultBaseVarTime(&tmp2, &sigs[i].c, &tmp3, &sigs[i].r)

		geToBytes(&b, &tmp2)
		copy(buf[j:k], b[:])
		j += 32
		k += 32
		tmp3 = *hashToEC(pubs[i])

		geDoubleScalarMultPrecompVarTime(&tmp2, &sigs[i].r, &tmp3, &sigs[i].c, &imagePre)
		geToBytes(&b, &tmp2)
		copy(buf[j:k], b[:])
		j += 32
		k += 32
		scAdd(&sum, &sum, &sigs[i].c)
	}

	hashToScalar(&h, buf)
	scSub(&h, &h, &sum)
	return !scIsNonZero(&h)
}
