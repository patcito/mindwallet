package crypto

type fe [10]int32

var (
	fe_d      = &fe{-10913610, 13857413, -15372611, 6949391, 114729, -8787816, -6275908, -3247719, -18696448, -12055116} /* d */
	fe_sqrtm1 = &fe{-32595792, -7943725, 9377950, 3500415, 12389472, -272473, -25146209, -2005654, 326686, 11406482}     /* sqrt(-1) */
	fe_d2     = &fe{-21827239, -5839606, -30745221, 13898782, 229458, 15978800, -12551817, -6495438, 29715968, 9444199}
)

func fe0(h *fe) {
	// h = 0
	h[0] = 0
	h[1] = 0
	h[2] = 0
	h[3] = 0
	h[4] = 0
	h[5] = 0
	h[6] = 0
	h[7] = 0
	h[8] = 0
	h[9] = 0
}

func fe1(h *fe) {
	// h = 1
	h[0] = 1
	h[1] = 0
	h[2] = 0
	h[3] = 0
	h[4] = 0
	h[5] = 0
	h[6] = 0
	h[7] = 0
	h[8] = 0
	h[9] = 0
}

func feAdd(h, f, g *fe) {
	// h = f + g
	// Can overlap h with f or g.
	//
	// Preconditions:
	//    |f| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.
	//    |g| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.1*2^26,1.1*2^25,1.1*2^26,1.1*2^25,etc.

	h[0] = f[0] + g[0]
	h[1] = f[1] + g[1]
	h[2] = f[2] + g[2]
	h[3] = f[3] + g[3]
	h[4] = f[4] + g[4]
	h[5] = f[5] + g[5]
	h[6] = f[6] + g[6]
	h[7] = f[7] + g[7]
	h[8] = f[8] + g[8]
	h[9] = f[9] + g[9]
}

func feSub(h, f, g *fe) {
	// h = f - g
	// Can overlap h with f or g.
	//
	// Preconditions:
	//    |f| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.
	//    |g| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.1*2^26,1.1*2^25,1.1*2^26,1.1*2^25,etc.

	h[0] = f[0] - g[0]
	h[1] = f[1] - g[1]
	h[2] = f[2] - g[2]
	h[3] = f[3] - g[3]
	h[4] = f[4] - g[4]
	h[5] = f[5] - g[5]
	h[6] = f[6] - g[6]
	h[7] = f[7] - g[7]
	h[8] = f[8] - g[8]
	h[9] = f[9] - g[9]
}

func feNeg(h, f *fe) {
	// h = -f
	//
	// Preconditions:
	//    |f| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.1*2^25,1.1*2^24,1.1*2^25,1.1*2^24,etc.

	h[0] = -f[0]
	h[1] = -f[1]
	h[2] = -f[2]
	h[3] = -f[3]
	h[4] = -f[4]
	h[5] = -f[5]
	h[6] = -f[6]
	h[7] = -f[7]
	h[8] = -f[8]
	h[9] = -f[9]
}

