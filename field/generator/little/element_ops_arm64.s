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
DATA q<>+0(SB)/8, $16657056631097876523
DATA q<>+8(SB)/8, $55
GLOBL q<>(SB), (RODATA+NOPTR), $16
// qInv0 q'[0]
DATA qInv0<>(SB)/8, $3679285146481172861
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8
#define storeVector(ePtr, e0, e1) \
	STP (e0, e1), 0(ePtr) \

// add(res, x, y *Element)
TEXT ·add(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R2, R3)

	// load operands and add mod 2^r
	LDP  0(R2), (R0, R4)
	LDP  0(R3), (R1, R5)
	ADDS R0, R1, R0
	ADCS R4, R5, R1

	// load modulus and subtract
	LDP  q<>+0(SB), (R2, R3)
	SUBS R2, R0, R2
	SBCS R3, R1, R3

	// reduce if necessary
	CSEL CS, R2, R0, R0
	CSEL CS, R3, R1, R1

	// store
	MOVD res+0(FP), R2
	storeVector(R2, R0, R1)
	RET

// sub(res, x, y *Element)
TEXT ·sub(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R2, R3)

	// load operands and subtract mod 2^r
	LDP  0(R2), (R0, R4)
	LDP  0(R3), (R1, R5)
	SUBS R1, R0, R0
	SBCS R5, R4, R1

	// load modulus and select
	MOVD $0, R4
	LDP  q<>+0(SB), (R2, R3)
	CSEL CS, R4, R2, R2
	CSEL CS, R4, R3, R3

	// augment (or not)
	ADDS R0, R2, R0
	ADCS R1, R3, R1

	// store
	MOVD res+0(FP), R2
	storeVector(R2, R0, R1)
	RET

// double(res, x *Element)
TEXT ·double(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R3, R2)

	// load operands and add mod 2^r
	LDP  0(R2), (R0, R1)
	ADDS R0, R0, R0
	ADCS R1, R1, R1

	// load modulus and subtract
	LDP  q<>+0(SB), (R2, R4)
	SUBS R2, R0, R2
	SBCS R4, R1, R4

	// reduce if necessary
	CSEL CS, R2, R0, R0
	CSEL CS, R4, R1, R1

	// store
	storeVector(R3, R0, R1)
	RET

// neg(res, x *Element)
TEXT ·neg(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R3, R2)

	// load operands and subtract
	MOVD $0, R6
	LDP  0(R2), (R0, R1)
	LDP  q<>+0(SB), (R4, R5)
	ORR  R0, R6, R6              // has x been 0 so far?
	ORR  R1, R6, R6
	SUBS R0, R4, R0
	SBCS R1, R5, R1
	TST  $0xffffffffffffffff, R6
	CSEL EQ, R6, R0, R0
	CSEL EQ, R6, R1, R1

	// store
	storeVector(R3, R0, R1)
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

#define loadVector(ePtr, e0, e1) \
	LDP 0(ePtr), (e0, e1) \

TEXT ·mul(SB), NOSPLIT, $0-24
	// mul(res, x, y)
#define _qInv0 R4
#define c0 R5
#define c1 R6
#define c2 R7
#define m R8
#define q0 R9
#define q1 R10
#define y0 R11
#define y1 R12
	// Load all of y
	LDP x+8(FP), (R2, R3)
	loadVector(R3, y0, y1)

#define z0 R3
#define z1 R13
	MOVD qInv0<>+0(SB), _qInv0 // Load qInv0

	// Load q
	LDP q<>+0(SB), (q0, q1)
	LDP 0(R2), (R0, R1)

	// Round 0
	MUL   R0, y0, c0
	UMULH R0, y0, c1
	MUL   _qInv0, c0, m
	madd0(c2, m, q0, c0)
	madd1(c1, c0, R0, y1, c1)
	madd3(z1, z0, m, q1, c0, c2, c1)

	// Round 1
	madd1(c1, c0, R1, y0, z0)
	MUL _qInv0, c0, m
	madd0(c2, m, q0, c0)
	madd2(c1, c0, R1, y1, c1, z1)
	madd3(z1, z0, m, q1, c0, c2, c1)

	// Reduce if necessary
	SUBS q0, z0, y0
	SBCS q1, z1, y1
	CSEL CS, y0, z0, z0
	CSEL CS, y1, z1, z1
	MOVD res+0(FP), R2  // zPtr
	storeVector(R2, z0, z1)
	RET

#undef _qInv0
#undef c0
#undef c1
#undef c2
#undef m
#undef q0
#undef q1
#undef y0
#undef y1
#undef z0
#undef z1
