package crypto

type ECScalar [32]byte

func (s *ECScalar) Check() bool { return scCheck((*[32]byte)(s)) }

func load3(b []byte) int64 {
	i := int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	return i
}

func load4(b []byte) int64 {
	i := int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 24
	return i
}

// scReduce performs some sort of unfathomable reduction
// on src[:64] to dst[:32]. If src and dst are the same
// slice it will not affect the output.
func scReduce(dst, src []byte) {
	// Input:
	//   s[0]+256*s[1]+...+256^63*s[63] = s
	//
	// Output:
	//   s[0]+256*s[1]+...+256^31*s[31] = s mod l
	//   where l = 2^252 + 27742317777372353535851937790883648493.

	s0 := int64(0x1fffff & load3(src[0:]))
	s1 := int64(0x1fffff & (load4(src[2:]) >> 5))
	s2 := int64(0x1fffff & (load3(src[5:]) >> 2))
	s3 := int64(0x1fffff & (load4(src[7:]) >> 7))
	s4 := int64(0x1fffff & (load4(src[10:]) >> 4))
	s5 := int64(0x1fffff & (load3(src[13:]) >> 1))
	s6 := int64(0x1fffff & (load4(src[15:]) >> 6))
	s7 := int64(0x1fffff & (load3(src[18:]) >> 3))
	s8 := int64(0x1fffff & load3(src[21:]))
	s9 := int64(0x1fffff & (load4(src[23:]) >> 5))
	s10 := int64(0x1fffff & (load3(src[26:]) >> 2))
	s11 := int64(0x1fffff & (load4(src[28:]) >> 7))
	s12 := int64(0x1fffff & (load4(src[31:]) >> 4))
	s13 := int64(0x1fffff & (load3(src[34:]) >> 1))
	s14 := int64(0x1fffff & (load4(src[36:]) >> 6))
	s15 := int64(0x1fffff & (load3(src[39:]) >> 3))
	s16 := int64(0x1fffff & load3(src[42:]))
	s17 := int64(0x1fffff & (load4(src[44:]) >> 5))
	s18 := int64(0x1fffff & (load3(src[47:]) >> 2))
	s19 := int64(0x1fffff & (load4(src[49:]) >> 7))
	s20 := int64(0x1fffff & (load4(src[52:]) >> 4))
	s21 := int64(0x1fffff & (load3(src[55:]) >> 1))
	s22 := int64(0x1fffff & (load4(src[57:]) >> 6))
	s23 := int64(load4(src[60:]) >> 3)
	var (
		carry0, carry1, carry2, carry3     int64
		carry4, carry5, carry6, carry7     int64
		carry8, carry9, carry10, carry11   int64
		carry12, carry13, carry14, carry15 int64
		carry16                            int64
	)

	s11 += s23 * 666643
	s12 += s23 * 470296
	s13 += s23 * 654183
	s14 -= s23 * 997805
	s15 += s23 * 136657
	s16 -= s23 * 683901

	s10 += s22 * 666643
	s11 += s22 * 470296
	s12 += s22 * 654183
	s13 -= s22 * 997805
	s14 += s22 * 136657
	s15 -= s22 * 683901

	s9 += s21 * 666643
	s10 += s21 * 470296
	s11 += s21 * 654183
	s12 -= s21 * 997805
	s13 += s21 * 136657
	s14 -= s21 * 683901

	s8 += s20 * 666643
	s9 += s20 * 470296
	s10 += s20 * 654183
	s11 -= s20 * 997805
	s12 += s20 * 136657
	s13 -= s20 * 683901

	s7 += s19 * 666643
	s8 += s19 * 470296
	s9 += s19 * 654183
	s10 -= s19 * 997805
	s11 += s19 * 136657
	s12 -= s19 * 683901

	s6 += s18 * 666643
	s7 += s18 * 470296
	s8 += s18 * 654183
	s9 -= s18 * 997805
	s10 += s18 * 136657
	s11 -= s18 * 683901

	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry12 = (s12 + (1 << 20)) >> 21
	s13 += carry12
	s12 -= carry12 << 21
	carry14 = (s14 + (1 << 20)) >> 21
	s15 += carry14
	s14 -= carry14 << 21
	carry16 = (s16 + (1 << 20)) >> 21
	s17 += carry16
	s16 -= carry16 << 21

	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21
	carry13 = (s13 + (1 << 20)) >> 21
	s14 += carry13
	s13 -= carry13 << 21
	carry15 = (s15 + (1 << 20)) >> 21
	s16 += carry15
	s15 -= carry15 << 21

	s5 += s17 * 666643
	s6 += s17 * 470296
	s7 += s17 * 654183
	s8 -= s17 * 997805
	s9 += s17 * 136657
	s10 -= s17 * 683901

	s4 += s16 * 666643
	s5 += s16 * 470296
	s6 += s16 * 654183
	s7 -= s16 * 997805
	s8 += s16 * 136657
	s9 -= s16 * 683901

	s3 += s15 * 666643
	s4 += s15 * 470296
	s5 += s15 * 654183
	s6 -= s15 * 997805
	s7 += s15 * 136657
	s8 -= s15 * 683901

	s2 += s14 * 666643
	s3 += s14 * 470296
	s4 += s14 * 654183
	s5 -= s14 * 997805
	s6 += s14 * 136657
	s7 -= s14 * 683901

	s1 += s13 * 666643
	s2 += s13 * 470296
	s3 += s13 * 654183
	s4 -= s13 * 997805
	s5 += s13 * 136657
	s6 -= s13 * 683901

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	dst[0] = byte(s0 >> 0)
	dst[1] = byte(s0 >> 8)
	dst[2] = byte((s0 >> 16) | (s1 << 5))
	dst[3] = byte(s1 >> 3)
	dst[4] = byte(s1 >> 11)
	dst[5] = byte((s1 >> 19) | (s2 << 2))
	dst[6] = byte(s2 >> 6)
	dst[7] = byte((s2 >> 14) | (s3 << 7))
	dst[8] = byte(s3 >> 1)
	dst[9] = byte(s3 >> 9)
	dst[10] = byte((s3 >> 17) | (s4 << 4))
	dst[11] = byte(s4 >> 4)
	dst[12] = byte(s4 >> 12)
	dst[13] = byte((s4 >> 20) | (s5 << 1))
	dst[14] = byte(s5 >> 7)
	dst[15] = byte((s5 >> 15) | (s6 << 6))
	dst[16] = byte(s6 >> 2)
	dst[17] = byte(s6 >> 10)
	dst[18] = byte((s6 >> 18) | (s7 << 3))
	dst[19] = byte(s7 >> 5)
	dst[20] = byte(s7 >> 13)
	dst[21] = byte(s8 >> 0)
	dst[22] = byte(s8 >> 8)
	dst[23] = byte((s8 >> 16) | (s9 << 5))
	dst[24] = byte(s9 >> 3)
	dst[25] = byte(s9 >> 11)
	dst[26] = byte((s9 >> 19) | (s10 << 2))
	dst[27] = byte(s10 >> 6)
	dst[28] = byte((s10 >> 14) | (s11 << 7))
	dst[29] = byte(s11 >> 1)
	dst[30] = byte(s11 >> 9)
	dst[31] = byte(s11 >> 17)
}

