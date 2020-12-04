package crypto

import "encoding/hex"

type geP2 struct {
	X, Y, Z fe
}

type geP3 struct {
	X, Y, Z, T fe
}

func (p *geP3) String() string {
	var b [32]byte
	geP3ToBytes(&b, p)
	return hex.EncodeToString(b[:])
}

type geP1P1 struct {
	X, Y, Z, T fe
}

type geP1P2 struct {
	X, Y, Z, T fe
}

type gePrecomp struct {
	yplusx, yminusx, xy2d fe
}

type geCached struct {
	YplusX, YminusX, Z, T2d fe
}

type geDsmp [8]geCached

func geMul8(r *geP1P1, t *geP2) {
	var u geP2
	geP2Dbl(r, t)
	geP1P1ToP2(&u, r)
	geP2Dbl(r, &u)
	geP1P1ToP2(&u, r)
	geP2Dbl(r, &u)
}

func geFromBytesVarTime(h *geP3, s []byte) bool {
	var u, v, vxx, check fe

	// From fe_frombytes.c
	h0 := int64(load4(s))
	h1 := int64(load3(s[4:]) << 6)
	h2 := int64(load3(s[7:]) << 5)
	h3 := int64(load3(s[10:]) << 3)
	h4 := int64(load3(s[13:]) << 2)
	h5 := int64(load4(s[16:]))
	h6 := int64(load3(s[20:]) << 7)
	h7 := int64(load3(s[23:]) << 5)
	h8 := int64(load3(s[26:]) << 4)
	h9 := int64((load3(s[29:]) & 0x7fffff) << 2)
	var (
		carry0, carry1, carry2, carry3, carry4 int64
		carry5, carry6, carry7, carry8, carry9 int64
	)

	// Validate the number to be canonical
	if h9 == 33554428 && h8 == 268435440 && h7 == 536870880 && h6 == 2147483520 && h5 == 4294967295 && h4 == 67108860 && h3 == 134217720 && h2 == 536870880 && h1 == 1073741760 && h0 >= 4294967277 {
		return false
	}

	carry9 = (h9 + (1 << 24)) >> 25
	h0 += carry9 * 19
	h9 -= carry9 << 25
	carry1 = (h1 + (1 << 24)) >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry3 = (h3 + (1 << 24)) >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry5 = (h5 + (1 << 24)) >> 25
	h6 += carry5
	h5 -= carry5 << 25
	carry7 = (h7 + (1 << 24)) >> 25
	h8 += carry7
	h7 -= carry7 << 25

	carry0 = (h0 + (1 << 25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry2 = (h2 + (1 << 25)) >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry4 = (h4 + (1 << 25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry6 = (h6 + (1 << 25)) >> 26
	h7 += carry6
	h6 -= carry6 << 26
	carry8 = (h8 + (1 << 25)) >> 26
	h9 += carry8
	h8 -= carry8 << 26

	h.Y[0] = int32(h0)
	h.Y[1] = int32(h1)
	h.Y[2] = int32(h2)
	h.Y[3] = int32(h3)
	h.Y[4] = int32(h4)
	h.Y[5] = int32(h5)
	h.Y[6] = int32(h6)
	h.Y[7] = int32(h7)
	h.Y[8] = int32(h8)
	h.Y[9] = int32(h9)

	// End fe_frombytes.c

	fe1(&h.Z)
	feSq(&u, &h.Y)
	feMul(&v, &u, feD)
	feSub(&u, &u, &h.Z) // u = y^2-1
	feAdd(&v, &v, &h.Z) // v = dy^2+1

	feDivPowM1(&h.X, &u, &v) // x = uv^3(uv^7)^((q-5)/8)

	feSq(&vxx, &h.X)
	feMul(&vxx, &vxx, &v)
	feSub(&check, &vxx, &u) // vx^2-u
	if feIsNonZero(&check) {
		feAdd(&check, &vxx, &u) // vx^2+u
		if feIsNonZero(&check) {
			return false
		}
		feMul(&h.X, &h.X, fe_sqrtm1)
	}

	if feIsNegative(&h.X) != (s[31] >> 7) {
		/* If x = 0, the sign must be positive */
		if !feIsNonZero(&h.X) {
			return false
		}
		feNeg(&h.X, &h.X)
	}

	feMul(&h.T, &h.X, &h.Y)
	return true

}

func geP20(h *geP2) {
	fe0(&h.X)
	fe1(&h.Y)
	fe1(&h.Z)
}

func geP30(h *geP3) {
	fe0(&h.X)
	fe1(&h.Y)
	fe1(&h.Z)
	fe0(&h.T)
}

func gePrecomp0(h *gePrecomp) {
	fe1(&h.yplusx)
	fe1(&h.yminusx)
	fe0(&h.xy2d)
}

func geCached0(r *geCached) {
	fe1(&r.YplusX)
	fe1(&r.YminusX)
	fe1(&r.Z)
	fe0(&r.T2d)
}

func gePrecompCmov(t, u *gePrecomp, b uint) {
	feCmov(&t.yplusx, &u.yplusx, b)
	feCmov(&t.yminusx, &u.yminusx, b)
	feCmov(&t.xy2d, &u.xy2d, b)
}

func geCachedCmov(t *geCached, u *geCached, b uint) {
	feCmov(&t.YplusX, &u.YplusX, b)
	feCmov(&t.YminusX, &u.YminusX, b)
	feCmov(&t.Z, &u.Z, b)
	feCmov(&t.T2d, &u.T2d, b)
}

func equal(b, c uint8) uint {
	// TODO benchmark with different types
	ub := uint8(b)
	uc := uint8(c)
	x := ub ^ uc   // 0: yes; 1..255: no
	y := uint32(x) // 0: yes; 1..255: no
	y -= 1         // 4294967295: yes; 0..254: no
	y >>= 31       // 1: yes; 0: no
	return uint(y)
}

func negative(b int8) uint {
	// TODO benchmark with different types
	x := uint64(b) // 18446744073709551361..18446744073709551615: yes; 0..255
	x >>= 63       // 1: yes; 0: no
	return uint(x)
}

func select_(t *gePrecomp, pos int, b int8) {
	var minust gePrecomp
	bnegative := negative(b)
	babs := uint8(byte(b) - ((byte(-bnegative) & byte(b)) << 1))

	gePrecomp0(t)

	gePrecompCmov(t, &geBase[pos][0], equal(babs, 1))
	gePrecompCmov(t, &geBase[pos][1], equal(babs, 2))
	gePrecompCmov(t, &geBase[pos][2], equal(babs, 3))
	gePrecompCmov(t, &geBase[pos][3], equal(babs, 4))
	gePrecompCmov(t, &geBase[pos][4], equal(babs, 5))
	gePrecompCmov(t, &geBase[pos][5], equal(babs, 6))
	gePrecompCmov(t, &geBase[pos][6], equal(babs, 7))
	gePrecompCmov(t, &geBase[pos][7], equal(babs, 8))
	feCopy(&minust.yplusx, &t.yminusx)
	feCopy(&minust.yminusx, &t.yplusx)
	feNeg(&minust.xy2d, &t.xy2d)
	gePrecompCmov(t, &minust, bnegative)
}

func geMadd(r *geP1P1, p *geP3, q *gePrecomp) {
	// r = p + q

	var t0 fe
	feAdd(&r.X, &p.Y, &p.X)
	feSub(&r.Y, &p.Y, &p.X)
	feMul(&r.Z, &r.X, &q.yplusx)
	feMul(&r.Y, &r.Y, &q.yminusx)
	feMul(&r.T, &q.xy2d, &p.T)
	feAdd(&t0, &p.Z, &p.Z)
	feSub(&r.X, &r.Z, &r.Y)
	feAdd(&r.Y, &r.Z, &r.Y)
	feAdd(&r.Z, &t0, &r.T)
	feSub(&r.T, &t0, &r.T)
}

func geMsub(r *geP1P1, p *geP3, q *gePrecomp) {
	// r = p - q
	var t0 fe
	feAdd(&r.X, &p.Y, &p.X)
	feSub(&r.Y, &p.Y, &p.X)
	feMul(&r.Z, &r.X, &q.yminusx)
	feMul(&r.Y, &r.Y, &q.yplusx)
	feMul(&r.T, &q.xy2d, &p.T)
	feAdd(&t0, &p.Z, &p.Z)
	feSub(&r.X, &r.Z, &r.Y)
	feAdd(&r.Y, &r.Z, &r.Y)
	feSub(&r.Z, &t0, &r.T)
	feAdd(&r.T, &t0, &r.T)
}

func geP1P1ToP2(r *geP2, p *geP1P1) {
	// r = p
	feMul(&r.X, &p.X, &p.T)
	feMul(&r.Y, &p.Y, &p.Z)
	feMul(&r.Z, &p.Z, &p.T)
}

func geP1P1ToP3(r *geP3, p *geP1P1) {
	// r = p
	feMul(&r.X, &p.X, &p.T)
	feMul(&r.Y, &p.Y, &p.Z)
	feMul(&r.Z, &p.Z, &p.T)
	feMul(&r.T, &p.X, &p.Y)
}

func geP3ToCached(r *geCached, p *geP3) {
	// r = p
	feAdd(&r.YplusX, &p.Y, &p.X)
	feSub(&r.YminusX, &p.Y, &p.X)
	feCopy(&r.Z, &p.Z)
	feMul(&r.T2d, &p.T, feD2)
}

func geP2Dbl(r *geP1P1, p *geP2) {
	// r = 2 * p
	var t0 fe
	feSq(&r.X, &p.X)
	feSq(&r.Z, &p.Y)
	feSq2(&r.T, &p.Z)
	feAdd(&r.Y, &p.X, &p.Y)
	feSq(&t0, &r.Y)
	feAdd(&r.Y, &r.Z, &r.X)
	feSub(&r.Z, &r.Z, &r.X)
	feSub(&r.X, &t0, &r.Y)
	feSub(&r.T, &r.T, &r.Z)
}

func geP3ToP2(r *geP2, p *geP3) {
	// r = p
	feCopy(&r.X, &p.X)
	feCopy(&r.Y, &p.Y)
	feCopy(&r.Z, &p.Z)
}

func geP3ToBytes(dst *[32]byte, h *geP3) {
	var recip, x, y fe
	feInvert(&recip, &h.Z)
	feMul(&x, &h.X, &recip)
	feMul(&y, &h.Y, &recip)
	feToBytes(dst, &y)
	dst[31] ^= feIsNegative(&x) << 7
}

func geP3Dbl(r *geP1P1, p *geP3) {
	// r = 2 * p
	var q geP2
	geP3ToP2(&q, p)
	geP2Dbl(r, &q)
}

func geScalarMultBase(h *geP3, a *[32]byte) {
	var (
		e     [64]int8
		carry int8
		r     geP1P1
		s     geP2
		t     gePrecomp
	)

	for i := 0; i < 32; i++ {
		e[2*i+0] = (int8(a[i]) >> 0) & 15
		e[2*i+1] = (int8(a[i]) >> 4) & 15
	}
	// each e[i] is between 0 and 15
	// e[63] is between 0 and 7
	for i := 0; i < 63; i++ {
		e[i] += carry
		carry = e[i] + 8
		carry >>= 4
		e[i] -= carry << 4
	}
	e[63] += carry
	// each e[i] is between -8 and 8

	geP30(h)
	for i := 1; i < 64; i += 2 {
		select_(&t, i/2, e[i])
		geMadd(&r, h, &t)
		geP1P1ToP3(h, &r)
	}

	geP3Dbl(&r, h)
	geP1P1ToP2(&s, &r)
	geP2Dbl(&r, &s)
	geP1P1ToP2(&s, &r)
	geP2Dbl(&r, &s)
	geP1P1ToP2(&s, &r)
	geP2Dbl(&r, &s)
	geP1P1ToP3(h, &r)

	for i := 0; i < 64; i += 2 {
		select_(&t, i/2, e[i])
		geMadd(&r, h, &t)
		geP1P1ToP3(h, &r)
	}
}

func geAdd(r *geP1P1, p *geP3, q *geCached) {
	var t0 fe
	feAdd(&r.X, &p.Y, &p.X)
	feSub(&r.Y, &p.Y, &p.X)
	feMul(&r.Z, &r.X, &q.YplusX)
	feMul(&r.Y, &r.Y, &q.YminusX)
	feMul(&r.T, &q.T2d, &p.T)
	feMul(&r.X, &p.Z, &q.Z)
	feAdd(&t0, &r.X, &r.X)
	feSub(&r.X, &r.Z, &r.Y)
	feAdd(&r.Y, &r.Z, &r.Y)
	feAdd(&r.Z, &t0, &r.T)
	feSub(&r.T, &t0, &r.T)
}

func geSub(r *geP1P1, p *geP3, q *geCached) {
	// r = p - q
	var t0 fe
	feAdd(&r.X, &p.Y, &p.X)
	feSub(&r.Y, &p.Y, &p.X)
	feMul(&r.Z, &r.X, &q.YminusX)
	feMul(&r.Y, &r.Y, &q.YplusX)
	feMul(&r.T, &q.T2d, &p.T)
	feMul(&r.X, &p.Z, &q.Z)
	feAdd(&t0, &r.X, &r.X)
	feSub(&r.X, &r.Z, &r.Y)
	feAdd(&r.Y, &r.Z, &r.Y)
	feSub(&r.Z, &t0, &r.T)
	feAdd(&r.T, &t0, &r.T)
}

func geScalarMult(r *geP2, a *[32]byte, A *geP3) {
	// Assumes that a[31] <= 127
	var (
		e                [64]int8
		i, carry, carry2 int
		Ai               [8]geCached // 1 * A, 2 * A, ..., 8 * A
		//ge_cached Ai[8];
		t geP1P1
		u geP3
	)

	carry = 0 // 0..1
	for i = 0; i < 31; i++ {
		carry += int(a[i])                     // 0..256
		carry2 = (carry + 8) >> 4              // 0..16
		e[2*i] = int8(carry - (carry2 << 4))   // -8..7
		carry = (carry2 + 8) >> 4              // 0..1
		e[2*i+1] = int8(carry2 - (carry << 4)) // -8..7
	}
	carry += int(a[31])                 // 0..128
	carry2 = (carry + 8) >> 4           // 0..8
	e[62] = int8(carry - (carry2 << 4)) // -8..7
	e[63] = int8(carry2)                // 0..8

	geP3ToCached(&Ai[0], A)
	for i = 0; i < 7; i++ {
		geAdd(&t, A, &Ai[i])
		geP1P1ToP3(&u, &t)
		geP3ToCached(&Ai[i+1], &u)
	}

	geP20(r)
	for i = 63; i >= 0; i-- {
		b := e[i]
		bnegative := negative(b)
		babs := uint8(byte(b) - ((byte(-bnegative) & byte(b)) << 1))
		var cur, minuscur geCached

		geP2Dbl(&t, r)
		geP1P1ToP2(r, &t)
		geP2Dbl(&t, r)
		geP1P1ToP2(r, &t)
		geP2Dbl(&t, r)
		geP1P1ToP2(r, &t)
		geP2Dbl(&t, r)
		geP1P1ToP3(&u, &t)
		geCached0(&cur)
		geCachedCmov(&cur, &Ai[0], equal(babs, 1))
		geCachedCmov(&cur, &Ai[1], equal(babs, 2))
		geCachedCmov(&cur, &Ai[2], equal(babs, 3))
		geCachedCmov(&cur, &Ai[3], equal(babs, 4))
		geCachedCmov(&cur, &Ai[4], equal(babs, 5))
		geCachedCmov(&cur, &Ai[5], equal(babs, 6))
		geCachedCmov(&cur, &Ai[6], equal(babs, 7))
		geCachedCmov(&cur, &Ai[7], equal(babs, 8))
		feCopy(&minuscur.YplusX, &cur.YminusX)
		feCopy(&minuscur.YminusX, &cur.YplusX)
		feCopy(&minuscur.Z, &cur.Z)
		feNeg(&minuscur.T2d, &cur.T2d)
		geCachedCmov(&cur, &minuscur, bnegative)
		geAdd(&t, &u, &cur)
		geP1P1ToP2(r, &t)
	}
}

func geToBytes(b *[32]byte, h *geP2) {
	var recip, x, y fe

	feInvert(&recip, &h.Z)
	feMul(&x, &h.X, &recip)
	feMul(&y, &h.Y, &recip)
	feToBytes(b, &y)
	b[31] ^= feIsNegative(&x) << 7
}

func geFromFeFromBytesVarTime(s []byte) *geP2 {
	var (
		u, v, w, x, y, z fe
		sign             uint8
	)

	r := new(geP2)

	h0 := load4(s)
	h1 := load3(s[4:]) << 6
	h2 := load3(s[7:]) << 5
	h3 := load3(s[10:]) << 3
	h4 := load3(s[13:]) << 2
	h5 := load4(s[16:])
	h6 := load3(s[20:]) << 7
	h7 := load3(s[23:]) << 5
	h8 := load3(s[26:]) << 4
	h9 := load3(s[29:]) << 2
	var (
		carry0, carry1, carry2, carry3, carry4 int64
		carry5, carry6, carry7, carry8, carry9 int64
	)

	carry9 = (h9 + (1 << 24)) >> 25
	h0 += carry9 * 19
	h9 -= carry9 << 25
	carry1 = (h1 + (1 << 24)) >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry3 = (h3 + (1 << 24)) >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry5 = (h5 + (1 << 24)) >> 25
	h6 += carry5
	h5 -= carry5 << 25
	carry7 = (h7 + (1 << 24)) >> 25
	h8 += carry7
	h7 -= carry7 << 25

	carry0 = (h0 + (1 << 25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry2 = (h2 + (1 << 25)) >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry4 = (h4 + (1 << 25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry6 = (h6 + (1 << 25)) >> 26
	h7 += carry6
	h6 -= carry6 << 26
	carry8 = (h8 + (1 << 25)) >> 26
	h9 += carry8
	h8 -= carry8 << 26

	u[0] = int32(h0)
	u[1] = int32(h1)
	u[2] = int32(h2)
	u[3] = int32(h3)
	u[4] = int32(h4)
	u[5] = int32(h5)
	u[6] = int32(h6)
	u[7] = int32(h7)
	u[8] = int32(h8)
	u[9] = int32(h9)

	feSq2(&v, &u) // 2 * u^2
	fe1(&w)
	feAdd(&w, &v, &w)        // w = 2 * u^2 + 1
	feSq(&x, &w)             // w^2
	feMul(&y, feMa2, &v)     // -2 * A^2 * u^2
	feAdd(&x, &x, &y)        // x = w^2 - 2 * A^2 * u^2
	feDivPowM1(&r.X, &w, &x) // (&w / x)^(m + 1)
	feSq(&y, &r.X)
	feMul(&x, &y, &x)
	feSub(&y, &w, &x)
	feCopy(&z, feMa)
	if feIsNonZero(&y) {
		feAdd(&y, &w, &x)
		if feIsNonZero(&y) {
			goto negative
		} else {
			feMul(&r.X, &r.X, feFFFB1)
		}
	} else {
		feMul(&r.X, &r.X, feFFFB2)
	}
	feMul(&r.X, &r.X, &u) // u * sqrt(2 * A * (A + 2) * w / x)
	feMul(&z, &z, &v)     // -2 * A * u^2
	sign = 0
	goto setsign
negative:
	feMul(&x, &x, feSqrtm1)
	feSub(&y, &w, &x)
	if feIsNonZero(&y) {
		feMul(&r.X, &r.X, feFFFB3)
	} else {
		feMul(&r.X, &r.X, feFFFB4)
	}
	// &r.X = sqrt(A * (A + 2) * w / x)
	// z = -A
	sign = 1
setsign:
	if feIsNegative(&r.X) != sign {
		//assert(feIsNonZero(&r.X))
		feNeg(&r.X, &r.X)
	}
	feAdd(&r.Z, &z, &w)
	feSub(&r.Y, &z, &w)
	feMul(&r.X, &r.X, &r.Z)
	/*
	   #if !defined(NDEBUG)
	     {
	       fe check_x, check_y, check_iz, check_v;
	       fe_invert(check_iz, &r.Z);
	       feMul(check_x, &r.X, check_iz);
	       feMul(check_y, &r.Y, check_iz);
	       feSq(check_x, check_x);
	       feSq(check_y, check_y);
	       feMul(check_v, check_x, check_y);
	       feMul(check_v, feD, check_v);
	       feAdd(check_v, check_v, check_x);
	       feSub(check_v, check_v, check_y);
	       fe1(check_x);
	       feAdd(check_v, check_v, check_x);
	       assert(!feIsNonZero(check_v));
	     }
	   #endif
	*/
	return r
}

func slide(r *[256]int8, a *[32]byte) {
	var i, b, k uint

	for i = 0; i < 256; i++ {
		r[i] = int8(1 & (int8(a[i>>3]) >> (i & 7)))
	}

	for i = 0; i < 256; i++ {
		if r[i] != 0 {
			for b = 1; b <= 6 && i+b < 256; b++ {
				if r[i+b] != 0 {
					if r[i]+(r[i+b]<<b) <= 15 {
						r[i] += r[i+b] << b
						r[i+b] = 0
					} else if r[i]-(r[i+b]<<b) >= -15 {
						r[i] -= r[i+b] << b
						for k = i + b; k < 256; k++ {
							if r[k] == 0 {
								r[k] = 1
								break
							}
							r[k] = 0
						}
					} else {
						break
					}
				}
			}
		}
	}
}

func geDsmPrecomp(r *geDsmp, s *geP3) {
	var (
		t     geP1P1
		s2, u geP3
	)
	geP3ToCached(&r[0], s)
	geP3Dbl(&t, s)
	geP1P1ToP3(&s2, &t)
	geAdd(&t, &s2, &r[0])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[1], &u)
	geAdd(&t, &s2, &r[1])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[2], &u)
	geAdd(&t, &s2, &r[2])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[3], &u)
	geAdd(&t, &s2, &r[3])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[4], &u)
	geAdd(&t, &s2, &r[4])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[5], &u)
	geAdd(&t, &s2, &r[5])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[6], &u)
	geAdd(&t, &s2, &r[6])
	geP1P1ToP3(&u, &t)
	geP3ToCached(&r[7], &u)
}

func geDoubleScalarMultBaseVarTime(r *geP2, a *[32]byte, A *geP3, b *[32]byte) {
	// r = a * A + b * B
	// where a = a[0]+256*a[1]+...+256^31 a[31].
	// and b = b[0]+256*b[1]+...+256^31 b[31].
	// B is the Ed25519 base point (x,4/5) with x positive.

	var (
		aslide, bslide [256]int8
		Ai             geDsmp // A, 3A, 5A, 7A, 9A, 11A, 13A, 15A
		t              geP1P1
		u              geP3
		i              int
	)

	slide(&aslide, a)
	slide(&bslide, b)
	geDsmPrecomp(&Ai, A)

	geP20(r)

	for i = 255; i >= 0; i-- {
		if aslide[i] != 0 || bslide[i] != 0 {
			break
		}
	}

	for ; i >= 0; i-- {
		geP2Dbl(&t, r)

		if aslide[i] > 0 {
			geP1P1ToP3(&u, &t)
			geAdd(&t, &u, &Ai[aslide[i]/2])
		} else if aslide[i] < 0 {
			geP1P1ToP3(&u, &t)
			geSub(&t, &u, &Ai[(-aslide[i])/2])
		}

		if bslide[i] > 0 {
			geP1P1ToP3(&u, &t)
			geMadd(&t, &u, &geBi[bslide[i]/2])
		} else if bslide[i] < 0 {
			geP1P1ToP3(&u, &t)
			geMsub(&t, &u, &geBi[(-bslide[i])/2])
		}

		geP1P1ToP2(r, &t)
	}

}

func geDoubleScalarMultPrecompVarTime(r *geP2, a *[32]byte, A *geP3, b *[32]byte, Bi *geDsmp) {
	var (
		aslide, bslide [256]int8
		Ai             geDsmp // A, 3A, 5A, 7A, 9A, 11A, 13A, 15A
		t              geP1P1
		u              geP3
		i              int
	)

	slide(&aslide, a)
	slide(&bslide, b)
	geDsmPrecomp(&Ai, A)

	geP20(r)

	for i = 255; i >= 0; i-- {
		if aslide[i] != 0 || bslide[i] != 0 {
			break
		}
	}

	for ; i >= 0; i-- {
		geP2Dbl(&t, r)

		if aslide[i] > 0 {
			geP1P1ToP3(&u, &t)
			geAdd(&t, &u, &Ai[aslide[i]/2])
		} else if aslide[i] < 0 {
			geP1P1ToP3(&u, &t)
			geSub(&t, &u, &Ai[(-aslide[i])/2])
		}

		if bslide[i] > 0 {
			geP1P1ToP3(&u, &t)
			geAdd(&t, &u, &Bi[bslide[i]/2])
		} else if bslide[i] < 0 {
			geP1P1ToP3(&u, &t)
			geSub(&t, &u, &Bi[(-bslide[i])/2])
		}

		geP1P1ToP2(r, &t)
	}
}
