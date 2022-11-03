//go:build !amd64 && !arm64
// +build !amd64,!arm64

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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fp

import (
	"math/bits"
)

// Mul z = x * y (mod q)
//
// x and y must be strictly inferior to q
func (z *Element) Mul(x, y *Element) *Element {

	// Implements CIOS multiplication -- section 2.3.2 of Tolga Acar's thesis
	// https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf
	//
	// The algorithm:
	//
	// for i=0 to N-1
	// 		C := 0
	// 		for j=0 to N-1
	// 			(C,t[j]) := t[j] + x[j]*y[i] + C
	// 		(t[N+1],t[N]) := t[N] + C
	//
	// 		C := 0
	// 		m := t[0]*q'[0] mod D
	// 		(C,_) := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 			(C,t[j-1]) := t[j] + m*q[j] + C
	//
	// 		(C,t[N-1]) := t[N] + C
	// 		t[N] := t[N+1] + C
	//
	// → N is the number of machine words needed to store the modulus q
	// → D is the word size. For example, on a 64-bit architecture D is 2	64
	// → x[i], y[i], q[i] is the ith word of the numbers x,y,q
	// → q'[0] is the lowest word of the number -q⁻¹ mod r. This quantity is pre-computed, as it does not depend on the inputs.
	// → t is a temporary array of size N+2
	// → C, S are machine words. A pair (C,S) refers to (hi-bits, lo-bits) of a two-word number
	//
	// As described here https://hackmd.io/@gnark/modular_multiplication we can get rid of one carry chain and simplify:
	//
	// for i=0 to N-1
	// 		(A,t[0]) := t[0] + x[0]*y[i]
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 			(A,t[j])  := t[j] + x[j]*y[i] + A
	// 			(C,t[j-1]) := t[j] + m*q[j] + C
	//
	// 		t[N-1] = C + A
	//
	// This optimization saves 5N + 2 additions in the algorithm, and can be used whenever the highest bit
	// of the modulus is zero (and not all of the remaining bits are set).

	var t [10]uint64
	var c [3]uint64
	{
		// round 0
		v := x[0]
		c[1], c[0] = bits.Mul64(v, y[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd1(v, y[1], c[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd1(v, y[2], c[1])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd1(v, y[3], c[1])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd1(v, y[4], c[1])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd1(v, y[5], c[1])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd1(v, y[6], c[1])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd1(v, y[7], c[1])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd1(v, y[8], c[1])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd1(v, y[9], c[1])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 1
		v := x[1]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 2
		v := x[2]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 3
		v := x[3]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 4
		v := x[4]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 5
		v := x[5]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 6
		v := x[6]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 7
		v := x[7]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 8
		v := x[8]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 9
		v := x[9]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], z[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], z[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		c[2], z[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, y[4], c[1], t[4])
		c[2], z[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, y[5], c[1], t[5])
		c[2], z[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, y[6], c[1], t[6])
		c[2], z[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, y[7], c[1], t[7])
		c[2], z[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, y[8], c[1], t[8])
		c[2], z[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, y[9], c[1], t[9])
		z[9], z[8] = madd3(m, q9, c[0], c[2], c[1])
	}

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], b = bits.Sub64(z[4], q4, b)
		z[5], b = bits.Sub64(z[5], q5, b)
		z[6], b = bits.Sub64(z[6], q6, b)
		z[7], b = bits.Sub64(z[7], q7, b)
		z[8], b = bits.Sub64(z[8], q8, b)
		z[9], _ = bits.Sub64(z[9], q9, b)
	}
	return z
}

// Square z = x * x (mod q)
//
// x must be strictly inferior to q
func (z *Element) Square(x *Element) *Element {
	// see Mul for algorithm documentation

	var t [10]uint64
	var c [3]uint64
	{
		// round 0
		v := x[0]
		c[1], c[0] = bits.Mul64(v, x[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd1(v, x[1], c[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd1(v, x[2], c[1])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd1(v, x[3], c[1])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd1(v, x[4], c[1])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd1(v, x[5], c[1])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd1(v, x[6], c[1])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd1(v, x[7], c[1])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd1(v, x[8], c[1])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd1(v, x[9], c[1])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 1
		v := x[1]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 2
		v := x[2]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 3
		v := x[3]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 4
		v := x[4]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 5
		v := x[5]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 6
		v := x[6]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 7
		v := x[7]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 8
		v := x[8]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], t[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], t[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], t[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], t[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], t[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		t[9], t[8] = madd3(m, q9, c[0], c[2], c[1])
	}
	{
		// round 9
		v := x[9]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], z[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], z[1] = madd2(m, q2, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], z[2] = madd2(m, q3, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		c[2], z[3] = madd2(m, q4, c[2], c[0])
		c[1], c[0] = madd2(v, x[5], c[1], t[5])
		c[2], z[4] = madd2(m, q5, c[2], c[0])
		c[1], c[0] = madd2(v, x[6], c[1], t[6])
		c[2], z[5] = madd2(m, q6, c[2], c[0])
		c[1], c[0] = madd2(v, x[7], c[1], t[7])
		c[2], z[6] = madd2(m, q7, c[2], c[0])
		c[1], c[0] = madd2(v, x[8], c[1], t[8])
		c[2], z[7] = madd2(m, q8, c[2], c[0])
		c[1], c[0] = madd2(v, x[9], c[1], t[9])
		z[9], z[8] = madd3(m, q9, c[0], c[2], c[1])
	}

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], b = bits.Sub64(z[4], q4, b)
		z[5], b = bits.Sub64(z[5], q5, b)
		z[6], b = bits.Sub64(z[6], q6, b)
		z[7], b = bits.Sub64(z[7], q7, b)
		z[8], b = bits.Sub64(z[8], q8, b)
		z[9], _ = bits.Sub64(z[9], q9, b)
	}
	return z
}