// reduce32 reduces src[:32] to dst[:32].
// If src and dst are the same slice it will not affect the output.
func reduce32(dst, src *[32]byte) {
	s0 := int64(0x1fffff & load3(src[0:]))
	s1 := int64(0x1fffff & (load4(src[2:]) >> 5))
	s2 := int64(0x1fffff & (load3(src[5:]) >> 2))
	s3 := int64(0x1fffff & (load4(src[7:]) >> 7))
	s4 := int64(0x1fffff & (load4(src[10:]) >> 4))
	s5 := int64(0x1fffff & (load3(src[13:]) >> 1))
	s6 := int64(0x1fffff & (load4(src[15:]) >> 6))
	s7 := int64(0x1fffff & (load3(src[18:]) >> 3))
	s8 := int64(0x1fffff & load3(src[21:]))
	s9 := int64(0x1fffff & (load4(src[23:]) >> 5))
	s10 := int64(0x1fffff & (load3(src[26:]) >> 2))
	s11 := (load4(src[28:]) >> 7)
	s12 := int64(0)

	var (
		carry0, carry1, carry2, carry3   int64
		carry4, carry5, carry6, carry7   int64
		carry8, carry9, carry10, carry11 int64
	)

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	dst[0] = byte(s0 >> 0)
	dst[1] = byte(s0 >> 8)
	dst[2] = byte((s0 >> 16) | (s1 << 5))
	dst[3] = byte(s1 >> 3)
	dst[4] = byte(s1 >> 11)
	dst[5] = byte((s1 >> 19) | (s2 << 2))
	dst[6] = byte(s2 >> 6)
	dst[7] = byte((s2 >> 14) | (s3 << 7))
	dst[8] = byte(s3 >> 1)
	dst[9] = byte(s3 >> 9)
	dst[10] = byte((s3 >> 17) | (s4 << 4))
	dst[11] = byte(s4 >> 4)
	dst[12] = byte(s4 >> 12)
	dst[13] = byte((s4 >> 20) | (s5 << 1))
	dst[14] = byte(s5 >> 7)
	dst[15] = byte((s5 >> 15) | (s6 << 6))
	dst[16] = byte(s6 >> 2)
	dst[17] = byte(s6 >> 10)
	dst[18] = byte((s6 >> 18) | (s7 << 3))
	dst[19] = byte(s7 >> 5)
	dst[20] = byte(s7 >> 13)
	dst[21] = byte(s8 >> 0)
	dst[22] = byte(s8 >> 8)
	dst[23] = byte((s8 >> 16) | (s9 << 5))
	dst[24] = byte(s9 >> 3)
	dst[25] = byte(s9 >> 11)
	dst[26] = byte((s9 >> 19) | (s10 << 2))
	dst[27] = byte(s10 >> 6)
	dst[28] = byte((s10 >> 14) | (s11 << 7))
	dst[29] = byte(s11 >> 1)
	dst[30] = byte(s11 >> 9)
	dst[31] = byte(s11 >> 17)
}

