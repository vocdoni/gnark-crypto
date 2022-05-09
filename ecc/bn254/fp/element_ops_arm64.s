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
DATA q<>+0(SB)/8, $4332616871279656263
DATA q<>+8(SB)/8, $10917124144477883021
DATA q<>+16(SB)/8, $13281191951274694749
DATA q<>+24(SB)/8, $3486998266802970665
GLOBL q<>(SB), (RODATA+NOPTR), $32
// qInv0 q'[0]
DATA qInv0<>(SB)/8, $9786893198990664585
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8
#define storeVector(ePtr, e0, e1, e2, e3) \
	STP (e0, e1), 0(ePtr)  \
	STP (e2, e3), 16(ePtr) \

// add(res, x, y *Element)
TEXT ·add(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R4, R5)

	// load operands and add mod 2^r
	LDP  0(R4), (R0, R6)
	LDP  0(R5), (R1, R7)
	ADDS R0, R1, R0
	ADCS R6, R7, R1
	LDP  16(R4), (R2, R6)
	LDP  16(R5), (R3, R7)
	ADCS R2, R3, R2
	ADCS R6, R7, R3

	// load modulus and subtract
	LDP  q<>+0(SB), (R4, R5)
	SUBS R4, R0, R4
	SBCS R5, R1, R5
	LDP  q<>+16(SB), (R6, R7)
	SBCS R6, R2, R6
	SBCS R7, R3, R7

	// reduce if necessary
	CSEL CS, R4, R0, R0
	CSEL CS, R5, R1, R1
	CSEL CS, R6, R2, R2
	CSEL CS, R7, R3, R3

	// store
	MOVD res+0(FP), R4
	storeVector(R4, R0, R1, R2, R3)
	RET

// sub(res, x, y *Element)
TEXT ·sub(SB), NOSPLIT, $0-24
	LDP x+8(FP), (R4, R5)

	// load operands and subtract mod 2^r
	LDP  0(R4), (R0, R6)
	LDP  0(R5), (R1, R7)
	SUBS R1, R0, R0
	SBCS R7, R6, R1
	LDP  16(R4), (R2, R6)
	LDP  16(R5), (R3, R7)
	SBCS R3, R2, R2
	SBCS R7, R6, R3

	// load modulus and select
	MOVD $0, R8
	LDP  q<>+0(SB), (R4, R5)
	CSEL CS, R8, R4, R4
	CSEL CS, R8, R5, R5
	LDP  q<>+16(SB), (R6, R7)
	CSEL CS, R8, R6, R6
	CSEL CS, R8, R7, R7

	// augment (or not)
	ADDS R0, R4, R0
	ADCS R1, R5, R1
	ADCS R2, R6, R2
	ADCS R3, R7, R3

	// store
	MOVD res+0(FP), R4
	storeVector(R4, R0, R1, R2, R3)
	RET

// double(res, x *Element)
TEXT ·double(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R5, R4)

	// load operands and add mod 2^r
	LDP  0(R4), (R0, R1)
	ADDS R0, R0, R0
	ADCS R1, R1, R1
	LDP  16(R4), (R2, R3)
	ADCS R2, R2, R2
	ADCS R3, R3, R3

	// load modulus and subtract
	LDP  q<>+0(SB), (R4, R6)
	SUBS R4, R0, R4
	SBCS R6, R1, R6
	LDP  q<>+16(SB), (R7, R8)
	SBCS R7, R2, R7
	SBCS R8, R3, R8

	// reduce if necessary
	CSEL CS, R4, R0, R0
	CSEL CS, R6, R1, R1
	CSEL CS, R7, R2, R2
	CSEL CS, R8, R3, R3

	// store
	storeVector(R5, R0, R1, R2, R3)
	RET

