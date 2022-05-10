// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "textflag.h"
#include "funcdata.h"

// modulus q
DATA q<>+0(SB)/8, $5004131540414146891
DATA q<>+8(SB)/8, $10989443069131923958
DATA q<>+16(SB)/8, $3
GLOBL q<>(SB), (RODATA+NOPTR), $24
// qInv0 q'[0]
DATA qInv0<>(SB)/8, $16838227660570584989
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8
#define storeVector(ePtr, e0, e1, e2) \
	STP  (e0, e1), 0(ePtr) \
	MOVD e2, 16(ePtr)      \

// add(res, x, y *Element)
TEXT ·add(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R3, R4)

	// load operands and add mod 2^r
	LDP  0(R3), (R0, R5)
	LDP  0(R4), (R1, R6)
	ADDS R0, R1, R0
	ADCS R5, R6, R1
	MOVD 16(R3), R2      // can't import these in pairs
	MOVD 16(R4), R5
	ADCS R2, R5, R2

	// load modulus and subtract
	LDP  q<>+0(SB), (R3, R4)
	SUBS R3, R0, R3
	SBCS R4, R1, R4
	MOVD q<>+16(SB), R5
	SBCS R5, R2, R5

	// reduce if necessary
	CSEL CS, R3, R0, R0
	CSEL CS, R4, R1, R1
	CSEL CS, R5, R2, R2

	// store
	MOVD res+0(FP), R3
	storeVector(R3, R0, R1, R2)
	RET

// sub(res, x, y *Element)
TEXT ·sub(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R3, R4)

	// load operands and subtract mod 2^r
	LDP  0(R3), (R0, R5)
	LDP  0(R4), (R1, R6)
	SUBS R1, R0, R0
	SBCS R6, R5, R1
	MOVD 16(R3), R2      // can't import these in pairs
	MOVD 16(R4), R5
	SBCS R5, R2, R2

	// load modulus and select
	MOVD $0, R6
	LDP  q<>+0(SB), (R3, R4)
	CSEL CS, R6, R3, R3
	CSEL CS, R6, R4, R4
	MOVD q<>+16(SB), R5
	CSEL CS, R6, R5, R5

	// augment (or not)
	ADDS R0, R3, R0
	ADCS R1, R4, R1
	ADCS R2, R5, R2

	// store
	MOVD res+0(FP), R3
	storeVector(R3, R0, R1, R2)
	RET

// double(res, x *Element)
TEXT ·double(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R4, R3)

	// load operands and add mod 2^r
	LDP  0(R3), (R0, R1)
	ADDS R0, R0, R0
	ADCS R1, R1, R1
	MOVD 16(R3), R2
	ADCS R2, R2, R2

	// load modulus and subtract
	LDP  q<>+0(SB), (R3, R5)
	SUBS R3, R0, R3
	SBCS R5, R1, R5
	MOVD q<>+16(SB), R6
	SBCS R6, R2, R6

	// reduce if necessary
	CSEL CS, R3, R0, R0
	CSEL CS, R5, R1, R1
	CSEL CS, R6, R2, R2

	// store
	storeVector(R4, R0, R1, R2)
	RET

// neg(res, x *Element)
TEXT ·neg(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R4, R3)

	// load operands and subtract
	MOVD $0, R7
	LDP  0(R3), (R0, R1)
	LDP  q<>+0(SB), (R5, R6)
	ORR  R0, R7, R7              // has x been 0 so far?
	ORR  R1, R7, R7
	SUBS R0, R5, R0
	SBCS R1, R6, R1
	MOVD 16(R3), R2              // can't import these in pairs
	MOVD q<>+16(SB), R5
	ORR  R2, R7, R7
	SBCS R2, R5, R2
	TST  $0xffffffffffffffff, R7
	CSEL EQ, R7, R0, R0
	CSEL EQ, R7, R1, R1
	CSEL EQ, R7, R2, R2

	// store
	storeVector(R4, R0, R1, R2)
	RET

// (hi, -) = a*b + c
#define madd0(hi, a, b, c) \
madd1(hi, hi, a, b, c) \

// (hi, lo) = a*b + c
// hi can be the same register as any other operand, including lo
// lo can't be the same register as any of the input
#define madd1(hi, lo, a, b, c) \
	MUL   a, b, lo   \
	ADDS  c, lo, lo  \
	UMULH a, b, hi   \
	ADC   $0, hi, hi \

// madd2 (hi, lo) = a*b + c + d
#define madd2(hi, lo, a, b, c, d) \
madd3(hi, lo, a, b, c, d, $0) \

// madd3 (hi, lo) = a*b + c + d + (e,0)
// hi can be the same register as a, b or c.
#define madd3(hi, lo, a, b, c, d, e) \
	MUL   a, b, lo   \
	ADDS  c, lo, lo  \
	UMULH a, b, hi   \
	ADC   $0, hi, hi \
	ADDS  d, lo, lo  \
	ADC   e, hi, hi  \

#define loadVector(ePtr, e0, e1, e2) \
	LDP  0(ePtr), (e0, e1) \
	MOVD 16(ePtr), e2      \

TEXT ·mul(SB), NOSPLIT, $0-24
	// mul(res, x, y)
#define _qInv0 R4
#define c0 R5
#define c1 R6
#define c2 R7
#define m R8
#define q0 R9
#define q1 R10
#define q2 R11
#define y0 R12
#define y1 R13
#define y2 R14
	// Load all of y
	LDP x+8(FP), (R2, R3)
loadVector(R3, y0, y1, y2)

#define z0 R3
#define z1 R15
#define z2 R16
	MOVD qInv0<>+0(SB), _qInv0 // Load qInv0

	// Load q
	LDP  q<>+0(SB), (q0, q1)
	MOVD q<>+16(SB), q2

	// Round 0
	LDP   0(R2), (R0, R1)
	MUL   R0, y0, c0
	UMULH R0, y0, c1
	MUL   _qInv0, c0, m
	madd0(c2, m, q0, c0)
	madd1(c1, c0, R0, y1, c1)
	madd2(c2, z0, m, q1, c2, c0)
	madd1(c1, c0, R0, y2, c1)
	madd3(z2, z1, m, q2, c0, c2, c1)

	// Round 1
	madd1(c1, c0, R1, y0, z0)
	MUL _qInv0, c0, m
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R1, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R1, y2, c1, z2)
	madd3(z2, z1, m, q2, c0, c2, c1)

	// Round 2
	MOVD 16(R2), R0
	madd1(c1, c0, R0, y0, z0)
	MUL  _qInv0, c0, m
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R0, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R0, y2, c1, z2)
	madd3(z2, z1, m, q2, c0, c2, c1)

	// Reduce if necessary
	SUBS q0, z0, y0
	SBCS q1, z1, y1
	SBCS q2, z2, y2
	CSEL CS, y0, z0, z0
	CSEL CS, y1, z1, z1
	CSEL CS, y2, z2, z2
	MOVD res+0(FP), R2  // zPtr
	storeVector(R2, z0, z1, z2)
	RET

#undef _qInv0
#undef c0
#undef c1
#undef c2
#undef m
#undef q0
#undef q1
#undef q2
#undef y0
#undef y1
#undef y2
#undef z0
#undef z1
#undef z2