func scIsNonZero(s *[32]byte) bool {
	return (((int(s[0]) | int(s[1]) | int(s[2]) | int(s[3]) | int(s[4]) | int(s[5]) | int(s[6]) | int(s[7]) | int(s[8]) |
		int(s[9]) | int(s[10]) | int(s[11]) | int(s[12]) | int(s[13]) | int(s[14]) | int(s[15]) | int(s[16]) | int(s[17]) |
		int(s[18]) | int(s[19]) | int(s[20]) | int(s[21]) | int(s[22]) | int(s[23]) | int(s[24]) | int(s[25]) | int(s[26]) |
		int(s[27]) | int(s[28]) | int(s[29]) | int(s[30]) | int(s[31])) - 1) >> 8) == 0
}

func signum(a int64) int64 {
	// Assumes that a != INT64_MIN
	return (a >> 63) - ((-a) >> 63)
}

func scCheck(s *[32]byte) bool {
	s0 := int64(load4(s[0:]))
	s1 := int64(load4(s[4:]))
	s2 := int64(load4(s[8:]))
	s3 := int64(load4(s[12:]))
	s4 := int64(load4(s[16:]))
	s5 := int64(load4(s[20:]))
	s6 := int64(load4(s[24:]))
	s7 := int64(load4(s[28:]))
	return 0 == ((signum(1559614444-s0) + (signum(1477600026-s1) << 1) + (signum(2734136534-s2) << 2) + (signum(350157278-s3) << 3) + (signum(-s4) << 4) + (signum(-s5) << 5) + (signum(-s6) << 6) + (signum(268435456-s7) << 7)) >> 8)
}

