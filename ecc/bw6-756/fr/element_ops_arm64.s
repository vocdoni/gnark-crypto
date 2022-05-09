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
DATA q<>+0(SB)/8, $11045256207009841153
DATA q<>+8(SB)/8, $14886639130118979584
DATA q<>+16(SB)/8, $10956628289047010687
DATA q<>+24(SB)/8, $9513184293603517222
DATA q<>+32(SB)/8, $6038022134869067682
DATA q<>+40(SB)/8, $283357621510263184
GLOBL q<>(SB), (RODATA+NOPTR), $48
// qInv0 q'[0]
DATA qInv0<>(SB)/8, $11045256207009841151
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8
#define storeVector(ePtr, e0, e1, e2, e3, e4, e5) \
	STP (e0, e1), 0(ePtr)  \
	STP (e2, e3), 16(ePtr) \
	STP (e4, e5), 32(ePtr) \

// add(res, x, y *Element)
TEXT ·add(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R6, R7)

	// load operands and add mod 2^r
	LDP  0(R6), (R0, R8)
	LDP  0(R7), (R1, R9)
	ADDS R0, R1, R0
	ADCS R8, R9, R1
	LDP  16(R6), (R2, R8)
	LDP  16(R7), (R3, R9)
	ADCS R2, R3, R2
	ADCS R8, R9, R3
	LDP  32(R6), (R4, R8)
	LDP  32(R7), (R5, R9)
	ADCS R4, R5, R4
	ADCS R8, R9, R5

	// load modulus and subtract
	LDP  q<>+0(SB), (R6, R7)
	SUBS R6, R0, R6
	SBCS R7, R1, R7
	LDP  q<>+16(SB), (R8, R9)
	SBCS R8, R2, R8
	SBCS R9, R3, R9
	LDP  q<>+32(SB), (R10, R11)
	SBCS R10, R4, R10
	SBCS R11, R5, R11

	// reduce if necessary
	CSEL CS, R6, R0, R0
	CSEL CS, R7, R1, R1
	CSEL CS, R8, R2, R2
	CSEL CS, R9, R3, R3
	CSEL CS, R10, R4, R4
	CSEL CS, R11, R5, R5

	// store
	MOVD res+0(FP), R6
	storeVector(R6, R0, R1, R2, R3, R4, R5)
	RET

// sub(res, x, y *Element)
TEXT ·sub(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R6, R7)

	// load operands and subtract mod 2^r
	LDP  0(R6), (R0, R8)
	LDP  0(R7), (R1, R9)
	SUBS R1, R0, R0
	SBCS R9, R8, R1
	LDP  16(R6), (R2, R8)
	LDP  16(R7), (R3, R9)
	SBCS R3, R2, R2
	SBCS R9, R8, R3
	LDP  32(R6), (R4, R8)
	LDP  32(R7), (R5, R9)
	SBCS R5, R4, R4
	SBCS R9, R8, R5

	// load modulus and select
	MOVD $0, R12
	LDP  q<>+0(SB), (R6, R7)
	CSEL CS, R12, R6, R6
	CSEL CS, R12, R7, R7
	LDP  q<>+16(SB), (R8, R9)
	CSEL CS, R12, R8, R8
	CSEL CS, R12, R9, R9
	LDP  q<>+32(SB), (R10, R11)
	CSEL CS, R12, R10, R10
	CSEL CS, R12, R11, R11

	// augment (or not)
	ADDS R0, R6, R0
	ADCS R1, R7, R1
	ADCS R2, R8, R2
	ADCS R3, R9, R3
	ADCS R4, R10, R4
	ADCS R5, R11, R5

	// store
	MOVD res+0(FP), R6
	storeVector(R6, R0, R1, R2, R3, R4, R5)
	RET