func feSq(h, f *fe) {
	// h = f * f
	// Can overlap h with f.
	//
	// Preconditions:
	//    |f| bounded by 1.65*2^26,1.65*2^25,1.65*2^26,1.65*2^25,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.01*2^25,1.01*2^24,1.01*2^25,1.01*2^24,etc.
	//

	f0 := f[0]
	f1 := f[1]
	f2 := f[2]
	f3 := f[3]
	f4 := f[4]
	f5 := f[5]
	f6 := f[6]
	f7 := f[7]
	f8 := f[8]
	f9 := f[9]
	f0_2 := 2 * f0
	f1_2 := 2 * f1
	f2_2 := 2 * f2
	f3_2 := 2 * f3
	f4_2 := 2 * f4
	f5_2 := 2 * f5
	f6_2 := 2 * f6
	f7_2 := 2 * f7
	f5_38 := 38 * f5 // 1.959375*2^30
	f6_19 := 19 * f6 // 1.959375*2^30
	f7_38 := 38 * f7 // 1.959375*2^30
	f8_19 := 19 * f8 // 1.959375*2^30
	f9_38 := 38 * f9 // 1.959375*2^30
	f0f0 := int64(f0) * int64(f0)
	f0f1_2 := int64(f0_2) * int64(f1)
	f0f2_2 := int64(f0_2) * int64(f2)
	f0f3_2 := int64(f0_2) * int64(f3)
	f0f4_2 := int64(f0_2) * int64(f4)
	f0f5_2 := int64(f0_2) * int64(f5)
	f0f6_2 := int64(f0_2) * int64(f6)
	f0f7_2 := int64(f0_2) * int64(f7)
	f0f8_2 := int64(f0_2) * int64(f8)
	f0f9_2 := int64(f0_2) * int64(f9)
	f1f1_2 := int64(f1_2) * int64(f1)
	f1f2_2 := int64(f1_2) * int64(f2)
	f1f3_4 := int64(f1_2) * int64(f3_2)
	f1f4_2 := int64(f1_2) * int64(f4)
	f1f5_4 := int64(f1_2) * int64(f5_2)
	f1f6_2 := int64(f1_2) * int64(f6)
	f1f7_4 := int64(f1_2) * int64(f7_2)
	f1f8_2 := int64(f1_2) * int64(f8)
	f1f9_76 := int64(f1_2) * int64(f9_38)
	f2f2 := int64(f2) * int64(f2)
	f2f3_2 := int64(f2_2) * int64(f3)
	f2f4_2 := int64(f2_2) * int64(f4)
	f2f5_2 := int64(f2_2) * int64(f5)
	f2f6_2 := int64(f2_2) * int64(f6)
	f2f7_2 := int64(f2_2) * int64(f7)
	f2f8_38 := int64(f2_2) * int64(f8_19)
	f2f9_38 := int64(f2) * int64(f9_38)
	f3f3_2 := int64(f3_2) * int64(f3)
	f3f4_2 := int64(f3_2) * int64(f4)
	f3f5_4 := int64(f3_2) * int64(f5_2)
	f3f6_2 := int64(f3_2) * int64(f6)
	f3f7_76 := int64(f3_2) * int64(f7_38)
	f3f8_38 := int64(f3_2) * int64(f8_19)
	f3f9_76 := int64(f3_2) * int64(f9_38)
	f4f4 := int64(f4) * int64(f4)
	f4f5_2 := int64(f4_2) * int64(f5)
	f4f6_38 := int64(f4_2) * int64(f6_19)
	f4f7_38 := int64(f4) * int64(f7_38)
	f4f8_38 := int64(f4_2) * int64(f8_19)
	f4f9_38 := int64(f4) * int64(f9_38)
	f5f5_38 := int64(f5) * int64(f5_38)
	f5f6_38 := int64(f5_2) * int64(f6_19)
	f5f7_76 := int64(f5_2) * int64(f7_38)
	f5f8_38 := int64(f5_2) * int64(f8_19)
	f5f9_76 := int64(f5_2) * int64(f9_38)
	f6f6_19 := int64(f6) * int64(f6_19)
	f6f7_38 := int64(f6) * int64(f7_38)
	f6f8_38 := int64(f6_2) * int64(f8_19)
	f6f9_38 := int64(f6) * int64(f9_38)
	f7f7_38 := int64(f7) * int64(f7_38)
	f7f8_38 := int64(f7_2) * int64(f8_19)
	f7f9_76 := int64(f7_2) * int64(f9_38)
	f8f8_19 := int64(f8) * int64(f8_19)
	f8f9_38 := int64(f8) * int64(f9_38)
	f9f9_38 := int64(f9) * int64(f9_38)
	h0 := f0f0 + f1f9_76 + f2f8_38 + f3f7_76 + f4f6_38 + f5f5_38
	h1 := f0f1_2 + f2f9_38 + f3f8_38 + f4f7_38 + f5f6_38
	h2 := f0f2_2 + f1f1_2 + f3f9_76 + f4f8_38 + f5f7_76 + f6f6_19
	h3 := f0f3_2 + f1f2_2 + f4f9_38 + f5f8_38 + f6f7_38
	h4 := f0f4_2 + f1f3_4 + f2f2 + f5f9_76 + f6f8_38 + f7f7_38
	h5 := f0f5_2 + f1f4_2 + f2f3_2 + f6f9_38 + f7f8_38
	h6 := f0f6_2 + f1f5_4 + f2f4_2 + f3f3_2 + f7f9_76 + f8f8_19
	h7 := f0f7_2 + f1f6_2 + f2f5_2 + f3f4_2 + f8f9_38
	h8 := f0f8_2 + f1f7_4 + f2f6_2 + f3f5_4 + f4f4 + f9f9_38
	h9 := f0f9_2 + f1f8_2 + f2f7_2 + f3f6_2 + f4f5_2

	var (
		carry0, carry1, carry2, carry3, carry4 int64
		carry5, carry6, carry7, carry8, carry9 int64
	)

	carry0 = (h0 + int64(1<<25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry4 = (h4 + int64(1<<25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26

	carry1 = (h1 + int64(1<<24)) >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry5 = (h5 + int64(1<<24)) >> 25
	h6 += carry5
	h5 -= carry5 << 25

	carry2 = (h2 + int64(1<<25)) >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry6 = (h6 + int64(1<<25)) >> 26
	h7 += carry6
	h6 -= carry6 << 26

	carry3 = (h3 + int64(1<<24)) >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry7 = (h7 + int64(1<<24)) >> 25
	h8 += carry7
	h7 -= carry7 << 25

	carry4 = (h4 + int64(1<<25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry8 = (h8 + int64(1<<25)) >> 26
	h9 += carry8
	h8 -= carry8 << 26

	carry9 = (h9 + int64(1<<24)) >> 25
	h0 += carry9 * 19
	h9 -= carry9 << 25

	carry0 = (h0 + int64(1<<25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26

	h[0] = int32(h0)
	h[1] = int32(h1)
	h[2] = int32(h2)
	h[3] = int32(h3)
	h[4] = int32(h4)
	h[5] = int32(h5)
	h[6] = int32(h6)
	h[7] = int32(h7)
	h[8] = int32(h8)
	h[9] = int32(h9)
}

func feMul(h, f, g *fe) {
	// h = f * g
	// Can overlap h with f or g.
	//
	// Preconditions:
	//    |f| bounded by 1.65*2^26,1.65*2^25,1.65*2^26,1.65*2^25,etc.
	//    |g| bounded by 1.65*2^26,1.65*2^25,1.65*2^26,1.65*2^25,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.01*2^25,1.01*2^24,1.01*2^25,1.01*2^24,etc.

	/*
	   Notes on implementation strategy:

	   Using schoolbook multiplication.
	   Karatsuba would save a little in some cost models.

	   Most multiplications by 2 and 19 are 32-bit precomputations;
	   cheaper than 64-bit postcomputations.

	   There is one remaining multiplication by 19 in the carry chain;
	   one *19 precomputation can be merged into this,
	   but the resulting data flow is considerably less clean.

	   There are 12 carries below.
	   10 of them are 2-way parallelizable and vectorizable.
	   Can get away with 11 carries, but then data flow is much deeper.

	   With tighter constraints on inputs can squeeze carries into int32.
	*/

	f0 := f[0]
	f1 := f[1]
	f2 := f[2]
	f3 := f[3]
	f4 := f[4]
	f5 := f[5]
	f6 := f[6]
	f7 := f[7]
	f8 := f[8]
	f9 := f[9]
	g0 := g[0]
	g1 := g[1]
	g2 := g[2]
	g3 := g[3]
	g4 := g[4]
	g5 := g[5]
	g6 := g[6]
	g7 := g[7]
	g8 := g[8]
	g9 := g[9]
	g1_19 := 19 * g1 /* 1.959375*2^29 */
	g2_19 := 19 * g2 /* 1.959375*2^30; still ok */
	g3_19 := 19 * g3
	g4_19 := 19 * g4
	g5_19 := 19 * g5
	g6_19 := 19 * g6
	g7_19 := 19 * g7
	g8_19 := 19 * g8
	g9_19 := 19 * g9
	f1_2 := 2 * f1
	f3_2 := 2 * f3
	f5_2 := 2 * f5
	f7_2 := 2 * f7
	f9_2 := 2 * f9
	f0g0 := int64(f0) * int64(g0)
	f0g1 := int64(f0) * int64(g1)
	f0g2 := int64(f0) * int64(g2)
	f0g3 := int64(f0) * int64(g3)
	f0g4 := int64(f0) * int64(g4)
	f0g5 := int64(f0) * int64(g5)
	f0g6 := int64(f0) * int64(g6)
	f0g7 := int64(f0) * int64(g7)
	f0g8 := int64(f0) * int64(g8)
	f0g9 := int64(f0) * int64(g9)
	f1g0 := int64(f1) * int64(g0)
	f1g1_2 := int64(f1_2) * int64(g1)
	f1g2 := int64(f1) * int64(g2)
	f1g3_2 := int64(f1_2) * int64(g3)
	f1g4 := int64(f1) * int64(g4)
	f1g5_2 := int64(f1_2) * int64(g5)
	f1g6 := int64(f1) * int64(g6)
	f1g7_2 := int64(f1_2) * int64(g7)
	f1g8 := int64(f1) * int64(g8)
	f1g9_38 := int64(f1_2) * int64(g9_19)
	f2g0 := int64(f2) * int64(g0)
	f2g1 := int64(f2) * int64(g1)
	f2g2 := int64(f2) * int64(g2)
	f2g3 := int64(f2) * int64(g3)
	f2g4 := int64(f2) * int64(g4)
	f2g5 := int64(f2) * int64(g5)
	f2g6 := int64(f2) * int64(g6)
	f2g7 := int64(f2) * int64(g7)
	f2g8_19 := int64(f2) * int64(g8_19)
	f2g9_19 := int64(f2) * int64(g9_19)
	f3g0 := int64(f3) * int64(g0)
	f3g1_2 := int64(f3_2) * int64(g1)
	f3g2 := int64(f3) * int64(g2)
	f3g3_2 := int64(f3_2) * int64(g3)
	f3g4 := int64(f3) * int64(g4)
	f3g5_2 := int64(f3_2) * int64(g5)
	f3g6 := int64(f3) * int64(g6)
	f3g7_38 := int64(f3_2) * int64(g7_19)
	f3g8_19 := int64(f3) * int64(g8_19)
	f3g9_38 := int64(f3_2) * int64(g9_19)
	f4g0 := int64(f4) * int64(g0)
	f4g1 := int64(f4) * int64(g1)
	f4g2 := int64(f4) * int64(g2)
	f4g3 := int64(f4) * int64(g3)
	f4g4 := int64(f4) * int64(g4)
	f4g5 := int64(f4) * int64(g5)
	f4g6_19 := int64(f4) * int64(g6_19)
	f4g7_19 := int64(f4) * int64(g7_19)
	f4g8_19 := int64(f4) * int64(g8_19)
	f4g9_19 := int64(f4) * int64(g9_19)
	f5g0 := int64(f5) * int64(g0)
	f5g1_2 := int64(f5_2) * int64(g1)
	f5g2 := int64(f5) * int64(g2)
	f5g3_2 := int64(f5_2) * int64(g3)
	f5g4 := int64(f5) * int64(g4)
	f5g5_38 := int64(f5_2) * int64(g5_19)
	f5g6_19 := int64(f5) * int64(g6_19)
	f5g7_38 := int64(f5_2) * int64(g7_19)
	f5g8_19 := int64(f5) * int64(g8_19)
	f5g9_38 := int64(f5_2) * int64(g9_19)
	f6g0 := int64(f6) * int64(g0)
	f6g1 := int64(f6) * int64(g1)
	f6g2 := int64(f6) * int64(g2)
	f6g3 := int64(f6) * int64(g3)
	f6g4_19 := int64(f6) * int64(g4_19)
	f6g5_19 := int64(f6) * int64(g5_19)
	f6g6_19 := int64(f6) * int64(g6_19)
	f6g7_19 := int64(f6) * int64(g7_19)
	f6g8_19 := int64(f6) * int64(g8_19)
	f6g9_19 := int64(f6) * int64(g9_19)
	f7g0 := int64(f7) * int64(g0)
	f7g1_2 := int64(f7_2) * int64(g1)
	f7g2 := int64(f7) * int64(g2)
	f7g3_38 := int64(f7_2) * int64(g3_19)
	f7g4_19 := int64(f7) * int64(g4_19)
	f7g5_38 := int64(f7_2) * int64(g5_19)
	f7g6_19 := int64(f7) * int64(g6_19)
	f7g7_38 := int64(f7_2) * int64(g7_19)
	f7g8_19 := int64(f7) * int64(g8_19)
	f7g9_38 := int64(f7_2) * int64(g9_19)
	f8g0 := int64(f8) * int64(g0)
	f8g1 := int64(f8) * int64(g1)
	f8g2_19 := int64(f8) * int64(g2_19)
	f8g3_19 := int64(f8) * int64(g3_19)
	f8g4_19 := int64(f8) * int64(g4_19)
	f8g5_19 := int64(f8) * int64(g5_19)
	f8g6_19 := int64(f8) * int64(g6_19)
	f8g7_19 := int64(f8) * int64(g7_19)
	f8g8_19 := int64(f8) * int64(g8_19)
	f8g9_19 := int64(f8) * int64(g9_19)
	f9g0 := int64(f9) * int64(g0)
	f9g1_38 := int64(f9_2) * int64(g1_19)
	f9g2_19 := int64(f9) * int64(g2_19)
	f9g3_38 := int64(f9_2) * int64(g3_19)
	f9g4_19 := int64(f9) * int64(g4_19)
	f9g5_38 := int64(f9_2) * int64(g5_19)
	f9g6_19 := int64(f9) * int64(g6_19)
	f9g7_38 := int64(f9_2) * int64(g7_19)
	f9g8_19 := int64(f9) * int64(g8_19)
	f9g9_38 := int64(f9_2) * int64(g9_19)
	h0 := f0g0 + f1g9_38 + f2g8_19 + f3g7_38 + f4g6_19 + f5g5_38 + f6g4_19 + f7g3_38 + f8g2_19 + f9g1_38
	h1 := f0g1 + f1g0 + f2g9_19 + f3g8_19 + f4g7_19 + f5g6_19 + f6g5_19 + f7g4_19 + f8g3_19 + f9g2_19
	h2 := f0g2 + f1g1_2 + f2g0 + f3g9_38 + f4g8_19 + f5g7_38 + f6g6_19 + f7g5_38 + f8g4_19 + f9g3_38
	h3 := f0g3 + f1g2 + f2g1 + f3g0 + f4g9_19 + f5g8_19 + f6g7_19 + f7g6_19 + f8g5_19 + f9g4_19
	h4 := f0g4 + f1g3_2 + f2g2 + f3g1_2 + f4g0 + f5g9_38 + f6g8_19 + f7g7_38 + f8g6_19 + f9g5_38
	h5 := f0g5 + f1g4 + f2g3 + f3g2 + f4g1 + f5g0 + f6g9_19 + f7g8_19 + f8g7_19 + f9g6_19
	h6 := f0g6 + f1g5_2 + f2g4 + f3g3_2 + f4g2 + f5g1_2 + f6g0 + f7g9_38 + f8g8_19 + f9g7_38
	h7 := f0g7 + f1g6 + f2g5 + f3g4 + f4g3 + f5g2 + f6g1 + f7g0 + f8g9_19 + f9g8_19
	h8 := f0g8 + f1g7_2 + f2g6 + f3g5_2 + f4g4 + f5g3_2 + f6g2 + f7g1_2 + f8g0 + f9g9_38
	h9 := f0g9 + f1g8 + f2g7 + f3g6 + f4g5 + f5g4 + f6g3 + f7g2 + f8g1 + f9g0
	var (
		carry0, carry1, carry2, carry3, carry4 int64
		carry5, carry6, carry7, carry8, carry9 int64
	)
	/*
	  |h0| <= (1.65*1.65*2^52*(1+19+19+19+19)+1.65*1.65*2^50*(38+38+38+38+38))
	    i.e. |h0| <= 1.4*2^60; narrower ranges for h2, h4, h6, h8
	  |h1| <= (1.65*1.65*2^51*(1+1+19+19+19+19+19+19+19+19))
	    i.e. |h1| <= 1.7*2^59; narrower ranges for h3, h5, h7, h9
	*/

	carry0 = (h0 + int64(1<<25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry4 = (h4 + int64(1<<25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	/* |h0| <= 2^25 */
	/* |h4| <= 2^25 */
	/* |h1| <= 1.71*2^59 */
	/* |h5| <= 1.71*2^59 */

	carry1 = (h1 + int64(1<<24)) >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry5 = (h5 + int64(1<<24)) >> 25
	h6 += carry5
	h5 -= carry5 << 25
	/* |h1| <= 2^24; from now on fits into int32 */
	/* |h5| <= 2^24; from now on fits into int32 */
	/* |h2| <= 1.41*2^60 */
	/* |h6| <= 1.41*2^60 */

	carry2 = (h2 + int64(1<<25)) >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry6 = (h6 + int64(1<<25)) >> 26
	h7 += carry6
	h6 -= carry6 << 26
	/* |h2| <= 2^25; from now on fits into int32 unchanged */
	/* |h6| <= 2^25; from now on fits into int32 unchanged */
	/* |h3| <= 1.71*2^59 */
	/* |h7| <= 1.71*2^59 */

	carry3 = (h3 + int64(1<<24)) >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry7 = (h7 + int64(1<<24)) >> 25
	h8 += carry7
	h7 -= carry7 << 25
	/* |h3| <= 2^24; from now on fits into int32 unchanged */
	/* |h7| <= 2^24; from now on fits into int32 unchanged */
	/* |h4| <= 1.72*2^34 */
	/* |h8| <= 1.41*2^60 */

	carry4 = (h4 + int64(1<<25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry8 = (h8 + int64(1<<25)) >> 26
	h9 += carry8
	h8 -= carry8 << 26
	/* |h4| <= 2^25; from now on fits into int32 unchanged */
	/* |h8| <= 2^25; from now on fits into int32 unchanged */
	/* |h5| <= 1.01*2^24 */
	/* |h9| <= 1.71*2^59 */

	carry9 = (h9 + int64(1<<24)) >> 25
	h0 += carry9 * 19
	h9 -= carry9 << 25
	/* |h9| <= 2^24; from now on fits into int32 unchanged */
	/* |h0| <= 1.1*2^39 */

	carry0 = (h0 + int64(1<<25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	/* |h0| <= 2^25; from now on fits into int32 unchanged */
	/* |h1| <= 1.01*2^24 */

	h[0] = int32(h0)
	h[1] = int32(h1)
	h[2] = int32(h2)
	h[3] = int32(h3)
	h[4] = int32(h4)
	h[5] = int32(h5)
	h[6] = int32(h6)
	h[7] = int32(h7)
	h[8] = int32(h8)
	h[9] = int32(h9)
}

func feDivPowM1(r, u, v *fe) {
	var v3, uv7, t0, t1, t2 fe
	var i int

	feSq(&v3, v)
	feMul(&v3, &v3, v) // v3 = v^3
	feSq(&uv7, &v3)
	feMul(&uv7, &uv7, v)
	feMul(&uv7, &uv7, u) // &uv7 = uv^7

	// fe_pow22523(uv7, &uv7);

	// From fe_pow22523.c

	feSq(&t0, &uv7)
	feSq(&t1, &t0)
	feSq(&t1, &t1)
	feMul(&t1, &uv7, &t1)
	feMul(&t0, &t0, &t1)
	feSq(&t0, &t0)
	feMul(&t0, &t1, &t0)
	feSq(&t1, &t0)
	for i = 0; i < 4; i++ {
		feSq(&t1, &t1)
	}
	feMul(&t0, &t1, &t0)
	feSq(&t1, &t0)
	for i = 0; i < 9; i++ {
		feSq(&t1, &t1)
	}
	feMul(&t1, &t1, &t0)
	feSq(&t2, &t1)
	for i = 0; i < 19; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t1, &t2, &t1)
	for i = 0; i < 10; i++ {
		feSq(&t1, &t1)
	}
	feMul(&t0, &t1, &t0)
	feSq(&t1, &t0)
	for i = 0; i < 49; i++ {
		feSq(&t1, &t1)
	}
	feMul(&t1, &t1, &t0)
	feSq(&t2, &t1)
	for i = 0; i < 99; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t1, &t2, &t1)
	for i = 0; i < 50; i++ {
		feSq(&t1, &t1)
	}
	feMul(&t0, &t1, &t0)
	feSq(&t0, &t0)
	feSq(&t0, &t0)
	feMul(&t0, &t0, &uv7)

	/* End fe_pow22523.c */
	/* &t0 = (uv^7)^((q-5)/8) */
	feMul(&t0, &t0, &v3)
	feMul(r, &t0, u) /* u^(m+1)v^(-(m+1)) */
}

func feIsNegative(f *fe) byte {
	var b [32]byte
	feToBytes(&b, f)
	return byte(b[0] & 1)
}

func feToBytes(dst *[32]byte, h *fe) {
	/*
	   Preconditions:
	     |h| bounded by 1.1*2^26,1.1*2^25,1.1*2^26,1.1*2^25,etc.

	   Write p=2^255-19; q=floor(h/p).
	   Basic claim: q = floor(2^(-255)(h + 19 2^(-25)h9 + 2^(-1))).

	   Proof:
	     Have |h|<=p so |q|<=1 so |19^2 2^(-255) q|<1/4.
	     Also have |h-2^230 h9|<2^231 so |19 2^(-255)(h-2^230 h9)|<1/4.

	     Write y=2^(-1)-19^2 2^(-255)q-19 2^(-255)(h-2^230 h9).
	     Then 0<y<1.

	     Write r=h-pq.
	     Have 0<=r<=p-1=2^255-20.
	     Thus 0<=r+19(2^-255)r<r+19(2^-255)2^255<=2^255-1.

	     Write x=r+19(2^-255)r+y.
	     Then 0<x<2^255 so floor(2^(-255)x) = 0 so floor(q+2^(-255)x) = q.

	     Have q+2^(-255)x = 2^(-255)(h + 19 2^(-25) h9 + 2^(-1))
	     so floor(2^(-255)(h + 19 2^(-25) h9 + 2^(-1))) = q.
	*/

	h0 := h[0]
	h1 := h[1]
	h2 := h[2]
	h3 := h[3]
	h4 := h[4]
	h5 := h[5]
	h6 := h[6]
	h7 := h[7]
	h8 := h[8]
	h9 := h[9]
	var (
		q                                      int32
		carry0, carry1, carry2, carry3, carry4 int32
		carry5, carry6, carry7, carry8, carry9 int32
	)

	q = (19*h9 + ((int32(1)) << 24)) >> 25
	q = (h0 + q) >> 26
	q = (h1 + q) >> 25
	q = (h2 + q) >> 26
	q = (h3 + q) >> 25
	q = (h4 + q) >> 26
	q = (h5 + q) >> 25
	q = (h6 + q) >> 26
	q = (h7 + q) >> 25
	q = (h8 + q) >> 26
	q = (h9 + q) >> 25

	/* Goal: Output h-(2^255-19)q, which is between 0 and 2^255-20. */
	h0 += 19 * q
	/* Goal: Output h-2^255 q, which is between 0 and 2^255-20. */

	carry0 = h0 >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry1 = h1 >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry2 = h2 >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry3 = h3 >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry4 = h4 >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry5 = h5 >> 25
	h6 += carry5
	h5 -= carry5 << 25
	carry6 = h6 >> 26
	h7 += carry6
	h6 -= carry6 << 26
	carry7 = h7 >> 25
	h8 += carry7
	h7 -= carry7 << 25
	carry8 = h8 >> 26
	h9 += carry8
	h8 -= carry8 << 26
	carry9 = h9 >> 25
	h9 -= carry9 << 25
	/* h10 = carry9 */

	/*
	  Goal: Output h0+...+2^255 h10-2^255 q, which is between 0 and 2^255-20.
	  Have h0+...+2^230 h9 between 0 and 2^255-1;
	  evidently 2^255 h10-2^255 q = 0.
	  Goal: Output h0+...+2^230 h9.
	*/

	dst[0] = byte(h0 >> 0)
	dst[1] = byte(h0 >> 8)
	dst[2] = byte(h0 >> 16)
	dst[3] = byte((h0 >> 24) | (h1 << 2))
	dst[4] = byte(h1 >> 6)
	dst[5] = byte(h1 >> 14)
	dst[6] = byte((h1 >> 22) | (h2 << 3))
	dst[7] = byte(h2 >> 5)
	dst[8] = byte(h2 >> 13)
	dst[9] = byte((h2 >> 21) | (h3 << 5))
	dst[10] = byte(h3 >> 3)
	dst[11] = byte(h3 >> 11)
	dst[12] = byte((h3 >> 19) | (h4 << 6))
	dst[13] = byte(h4 >> 2)
	dst[14] = byte(h4 >> 10)
	dst[15] = byte(h4 >> 18)
	dst[16] = byte(h5 >> 0)
	dst[17] = byte(h5 >> 8)
	dst[18] = byte(h5 >> 16)
	dst[19] = byte((h5 >> 24) | (h6 << 1))
	dst[20] = byte(h6 >> 7)
	dst[21] = byte(h6 >> 15)
	dst[22] = byte((h6 >> 23) | (h7 << 3))
	dst[23] = byte(h7 >> 5)
	dst[24] = byte(h7 >> 13)
	dst[25] = byte((h7 >> 21) | (h8 << 4))
	dst[26] = byte(h8 >> 4)
	dst[27] = byte(h8 >> 12)
	dst[28] = byte((h8 >> 20) | (h9 << 6))
	dst[29] = byte(h9 >> 2)
	dst[30] = byte(h9 >> 10)
	dst[31] = byte(h9 >> 18)
}

func feIsNonZero(f *fe) bool {
	var s [32]byte
	feToBytes(&s, f)
	return (((int)(s[0]|s[1]|s[2]|s[3]|s[4]|s[5]|s[6]|s[7]|s[8]|
		s[9]|s[10]|s[11]|s[12]|s[13]|s[14]|s[15]|s[16]|s[17]|
		s[18]|s[19]|s[20]|s[21]|s[22]|s[23]|s[24]|s[25]|s[26]|
		s[27]|s[28]|s[29]|s[30]|s[31]) - 1) >> 8) == 0
}

func feCmov(f, g *fe, b uint) {
	f0 := f[0]
	f1 := f[1]
	f2 := f[2]
	f3 := f[3]
	f4 := f[4]
	f5 := f[5]
	f6 := f[6]
	f7 := f[7]
	f8 := f[8]
	f9 := f[9]
	g0 := g[0]
	g1 := g[1]
	g2 := g[2]
	g3 := g[3]
	g4 := g[4]
	g5 := g[5]
	g6 := g[6]
	g7 := g[7]
	g8 := g[8]
	g9 := g[9]
	x0 := f0 ^ g0
	x1 := f1 ^ g1
	x2 := f2 ^ g2
	x3 := f3 ^ g3
	x4 := f4 ^ g4
	x5 := f5 ^ g5
	x6 := f6 ^ g6
	x7 := f7 ^ g7
	x8 := f8 ^ g8
	x9 := f9 ^ g9
	// assert((((b - 1) & ~b) | ((b - 2) & ~(b - 1))) == (unsigned int) -1);
	b = -b
	x0 &= int32(b)
	x1 &= int32(b)
	x2 &= int32(b)
	x3 &= int32(b)
	x4 &= int32(b)
	x5 &= int32(b)
	x6 &= int32(b)
	x7 &= int32(b)
	x8 &= int32(b)
	x9 &= int32(b)
	f[0] = f0 ^ x0
	f[1] = f1 ^ x1
	f[2] = f2 ^ x2
	f[3] = f3 ^ x3
	f[4] = f4 ^ x4
	f[5] = f5 ^ x5
	f[6] = f6 ^ x6
	f[7] = f7 ^ x7
	f[8] = f8 ^ x8
	f[9] = f9 ^ x9
}

func feCopy(h, f *fe) {
	// TODO benchmark
	h[0] = f[0]
	h[1] = f[1]
	h[2] = f[2]
	h[3] = f[3]
	h[4] = f[4]
	h[5] = f[5]
	h[6] = f[6]
	h[7] = f[7]
	h[8] = f[8]
	h[9] = f[9]
}

func feSq2(h, f *fe) {
	// h = 2 * f * f
	// Can overlap h with f.
	//
	// Preconditions:
	//    |f| bounded by 1.65*2^26,1.65*2^25,1.65*2^26,1.65*2^25,etc.
	//
	// Postconditions:
	//    |h| bounded by 1.01*2^25,1.01*2^24,1.01*2^25,1.01*2^24,etc.

	f0 := f[0]
	f1 := f[1]
	f2 := f[2]
	f3 := f[3]
	f4 := f[4]
	f5 := f[5]
	f6 := f[6]
	f7 := f[7]
	f8 := f[8]
	f9 := f[9]
	f0_2 := 2 * f0
	f1_2 := 2 * f1
	f2_2 := 2 * f2
	f3_2 := 2 * f3
	f4_2 := 2 * f4
	f5_2 := 2 * f5
	f6_2 := 2 * f6
	f7_2 := 2 * f7
	f5_38 := 38 * f5 // 1.959375*2^30
	f6_19 := 19 * f6 // 1.959375*2^30
	f7_38 := 38 * f7 // 1.959375*2^30
	f8_19 := 19 * f8 // 1.959375*2^30
	f9_38 := 38 * f9 // 1.959375*2^30
	f0f0 := int64(f0) * int64(f0)
	f0f1_2 := int64(f0_2) * int64(f1)
	f0f2_2 := int64(f0_2) * int64(f2)
	f0f3_2 := int64(f0_2) * int64(f3)
	f0f4_2 := int64(f0_2) * int64(f4)
	f0f5_2 := int64(f0_2) * int64(f5)
	f0f6_2 := int64(f0_2) * int64(f6)
	f0f7_2 := int64(f0_2) * int64(f7)
	f0f8_2 := int64(f0_2) * int64(f8)
	f0f9_2 := int64(f0_2) * int64(f9)
	f1f1_2 := int64(f1_2) * int64(f1)
	f1f2_2 := int64(f1_2) * int64(f2)
	f1f3_4 := int64(f1_2) * int64(f3_2)
	f1f4_2 := int64(f1_2) * int64(f4)
	f1f5_4 := int64(f1_2) * int64(f5_2)
	f1f6_2 := int64(f1_2) * int64(f6)
	f1f7_4 := int64(f1_2) * int64(f7_2)
	f1f8_2 := int64(f1_2) * int64(f8)
	f1f9_76 := int64(f1_2) * int64(f9_38)
	f2f2 := int64(f2) * int64(f2)
	f2f3_2 := int64(f2_2) * int64(f3)
	f2f4_2 := int64(f2_2) * int64(f4)
	f2f5_2 := int64(f2_2) * int64(f5)
	f2f6_2 := int64(f2_2) * int64(f6)
	f2f7_2 := int64(f2_2) * int64(f7)
	f2f8_38 := int64(f2_2) * int64(f8_19)
	f2f9_38 := int64(f2) * int64(f9_38)
	f3f3_2 := int64(f3_2) * int64(f3)
	f3f4_2 := int64(f3_2) * int64(f4)
	f3f5_4 := int64(f3_2) * int64(f5_2)
	f3f6_2 := int64(f3_2) * int64(f6)
	f3f7_76 := int64(f3_2) * int64(f7_38)
	f3f8_38 := int64(f3_2) * int64(f8_19)
	f3f9_76 := int64(f3_2) * int64(f9_38)
	f4f4 := int64(f4) * int64(f4)
	f4f5_2 := int64(f4_2) * int64(f5)
	f4f6_38 := int64(f4_2) * int64(f6_19)
	f4f7_38 := int64(f4) * int64(f7_38)
	f4f8_38 := int64(f4_2) * int64(f8_19)
	f4f9_38 := int64(f4) * int64(f9_38)
	f5f5_38 := int64(f5) * int64(f5_38)
	f5f6_38 := int64(f5_2) * int64(f6_19)
	f5f7_76 := int64(f5_2) * int64(f7_38)
	f5f8_38 := int64(f5_2) * int64(f8_19)
	f5f9_76 := int64(f5_2) * int64(f9_38)
	f6f6_19 := int64(f6) * int64(f6_19)
	f6f7_38 := int64(f6) * int64(f7_38)
	f6f8_38 := int64(f6_2) * int64(f8_19)
	f6f9_38 := int64(f6) * int64(f9_38)
	f7f7_38 := int64(f7) * int64(f7_38)
	f7f8_38 := int64(f7_2) * int64(f8_19)
	f7f9_76 := int64(f7_2) * int64(f9_38)
	f8f8_19 := int64(f8) * int64(f8_19)
	f8f9_38 := int64(f8) * int64(f9_38)
	f9f9_38 := int64(f9) * int64(f9_38)
	h0 := f0f0 + f1f9_76 + f2f8_38 + f3f7_76 + f4f6_38 + f5f5_38
	h1 := f0f1_2 + f2f9_38 + f3f8_38 + f4f7_38 + f5f6_38
	h2 := f0f2_2 + f1f1_2 + f3f9_76 + f4f8_38 + f5f7_76 + f6f6_19
	h3 := f0f3_2 + f1f2_2 + f4f9_38 + f5f8_38 + f6f7_38
	h4 := f0f4_2 + f1f3_4 + f2f2 + f5f9_76 + f6f8_38 + f7f7_38
	h5 := f0f5_2 + f1f4_2 + f2f3_2 + f6f9_38 + f7f8_38
	h6 := f0f6_2 + f1f5_4 + f2f4_2 + f3f3_2 + f7f9_76 + f8f8_19
	h7 := f0f7_2 + f1f6_2 + f2f5_2 + f3f4_2 + f8f9_38
	h8 := f0f8_2 + f1f7_4 + f2f6_2 + f3f5_4 + f4f4 + f9f9_38
	h9 := f0f9_2 + f1f8_2 + f2f7_2 + f3f6_2 + f4f5_2
	var (
		carry0, carry1, carry2, carry3, carry4 int64
		carry5, carry6, carry7, carry8, carry9 int64
	)

	h0 += h0
	h1 += h1
	h2 += h2
	h3 += h3
	h4 += h4
	h5 += h5
	h6 += h6
	h7 += h7
	h8 += h8
	h9 += h9

	carry0 = (h0 + (1 << 25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26
	carry4 = (h4 + (1 << 25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26

	carry1 = (h1 + (1 << 24)) >> 25
	h2 += carry1
	h1 -= carry1 << 25
	carry5 = (h5 + (1 << 24)) >> 25
	h6 += carry5
	h5 -= carry5 << 25

	carry2 = (h2 + (1 << 25)) >> 26
	h3 += carry2
	h2 -= carry2 << 26
	carry6 = (h6 + (1 << 25)) >> 26
	h7 += carry6
	h6 -= carry6 << 26

	carry3 = (h3 + (1 << 24)) >> 25
	h4 += carry3
	h3 -= carry3 << 25
	carry7 = (h7 + (1 << 24)) >> 25
	h8 += carry7
	h7 -= carry7 << 25

	carry4 = (h4 + (1 << 25)) >> 26
	h5 += carry4
	h4 -= carry4 << 26
	carry8 = (h8 + (1 << 25)) >> 26
	h9 += carry8
	h8 -= carry8 << 26

	carry9 = (h9 + (1 << 24)) >> 25
	h0 += carry9 * 19
	h9 -= carry9 << 25

	carry0 = (h0 + (1 << 25)) >> 26
	h1 += carry0
	h0 -= carry0 << 26

	h[0] = int32(h0)
	h[1] = int32(h1)
	h[2] = int32(h2)
	h[3] = int32(h3)
	h[4] = int32(h4)
	h[5] = int32(h5)
	h[6] = int32(h6)
	h[7] = int32(h7)
	h[8] = int32(h8)
	h[9] = int32(h9)
}

func feInvert(out, z *fe) {
	var t0, t1, t2, t3 fe
	var i int

	feSq(&t0, z)
	feSq(&t1, &t0)
	feSq(&t1, &t1)
	feMul(&t1, z, &t1)
	feMul(&t0, &t0, &t1)
	feSq(&t2, &t0)
	feMul(&t1, &t1, &t2)
	feSq(&t2, &t1)
	for i = 0; i < 4; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t1, &t2, &t1)
	feSq(&t2, &t1)
	for i = 0; i < 9; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t2, &t2, &t1)
	feSq(&t3, &t2)
	for i = 0; i < 19; i++ {
		feSq(&t3, &t3)
	}
	feMul(&t2, &t3, &t2)
	feSq(&t2, &t2)
	for i = 0; i < 9; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t1, &t2, &t1)
	feSq(&t2, &t1)
	for i = 0; i < 49; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t2, &t2, &t1)
	feSq(&t3, &t2)
	for i = 0; i < 99; i++ {
		feSq(&t3, &t3)
	}
	feMul(&t2, &t3, &t2)
	feSq(&t2, &t2)
	for i = 0; i < 49; i++ {
		feSq(&t2, &t2)
	}
	feMul(&t1, &t2, &t1)
	feSq(&t1, &t1)
	for i = 0; i < 4; i++ {
		feSq(&t1, &t1)
	}
	feMul(out, &t1, &t0)

	return
}