func scAdd(s, a, b *[32]byte) {
	a0 := 0x1fffff & load3(a[0:])
	a1 := 0x1fffff & (load4(a[2:]) >> 5)
	a2 := 0x1fffff & (load3(a[5:]) >> 2)
	a3 := 0x1fffff & (load4(a[7:]) >> 7)
	a4 := 0x1fffff & (load4(a[10:]) >> 4)
	a5 := 0x1fffff & (load3(a[13:]) >> 1)
	a6 := 0x1fffff & (load4(a[15:]) >> 6)
	a7 := 0x1fffff & (load3(a[18:]) >> 3)
	a8 := 0x1fffff & load3(a[21:])
	a9 := 0x1fffff & (load4(a[23:]) >> 5)
	a10 := 0x1fffff & (load3(a[26:]) >> 2)
	a11 := (load4(a[28:]) >> 7)
	b0 := 0x1fffff & load3(b[0:])
	b1 := 0x1fffff & (load4(b[2:]) >> 5)
	b2 := 0x1fffff & (load3(b[5:]) >> 2)
	b3 := 0x1fffff & (load4(b[7:]) >> 7)
	b4 := 0x1fffff & (load4(b[10:]) >> 4)
	b5 := 0x1fffff & (load3(b[13:]) >> 1)
	b6 := 0x1fffff & (load4(b[15:]) >> 6)
	b7 := 0x1fffff & (load3(b[18:]) >> 3)
	b8 := 0x1fffff & load3(b[21:])
	b9 := 0x1fffff & (load4(b[23:]) >> 5)
	b10 := 0x1fffff & (load3(b[26:]) >> 2)
	b11 := (load4(b[28:]) >> 7)
	s0 := a0 + b0
	s1 := a1 + b1
	s2 := a2 + b2
	s3 := a3 + b3
	s4 := a4 + b4
	s5 := a5 + b5
	s6 := a6 + b6
	s7 := a7 + b7
	s8 := a8 + b8
	s9 := a9 + b9
	s10 := a10 + b10
	s11 := a11 + b11
	s12 := int64(0)
	var (
		carry0, carry1, carry2, carry3   int64
		carry4, carry5, carry6, carry7   int64
		carry8, carry9, carry10, carry11 int64
	)

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	s[0] = byte(s0 >> 0)
	s[1] = byte(s0 >> 8)
	s[2] = byte((s0 >> 16) | (s1 << 5))
	s[3] = byte(s1 >> 3)
	s[4] = byte(s1 >> 11)
	s[5] = byte((s1 >> 19) | (s2 << 2))
	s[6] = byte(s2 >> 6)
	s[7] = byte((s2 >> 14) | (s3 << 7))
	s[8] = byte(s3 >> 1)
	s[9] = byte(s3 >> 9)
	s[10] = byte((s3 >> 17) | (s4 << 4))
	s[11] = byte(s4 >> 4)
	s[12] = byte(s4 >> 12)
	s[13] = byte((s4 >> 20) | (s5 << 1))
	s[14] = byte(s5 >> 7)
	s[15] = byte((s5 >> 15) | (s6 << 6))
	s[16] = byte(s6 >> 2)
	s[17] = byte(s6 >> 10)
	s[18] = byte((s6 >> 18) | (s7 << 3))
	s[19] = byte(s7 >> 5)
	s[20] = byte(s7 >> 13)
	s[21] = byte(s8 >> 0)
	s[22] = byte(s8 >> 8)
	s[23] = byte((s8 >> 16) | (s9 << 5))
	s[24] = byte(s9 >> 3)
	s[25] = byte(s9 >> 11)
	s[26] = byte((s9 >> 19) | (s10 << 2))
	s[27] = byte(s10 >> 6)
	s[28] = byte((s10 >> 14) | (s11 << 7))
	s[29] = byte(s11 >> 1)
	s[30] = byte(s11 >> 9)
	s[31] = byte(s11 >> 17)
}