// double(res, x *Element)
TEXT ·double(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R7, R6)

	// load operands and add mod 2^r
	LDP  0(R6), (R0, R1)
	ADDS R0, R0, R0
	ADCS R1, R1, R1
	LDP  16(R6), (R2, R3)
	ADCS R2, R2, R2
	ADCS R3, R3, R3
	LDP  32(R6), (R4, R5)
	ADCS R4, R4, R4
	ADCS R5, R5, R5

	// load modulus and subtract
	LDP  q<>+0(SB), (R6, R8)
	SUBS R6, R0, R6
	SBCS R8, R1, R8
	LDP  q<>+16(SB), (R9, R10)
	SBCS R9, R2, R9
	SBCS R10, R3, R10
	LDP  q<>+32(SB), (R11, R12)
	SBCS R11, R4, R11
	SBCS R12, R5, R12

	// reduce if necessary
	CSEL CS, R6, R0, R0
	CSEL CS, R8, R1, R1
	CSEL CS, R9, R2, R2
	CSEL CS, R10, R3, R3
	CSEL CS, R11, R4, R4
	CSEL CS, R12, R5, R5

	// store
	storeVector(R7, R0, R1, R2, R3, R4, R5)
	RET

// neg(res, x *Element)
TEXT ·neg(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R7, R6)

	// load operands and subtract
	MOVD $0, R10
	LDP  0(R6), (R0, R1)
	LDP  q<>+0(SB), (R8, R9)
	ORR  R0, R10, R10             // has x been 0 so far?
	ORR  R1, R10, R10
	SUBS R0, R8, R0
	SBCS R1, R9, R1
	LDP  16(R6), (R2, R3)
	LDP  q<>+16(SB), (R8, R9)
	ORR  R2, R10, R10             // has x been 0 so far?
	ORR  R3, R10, R10
	SBCS R2, R8, R2
	SBCS R3, R9, R3
	LDP  32(R6), (R4, R5)
	LDP  q<>+32(SB), (R8, R9)
	ORR  R4, R10, R10             // has x been 0 so far?
	ORR  R5, R10, R10
	SBCS R4, R8, R4
	SBCS R5, R9, R5
	TST  $0xffffffffffffffff, R10
	CSEL EQ, R10, R0, R0
	CSEL EQ, R10, R1, R1
	CSEL EQ, R10, R2, R2
	CSEL EQ, R10, R3, R3
	CSEL EQ, R10, R4, R4
	CSEL EQ, R10, R5, R5

	// store
	storeVector(R7, R0, R1, R2, R3, R4, R5)
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
madd3(a, b, c, d, $0, hi, lo) \

// madd3 (hi, lo) = a*b + c + d + (e,0)
#define madd3(hi, lo, a, b, c, d, e) \
	MUL   a, b, lo   \
	UMULH a, b, hi   \
	ADDS  c, lo, lo  \
	ADC   $0, hi, hi \
	ADDS  d, lo, lo  \
	ADC   e, hi, hi  \

#define loadVector(ePtr, e0, e1, e2, e3, e4, e5) \
	LDP 0(ePtr), (e0, e1)  \
	LDP 16(ePtr), (e2, e3) \
	LDP 32(ePtr), (e4, e5) \

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
#define q3 R12
#define q4 R13
#define q5 R14
#define y0 R15
#define y1 R16
#define y2 R17
#define y3 R19
#define y4 R20
#define y5 R21
	// Load all of y
	LDP x+8(FP), (R2, R3)
	loadVector(R3, y0, y1, y2, y3, y4, y5)

