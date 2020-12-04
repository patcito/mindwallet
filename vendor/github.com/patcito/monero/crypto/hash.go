package crypto

func hashToScalar(s *[32]byte, b []byte) {
	h := NewHash()
	h.Write(b)
	digest := make([]byte, 64)
	h.Sum(digest[:0])

	scReduce(s[:], digest)
}

func hashToEC(key *[32]byte) *geP3 {
	var (
		point2 geP1P1
	)
	r := new(geP3)
	h := NewHash()
	h.Write(key[:])
	digest := h.Sum(nil)
	point := geFromFeFromBytesVarTime(digest)

	geMul8(&point2, point)
	geP1P1ToP3(r, &point2)
	return r
}