func scSub(s, a, b *[32]byte) {
	a0 := 2097151 & load3(a[:])
	a1 := 2097151 & (load4(a[2:]) >> 5)
	a2 := 2097151 & (load3(a[5:]) >> 2)
	a3 := 2097151 & (load4(a[7:]) >> 7)
	a4 := 2097151 & (load4(a[10:]) >> 4)
	a5 := 2097151 & (load3(a[13:]) >> 1)
	a6 := 2097151 & (load4(a[15:]) >> 6)
	a7 := 2097151 & (load3(a[18:]) >> 3)
	a8 := 2097151 & load3(a[21:])
	a9 := 2097151 & (load4(a[23:]) >> 5)
	a10 := 2097151 & (load3(a[26:]) >> 2)
	a11 := (load4(a[28:]) >> 7)
	b0 := 2097151 & load3(b[:])
	b1 := 2097151 & (load4(b[2:]) >> 5)
	b2 := 2097151 & (load3(b[5:]) >> 2)
	b3 := 2097151 & (load4(b[7:]) >> 7)
	b4 := 2097151 & (load4(b[10:]) >> 4)
	b5 := 2097151 & (load3(b[13:]) >> 1)
	b6 := 2097151 & (load4(b[15:]) >> 6)
	b7 := 2097151 & (load3(b[18:]) >> 3)
	b8 := 2097151 & load3(b[21:])
	b9 := 2097151 & (load4(b[23:]) >> 5)
	b10 := 2097151 & (load3(b[26:]) >> 2)
	b11 := (load4(b[28:]) >> 7)
	s0 := a0 - b0
	s1 := a1 - b1
	s2 := a2 - b2
	s3 := a3 - b3
	s4 := a4 - b4
	s5 := a5 - b5
	s6 := a6 - b6
	s7 := a7 - b7
	s8 := a8 - b8
	s9 := a9 - b9
	s10 := a10 - b10
	s11 := a11 - b11
	s12 := int64(0)
	var (
		carry0,
		carry1,
		carry2,
		carry3,
		carry4,
		carry5,
		carry6,
		carry7,
		carry8,
		carry9,
		carry10,
		carry11 int64
	)

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	s[0] = byte(s0 >> 0)
	s[1] = byte(s0 >> 8)
	s[2] = byte((s0 >> 16) | (s1 << 5))
	s[3] = byte(s1 >> 3)
	s[4] = byte(s1 >> 11)
	s[5] = byte((s1 >> 19) | (s2 << 2))
	s[6] = byte(s2 >> 6)
	s[7] = byte((s2 >> 14) | (s3 << 7))
	s[8] = byte(s3 >> 1)
	s[9] = byte(s3 >> 9)
	s[10] = byte((s3 >> 17) | (s4 << 4))
	s[11] = byte(s4 >> 4)
	s[12] = byte(s4 >> 12)
	s[13] = byte((s4 >> 20) | (s5 << 1))
	s[14] = byte(s5 >> 7)
	s[15] = byte((s5 >> 15) | (s6 << 6))
	s[16] = byte(s6 >> 2)
	s[17] = byte(s6 >> 10)
	s[18] = byte((s6 >> 18) | (s7 << 3))
	s[19] = byte(s7 >> 5)
	s[20] = byte(s7 >> 13)
	s[21] = byte(s8 >> 0)
	s[22] = byte(s8 >> 8)
	s[23] = byte((s8 >> 16) | (s9 << 5))
	s[24] = byte(s9 >> 3)
	s[25] = byte(s9 >> 11)
	s[26] = byte((s9 >> 19) | (s10 << 2))
	s[27] = byte(s10 >> 6)
	s[28] = byte((s10 >> 14) | (s11 << 7))
	s[29] = byte(s11 >> 1)
	s[30] = byte(s11 >> 9)
	s[31] = byte(s11 >> 17)
}