// neg(res, x *Element)
TEXT ·neg(SB), NOSPLIT, $0-16
	LDP res+0(FP), (R5, R4)

	// load operands and subtract
	MOVD $0, R8
	LDP  0(R4), (R0, R1)
	LDP  q<>+0(SB), (R6, R7)
	ORR  R0, R8, R8              // has x been 0 so far?
	ORR  R1, R8, R8
	SUBS R0, R6, R0
	SBCS R1, R7, R1
	LDP  16(R4), (R2, R3)
	LDP  q<>+16(SB), (R6, R7)
	ORR  R2, R8, R8              // has x been 0 so far?
	ORR  R3, R8, R8
	SBCS R2, R6, R2
	SBCS R3, R7, R3
	TST  $0xffffffffffffffff, R8
	CSEL EQ, R8, R0, R0
	CSEL EQ, R8, R1, R1
	CSEL EQ, R8, R2, R2
	CSEL EQ, R8, R3, R3

	// store
	storeVector(R5, R0, R1, R2, R3)
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

#define loadVector(ePtr, e0, e1, e2, e3) \
	LDP 0(ePtr), (e0, e1)  \
	LDP 16(ePtr), (e2, e3) \

// mul(res, x, y)
TEXT ·mul(SB), NOSPLIT, $0-24
	// Load all of y
	LDP  x+8(FP), (R2, R3)
	loadVector(R3, R13, R14, R15, R16)
	MOVD qInv0<>+0(SB), R4 // Load qInv0

	// Load q
	LDP   q<>+0(SB), (R9, R10)   // R9, R10 = q[0], q[1]
	LDP   q<>+16(SB), (R11, R12) // R11, R12 = q[2], q[3]
	LDP   0(R2), (R0, R1)        // R0, R1 = x[0], x[1]
	MUL   R0, R13, R5
	UMULH R0, R13, R6
	MUL   R4, R8, R5
	madd0(R7, R8, R9, R5)
	madd1(R6, R5, R0, R14, R6)
	madd2(R7, R17, R8, R10, R7, R5)
	madd1(R6, R5, R0, R15, R6)
	madd2(R7, R19, R8, R11, R7, R5)
	madd1(R6, R5, R0, R16, R6)
	madd3(R20, R19, R8, R12, R5, R7, R6)
	madd1(R6, R5, R1, R13, R3)
	MUL   R4, R8, R5
	madd0(R7, R8, R9, R5)
	madd2(R6, R5, R1, R14, R6, R17)
	madd2(R7, R17, R8, R10, R7, R5)
	madd2(R6, R5, R1, R15, R6, R19)
	madd2(R7, R19, R8, R11, R7, R5)
	madd2(R6, R5, R1, R16, R6, R20)
	madd3(R20, R19, R8, R12, R5, R7, R6)
	LDP   16(R2), (R0, R1)       // R0, R1 = x[2], x[3]
	madd1(R6, R5, R0, R13, R3)
	MUL   R4, R8, R5
	madd0(R7, R8, R9, R5)
	madd2(R6, R5, R0, R14, R6, R17)
	madd2(R7, R17, R8, R10, R7, R5)
	madd2(R6, R5, R0, R15, R6, R19)
	madd2(R7, R19, R8, R11, R7, R5)
	madd2(R6, R5, R0, R16, R6, R20)
	madd3(R20, R19, R8, R12, R5, R7, R6)
	madd1(R6, R5, R1, R13, R3)
	MUL   R4, R8, R5
	madd0(R7, R8, R9, R5)
	madd2(R6, R5, R1, R14, R6, R17)
	madd2(R7, R17, R8, R10, R7, R5)
	madd2(R6, R5, R1, R15, R6, R19)
	madd2(R7, R19, R8, R11, R7, R5)
	madd2(R6, R5, R1, R16, R6, R20)
	madd3(R20, R19, R8, R12, R5, R7, R6)

	// Reduce if necessary
	SUBS R9, R3, R13
	SBCS R10, R17, R14
	SBCS R11, R19, R15
	SBCS R12, R20, R16
	CSEL CS, R13, R3, R3
	CSEL CS, R14, R17, R17
	CSEL CS, R15, R19, R19
	CSEL CS, R16, R20, R20
	MOVD res+0(FP), R2     // zPtr
	storeVector(R2, R3, R17, R19, R20)
	RET