#define z0 R3
#define z1 R22
#define z2 R23
#define z3 R24
#define z4 R25
#define z5 R26
	MOVD qInv0<>+0(SB), _qInv0 // Load qInv0

	// Load q
	LDP q<>+0(SB), (q0, q1)
	LDP q<>+16(SB), (q2, q3)
	LDP q<>+32(SB), (q4, q5)
	LDP 0(R2), (R0, R1)

	// Round 0
	MUL   R0, y0, c0
	UMULH R0, y0, c1
	MUL   _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd1(c1, c0, R0, y1, c1)
	madd2(c2, z1, m, q1, c2, c0)
	madd1(c1, c0, R0, y2, c1)
	madd2(c2, z2, m, q2, c2, c0)
	madd1(c1, c0, R0, y3, c1)
	madd2(c2, z3, m, q3, c2, c0)
	madd1(c1, c0, R0, y4, c1)
	madd2(c2, z4, m, q4, c2, c0)
	madd1(c1, c0, R0, y5, c1)
	madd3(z5, z4, m, q5, c0, c2, c1)

	// Round 1
	madd1(c1, c0, R1, y0, z0)
	MUL _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R1, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R1, y2, c1, z2)
	madd2(c2, z2, m, q2, c2, c0)
	madd2(c1, c0, R1, y3, c1, z3)
	madd2(c2, z3, m, q3, c2, c0)
	madd2(c1, c0, R1, y4, c1, z4)
	madd2(c2, z4, m, q4, c2, c0)
	madd2(c1, c0, R1, y5, c1, z5)
	madd3(z5, z4, m, q5, c0, c2, c1)
	LDP 16(R2), (R0, R1)

	// Round 2
	madd1(c1, c0, R0, y0, z0)
	MUL _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R0, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R0, y2, c1, z2)
	madd2(c2, z2, m, q2, c2, c0)
	madd2(c1, c0, R0, y3, c1, z3)
	madd2(c2, z3, m, q3, c2, c0)
	madd2(c1, c0, R0, y4, c1, z4)
	madd2(c2, z4, m, q4, c2, c0)
	madd2(c1, c0, R0, y5, c1, z5)
	madd3(z5, z4, m, q5, c0, c2, c1)

	// Round 3
	madd1(c1, c0, R1, y0, z0)
	MUL _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R1, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R1, y2, c1, z2)
	madd2(c2, z2, m, q2, c2, c0)
	madd2(c1, c0, R1, y3, c1, z3)
	madd2(c2, z3, m, q3, c2, c0)
	madd2(c1, c0, R1, y4, c1, z4)
	madd2(c2, z4, m, q4, c2, c0)
	madd2(c1, c0, R1, y5, c1, z5)
	madd3(z5, z4, m, q5, c0, c2, c1)
	LDP 32(R2), (R0, R1)

	// Round 4
	madd1(c1, c0, R0, y0, z0)
	MUL _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R0, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R0, y2, c1, z2)
	madd2(c2, z2, m, q2, c2, c0)
	madd2(c1, c0, R0, y3, c1, z3)
	madd2(c2, z3, m, q3, c2, c0)
	madd2(c1, c0, R0, y4, c1, z4)
	madd2(c2, z4, m, q4, c2, c0)
	madd2(c1, c0, R0, y5, c1, z5)
	madd3(z5, z4, m, q5, c0, c2, c1)

	// Round 5
	madd1(c1, c0, R1, y0, z0)
	MUL _qInv0, m, c0
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R1, y1, c1, z1)
	madd2(c2, z1, m, q1, c2, c0)
	madd2(c1, c0, R1, y2, c1, z2)
	madd2(c2, z2, m, q2, c2, c0)
	madd2(c1, c0, R1, y3, c1, z3)
	madd2(c2, z3, m, q3, c2, c0)
	madd2(c1, c0, R1, y4, c1, z4)
	madd2(c2, z4, m, q4, c2, c0)
	madd2(c1, c0, R1, y5, c1, z5)
	madd3(z5, z4, m, q5, c0, c2, c1)

	// Reduce if necessary
	SUBS q0, z0, y0
	SBCS q1, z1, y1
	SBCS q2, z2, y2
	SBCS q3, z3, y3
	SBCS q4, z4, y4
	SBCS q5, z5, y5
	CSEL CS, y0, z0, z0
	CSEL CS, y1, z1, z1
	CSEL CS, y2, z2, z2
	CSEL CS, y3, z3, z3
	CSEL CS, y4, z4, z4
	CSEL CS, y5, z5, z5
	MOVD res+0(FP), R2  // zPtr
	storeVector(R2, z0, z1, z2, z3, z4, z5)
	RET

#undef q5, c1, q3, z0, z1, z5, q0, y3, m, q1, q4, y4, z3, c0, c2, y0, y1, y2, y5, z2, z4, _qInv0, q2