func scMulSub(s, a *[32]byte, b, c []byte) {
	// Input:
	//   a[0]+256*a[1]+...+256^31*a[31] = a
	//   b[0]+256*b[1]+...+256^31*b[31] = b
	//   c[0]+256*c[1]+...+256^31*c[31] = c
	//
	// Output:
	//   s[0]+256*s[1]+...+256^31*s[31] = (c-ab) mod l
	//   where l = 2^252 + 27742317777372353535851937790883648493.

	a0 := 2097151 & load3(a[:])
	a1 := 2097151 & (load4(a[2:]) >> 5)
	a2 := 2097151 & (load3(a[5:]) >> 2)
	a3 := 2097151 & (load4(a[7:]) >> 7)
	a4 := 2097151 & (load4(a[10:]) >> 4)
	a5 := 2097151 & (load3(a[13:]) >> 1)
	a6 := 2097151 & (load4(a[15:]) >> 6)
	a7 := 2097151 & (load3(a[18:]) >> 3)
	a8 := 2097151 & load3(a[21:])
	a9 := 2097151 & (load4(a[23:]) >> 5)
	a10 := 2097151 & (load3(a[26:]) >> 2)
	a11 := (load4(a[28:]) >> 7)
	b0 := 2097151 & load3(b)
	b1 := 2097151 & (load4(b[2:]) >> 5)
	b2 := 2097151 & (load3(b[5:]) >> 2)
	b3 := 2097151 & (load4(b[7:]) >> 7)
	b4 := 2097151 & (load4(b[10:]) >> 4)
	b5 := 2097151 & (load3(b[13:]) >> 1)
	b6 := 2097151 & (load4(b[15:]) >> 6)
	b7 := 2097151 & (load3(b[18:]) >> 3)
	b8 := 2097151 & load3(b[21:])
	b9 := 2097151 & (load4(b[23:]) >> 5)
	b10 := 2097151 & (load3(b[26:]) >> 2)
	b11 := (load4(b[28:]) >> 7)
	c0 := 2097151 & load3(c)
	c1 := 2097151 & (load4(c[2:]) >> 5)
	c2 := 2097151 & (load3(c[5:]) >> 2)
	c3 := 2097151 & (load4(c[7:]) >> 7)
	c4 := 2097151 & (load4(c[10:]) >> 4)
	c5 := 2097151 & (load3(c[13:]) >> 1)
	c6 := 2097151 & (load4(c[15:]) >> 6)
	c7 := 2097151 & (load3(c[18:]) >> 3)
	c8 := 2097151 & load3(c[21:])
	c9 := 2097151 & (load4(c[23:]) >> 5)
	c10 := 2097151 & (load3(c[26:]) >> 2)
	c11 := (load4(c[28:]) >> 7)
	var (
		s0 = int64(0)
		s1,
		s2,
		s3,
		s4,
		s5,
		s6,
		s7,
		s8,
		s9,
		s10,
		s11,
		s12,
		s13,
		s14,
		s15,
		s16,
		s17,
		s18,
		s19,
		s20,
		s21,
		s22,
		s23,
		carry0,
		carry1,
		carry2,
		carry3,
		carry4,
		carry5,
		carry6,
		carry7,
		carry8,
		carry9,
		carry10,
		carry11,
		carry12,
		carry13,
		carry14,
		carry15,
		carry16,
		carry17,
		carry18,
		carry19,
		carry20,
		carry21,
		carry22 int64
	)

	s0 = c0 - a0*b0
	s1 = c1 - (a0*b1 + a1*b0)
	s2 = c2 - (a0*b2 + a1*b1 + a2*b0)
	s3 = c3 - (a0*b3 + a1*b2 + a2*b1 + a3*b0)
	s4 = c4 - (a0*b4 + a1*b3 + a2*b2 + a3*b1 + a4*b0)
	s5 = c5 - (a0*b5 + a1*b4 + a2*b3 + a3*b2 + a4*b1 + a5*b0)
	s6 = c6 - (a0*b6 + a1*b5 + a2*b4 + a3*b3 + a4*b2 + a5*b1 + a6*b0)
	s7 = c7 - (a0*b7 + a1*b6 + a2*b5 + a3*b4 + a4*b3 + a5*b2 + a6*b1 + a7*b0)
	s8 = c8 - (a0*b8 + a1*b7 + a2*b6 + a3*b5 + a4*b4 + a5*b3 + a6*b2 + a7*b1 + a8*b0)
	s9 = c9 - (a0*b9 + a1*b8 + a2*b7 + a3*b6 + a4*b5 + a5*b4 + a6*b3 + a7*b2 + a8*b1 + a9*b0)
	s10 = c10 - (a0*b10 + a1*b9 + a2*b8 + a3*b7 + a4*b6 + a5*b5 + a6*b4 + a7*b3 + a8*b2 + a9*b1 + a10*b0)
	s11 = c11 - (a0*b11 + a1*b10 + a2*b9 + a3*b8 + a4*b7 + a5*b6 + a6*b5 + a7*b4 + a8*b3 + a9*b2 + a10*b1 + a11*b0)
	s12 = -(a1*b11 + a2*b10 + a3*b9 + a4*b8 + a5*b7 + a6*b6 + a7*b5 + a8*b4 + a9*b3 + a10*b2 + a11*b1)
	s13 = -(a2*b11 + a3*b10 + a4*b9 + a5*b8 + a6*b7 + a7*b6 + a8*b5 + a9*b4 + a10*b3 + a11*b2)
	s14 = -(a3*b11 + a4*b10 + a5*b9 + a6*b8 + a7*b7 + a8*b6 + a9*b5 + a10*b4 + a11*b3)
	s15 = -(a4*b11 + a5*b10 + a6*b9 + a7*b8 + a8*b7 + a9*b6 + a10*b5 + a11*b4)
	s16 = -(a5*b11 + a6*b10 + a7*b9 + a8*b8 + a9*b7 + a10*b6 + a11*b5)
	s17 = -(a6*b11 + a7*b10 + a8*b9 + a9*b8 + a10*b7 + a11*b6)
	s18 = -(a7*b11 + a8*b10 + a9*b9 + a10*b8 + a11*b7)
	s19 = -(a8*b11 + a9*b10 + a10*b9 + a11*b8)
	s20 = -(a9*b11 + a10*b10 + a11*b9)
	s21 = -(a10*b11 + a11*b10)
	s22 = -a11 * b11
	s23 = 0

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry12 = (s12 + (1 << 20)) >> 21
	s13 += carry12
	s12 -= carry12 << 21
	carry14 = (s14 + (1 << 20)) >> 21
	s15 += carry14
	s14 -= carry14 << 21
	carry16 = (s16 + (1 << 20)) >> 21
	s17 += carry16
	s16 -= carry16 << 21
	carry18 = (s18 + (1 << 20)) >> 21
	s19 += carry18
	s18 -= carry18 << 21
	carry20 = (s20 + (1 << 20)) >> 21
	s21 += carry20
	s20 -= carry20 << 21
	carry22 = (s22 + (1 << 20)) >> 21
	s23 += carry22
	s22 -= carry22 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21
	carry13 = (s13 + (1 << 20)) >> 21
	s14 += carry13
	s13 -= carry13 << 21
	carry15 = (s15 + (1 << 20)) >> 21
	s16 += carry15
	s15 -= carry15 << 21
	carry17 = (s17 + (1 << 20)) >> 21
	s18 += carry17
	s17 -= carry17 << 21
	carry19 = (s19 + (1 << 20)) >> 21
	s20 += carry19
	s19 -= carry19 << 21
	carry21 = (s21 + (1 << 20)) >> 21
	s22 += carry21
	s21 -= carry21 << 21

	s11 += s23 * 666643
	s12 += s23 * 470296
	s13 += s23 * 654183
	s14 -= s23 * 997805
	s15 += s23 * 136657
	s16 -= s23 * 683901

	s10 += s22 * 666643
	s11 += s22 * 470296
	s12 += s22 * 654183
	s13 -= s22 * 997805
	s14 += s22 * 136657
	s15 -= s22 * 683901

	s9 += s21 * 666643
	s10 += s21 * 470296
	s11 += s21 * 654183
	s12 -= s21 * 997805
	s13 += s21 * 136657
	s14 -= s21 * 683901

	s8 += s20 * 666643
	s9 += s20 * 470296
	s10 += s20 * 654183
	s11 -= s20 * 997805
	s12 += s20 * 136657
	s13 -= s20 * 683901

	s7 += s19 * 666643
	s8 += s19 * 470296
	s9 += s19 * 654183
	s10 -= s19 * 997805
	s11 += s19 * 136657
	s12 -= s19 * 683901

	s6 += s18 * 666643
	s7 += s18 * 470296
	s8 += s18 * 654183
	s9 -= s18 * 997805
	s10 += s18 * 136657
	s11 -= s18 * 683901

	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry12 = (s12 + (1 << 20)) >> 21
	s13 += carry12
	s12 -= carry12 << 21
	carry14 = (s14 + (1 << 20)) >> 21
	s15 += carry14
	s14 -= carry14 << 21
	carry16 = (s16 + (1 << 20)) >> 21
	s17 += carry16
	s16 -= carry16 << 21

	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21
	carry13 = (s13 + (1 << 20)) >> 21
	s14 += carry13
	s13 -= carry13 << 21
	carry15 = (s15 + (1 << 20)) >> 21
	s16 += carry15
	s15 -= carry15 << 21

	s5 += s17 * 666643
	s6 += s17 * 470296
	s7 += s17 * 654183
	s8 -= s17 * 997805
	s9 += s17 * 136657
	s10 -= s17 * 683901

	s4 += s16 * 666643
	s5 += s16 * 470296
	s6 += s16 * 654183
	s7 -= s16 * 997805
	s8 += s16 * 136657
	s9 -= s16 * 683901

	s3 += s15 * 666643
	s4 += s15 * 470296
	s5 += s15 * 654183
	s6 -= s15 * 997805
	s7 += s15 * 136657
	s8 -= s15 * 683901

	s2 += s14 * 666643
	s3 += s14 * 470296
	s4 += s14 * 654183
	s5 -= s14 * 997805
	s6 += s14 * 136657
	s7 -= s14 * 683901

	s1 += s13 * 666643
	s2 += s13 * 470296
	s3 += s13 * 654183
	s4 -= s13 * 997805
	s5 += s13 * 136657
	s6 -= s13 * 683901

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 = (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 = (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 = (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 = (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 = (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 = (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 = (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 = (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 = (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 = (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 = (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	s[0] = byte(s0 >> 0)
	s[1] = byte(s0 >> 8)
	s[2] = byte((s0 >> 16) | (s1 << 5))
	s[3] = byte(s1 >> 3)
	s[4] = byte(s1 >> 11)
	s[5] = byte((s1 >> 19) | (s2 << 2))
	s[6] = byte(s2 >> 6)
	s[7] = byte((s2 >> 14) | (s3 << 7))
	s[8] = byte(s3 >> 1)
	s[9] = byte(s3 >> 9)
	s[10] = byte((s3 >> 17) | (s4 << 4))
	s[11] = byte(s4 >> 4)
	s[12] = byte(s4 >> 12)
	s[13] = byte((s4 >> 20) | (s5 << 1))
	s[14] = byte(s5 >> 7)
	s[15] = byte((s5 >> 15) | (s6 << 6))
	s[16] = byte(s6 >> 2)
	s[17] = byte(s6 >> 10)
	s[18] = byte((s6 >> 18) | (s7 << 3))
	s[19] = byte(s7 >> 5)
	s[20] = byte(s7 >> 13)
	s[21] = byte(s8 >> 0)
	s[22] = byte(s8 >> 8)
	s[23] = byte((s8 >> 16) | (s9 << 5))
	s[24] = byte(s9 >> 3)
	s[25] = byte(s9 >> 11)
	s[26] = byte((s9 >> 19) | (s10 << 2))
	s[27] = byte(s10 >> 6)
	s[28] = byte((s10 >> 14) | (s11 << 7))
	s[29] = byte(s11 >> 1)
	s[30] = byte(s11 >> 9)
	s[31] = byte(s11 >> 17)
}
