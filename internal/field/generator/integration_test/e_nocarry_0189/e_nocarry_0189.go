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

package integration

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"math/bits"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// e_nocarry_0189 represents a field element stored on 3 words (uint64)
//
// e_nocarry_0189 are assumed to be in Montgomery form in all methods.
//
// Modulus q =
//
//	q[base10] = 723492700717204643878506187592435902978572531888590686927
//	q[base16] = 0x1d819dc22c27779172bb251bfc52ed436e08691dd8d00acf
//
// # Warning
//
// This code has not been audited and is provided as-is. In particular, there is no security guarantees such as constant time implementation or side-channel attack resistance.
type e_nocarry_0189 [3]uint64

const (
	Limbs = 3         // number of 64 bits words needed to represent a e_nocarry_0189
	Bits  = 189       // number of bits needed to represent a e_nocarry_0189
	Bytes = Limbs * 8 // number of bytes needed to represent a e_nocarry_0189
)

// Field modulus q
const (
	q0 uint64 = 7928702720898239183
	q1 uint64 = 8267242343096315203
	q2 uint64 = 2126153956385585041
)

var qe_nocarry_0189 = e_nocarry_0189{
	q0,
	q1,
	q2,
}

var _modulus big.Int // q stored as big.Int

// Modulus returns q as a big.Int
//
//	q[base10] = 723492700717204643878506187592435902978572531888590686927
//	q[base16] = 0x1d819dc22c27779172bb251bfc52ed436e08691dd8d00acf
func Modulus() *big.Int {
	return new(big.Int).Set(&_modulus)
}

// q + r'.r = 1, i.e., qInvNeg = - q⁻¹ mod r
// used for Montgomery reduction
const qInvNeg uint64 = 3804046338715108305

var bigIntPool = sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

func init() {
	_modulus.SetString("1d819dc22c27779172bb251bfc52ed436e08691dd8d00acf", 16)
}

// Newe_nocarry_0189 returns a new e_nocarry_0189 from a uint64 value
//
// it is equivalent to
//
//	var v e_nocarry_0189
//	v.SetUint64(...)
func Newe_nocarry_0189(v uint64) e_nocarry_0189 {
	z := e_nocarry_0189{v}
	z.Mul(&z, &rSquare)
	return z
}

// SetUint64 sets z to v and returns z
func (z *e_nocarry_0189) SetUint64(v uint64) *e_nocarry_0189 {
	//  sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
	*z = e_nocarry_0189{v}
	return z.Mul(z, &rSquare) // z.ToMont()
}

// SetInt64 sets z to v and returns z
func (z *e_nocarry_0189) SetInt64(v int64) *e_nocarry_0189 {

	// absolute value of v
	m := v >> 63
	z.SetUint64(uint64((v ^ m) - m))

	if m != 0 {
		// v is negative
		z.Neg(z)
	}

	return z
}

// Set z = x and returns z
func (z *e_nocarry_0189) Set(x *e_nocarry_0189) *e_nocarry_0189 {
	z[0] = x[0]
	z[1] = x[1]
	z[2] = x[2]
	return z
}

// SetInterface converts provided interface into e_nocarry_0189
// returns an error if provided type is not supported
// supported types:
//
//	e_nocarry_0189
//	*e_nocarry_0189
//	uint64
//	int
//	string (see SetString for valid formats)
//	*big.Int
//	big.Int
//	[]byte
func (z *e_nocarry_0189) SetInterface(i1 interface{}) (*e_nocarry_0189, error) {
	if i1 == nil {
		return nil, errors.New("can't set integration.e_nocarry_0189 with <nil>")
	}

	switch c1 := i1.(type) {
	case e_nocarry_0189:
		return z.Set(&c1), nil
	case *e_nocarry_0189:
		if c1 == nil {
			return nil, errors.New("can't set integration.e_nocarry_0189 with <nil>")
		}
		return z.Set(c1), nil
	case uint8:
		return z.SetUint64(uint64(c1)), nil
	case uint16:
		return z.SetUint64(uint64(c1)), nil
	case uint32:
		return z.SetUint64(uint64(c1)), nil
	case uint:
		return z.SetUint64(uint64(c1)), nil
	case uint64:
		return z.SetUint64(c1), nil
	case int8:
		return z.SetInt64(int64(c1)), nil
	case int16:
		return z.SetInt64(int64(c1)), nil
	case int32:
		return z.SetInt64(int64(c1)), nil
	case int64:
		return z.SetInt64(c1), nil
	case int:
		return z.SetInt64(int64(c1)), nil
	case string:
		return z.SetString(c1)
	case *big.Int:
		if c1 == nil {
			return nil, errors.New("can't set integration.e_nocarry_0189 with <nil>")
		}
		return z.SetBigInt(c1), nil
	case big.Int:
		return z.SetBigInt(&c1), nil
	case []byte:
		return z.SetBytes(c1), nil
	default:
		return nil, errors.New("can't set integration.e_nocarry_0189 from type " + reflect.TypeOf(i1).String())
	}
}

// SetZero z = 0
func (z *e_nocarry_0189) SetZero() *e_nocarry_0189 {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *e_nocarry_0189) SetOne() *e_nocarry_0189 {
	z[0] = 10357354527652293000
	z[1] = 7649037550067684836
	z[2] = 1437512422624871284
	return z
}

// Div z = x*y⁻¹ (mod q)
func (z *e_nocarry_0189) Div(x, y *e_nocarry_0189) *e_nocarry_0189 {
	var yInv e_nocarry_0189
	yInv.Inverse(y)
	z.Mul(x, &yInv)
	return z
}

// Bit returns the i'th bit, with lsb == bit 0.
//
// It is the responsibility of the caller to convert from Montgomery to Regular form if needed.
func (z *e_nocarry_0189) Bit(i uint64) uint64 {
	j := i / 64
	if j >= 3 {
		return 0
	}
	return uint64(z[j] >> (i % 64) & 1)
}

// Equal returns z == x; constant-time
func (z *e_nocarry_0189) Equal(x *e_nocarry_0189) bool {
	return z.NotEqual(x) == 0
}

// NotEqual returns 0 if and only if z == x; constant-time
func (z *e_nocarry_0189) NotEqual(x *e_nocarry_0189) uint64 {
	return (z[2] ^ x[2]) | (z[1] ^ x[1]) | (z[0] ^ x[0])
}

// IsZero returns z == 0
func (z *e_nocarry_0189) IsZero() bool {
	return (z[2] | z[1] | z[0]) == 0
}

// IsOne returns z == 1
func (z *e_nocarry_0189) IsOne() bool {
	return (z[2] ^ 1437512422624871284 | z[1] ^ 7649037550067684836 | z[0] ^ 10357354527652293000) == 0
}

// IsUint64 reports whether z can be represented as an uint64.
func (z *e_nocarry_0189) IsUint64() bool {
	zz := *z
	zz.FromMont()
	return zz.FitsOnOneWord()
}

// Uint64 returns the uint64 representation of x. If x cannot be represented in a uint64, the result is undefined.
func (z *e_nocarry_0189) Uint64() uint64 {
	zz := *z
	zz.FromMont()
	return zz[0]
}

// FitsOnOneWord reports whether z words (except the least significant word) are 0
//
// It is the responsibility of the caller to convert from Montgomery to Regular form if needed.
func (z *e_nocarry_0189) FitsOnOneWord() bool {
	return (z[2] | z[1]) == 0
}

// Cmp compares (lexicographic order) z and x and returns:
//
//	-1 if z <  x
//	 0 if z == x
//	+1 if z >  x
func (z *e_nocarry_0189) Cmp(x *e_nocarry_0189) int {
	_z := *z
	_x := *x
	_z.FromMont()
	_x.FromMont()
	if _z[2] > _x[2] {
		return 1
	} else if _z[2] < _x[2] {
		return -1
	}
	if _z[1] > _x[1] {
		return 1
	} else if _z[1] < _x[1] {
		return -1
	}
	if _z[0] > _x[0] {
		return 1
	} else if _z[0] < _x[0] {
		return -1
	}
	return 0
}

// LexicographicallyLargest returns true if this element is strictly lexicographically
// larger than its negation, false otherwise
func (z *e_nocarry_0189) LexicographicallyLargest() bool {
	// adapted from github.com/zkcrypto/bls12_381
	// we check if the element is larger than (q-1) / 2
	// if z - (((q -1) / 2) + 1) have no underflow, then z > (q-1) / 2

	_z := *z
	_z.FromMont()

	var b uint64
	_, b = bits.Sub64(_z[0], 13187723397303895400, 0)
	_, b = bits.Sub64(_z[1], 13356993208402933409, b)
	_, b = bits.Sub64(_z[2], 1063076978192792520, b)

	return b == 0
}

// SetRandom sets z to a uniform random value in [0, q).
//
// This might error only if reading from crypto/rand.Reader errors,
// in which case, value of z is undefined.
func (z *e_nocarry_0189) SetRandom() (*e_nocarry_0189, error) {
	// this code is generated for all modulus
	// and derived from go/src/crypto/rand/util.go

	// l is number of limbs * 8; the number of bytes needed to reconstruct 3 uint64
	const l = 24

	// bitLen is the maximum bit length needed to encode a value < q.
	const bitLen = 189

	// k is the maximum byte length needed to encode a value < q.
	const k = (bitLen + 7) / 8

	// b is the number of bits in the most significant byte of q-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	var bytes [l]byte

	for {
		// note that bytes[k:l] is always 0
		if _, err := io.ReadFull(rand.Reader, bytes[:k]); err != nil {
			return nil, err
		}

		// Clear unused bits in in the most signicant byte to increase probability
		// that the candidate is < q.
		bytes[k-1] &= uint8(int(1<<b) - 1)
		z[0] = binary.LittleEndian.Uint64(bytes[0:8])
		z[1] = binary.LittleEndian.Uint64(bytes[8:16])
		z[2] = binary.LittleEndian.Uint64(bytes[16:24])

		if !z.smallerThanModulus() {
			continue // ignore the candidate and re-sample
		}

		return z, nil
	}
}

// smallerThanModulus returns true if z < q
// This is not constant time
func (z *e_nocarry_0189) smallerThanModulus() bool {
	return (z[2] < q2 || (z[2] == q2 && (z[1] < q1 || (z[1] == q1 && (z[0] < q0)))))
}

// One returns 1
func One() e_nocarry_0189 {
	var one e_nocarry_0189
	one.SetOne()
	return one
}

// Halve sets z to z / 2 (mod q)
func (z *e_nocarry_0189) Halve() {
	var carry uint64

	if z[0]&1 == 1 {
		// z = z + q
		z[0], carry = bits.Add64(z[0], q0, 0)
		z[1], carry = bits.Add64(z[1], q1, carry)
		z[2], _ = bits.Add64(z[2], q2, carry)

	}
	// z = z >> 1
	z[0] = z[0]>>1 | z[1]<<63
	z[1] = z[1]>>1 | z[2]<<63
	z[2] >>= 1

}

// Mul z = x * y (mod q)
//
// x and y must be strictly inferior to q
func (z *e_nocarry_0189) Mul(x, y *e_nocarry_0189) *e_nocarry_0189 {
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
	mul(z, x, y)
	return z
}

// Square z = x * x (mod q)
//
// x must be strictly inferior to q
func (z *e_nocarry_0189) Square(x *e_nocarry_0189) *e_nocarry_0189 {
	// see Mul for algorithm documentation
	mul(z, x, x)
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *e_nocarry_0189) FromMont() *e_nocarry_0189 {
	fromMont(z)
	return z
}

// Add z = x + y (mod q)
func (z *e_nocarry_0189) Add(x, y *e_nocarry_0189) *e_nocarry_0189 {

	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], _ = bits.Add64(x[2], y[2], carry)

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}
	return z
}

// Double z = x + x (mod q), aka Lsh 1
func (z *e_nocarry_0189) Double(x *e_nocarry_0189) *e_nocarry_0189 {

	var carry uint64
	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], _ = bits.Add64(x[2], x[2], carry)

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}
	return z
}

// Sub z = x - y (mod q)
func (z *e_nocarry_0189) Sub(x, y *e_nocarry_0189) *e_nocarry_0189 {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	z[1], b = bits.Sub64(x[1], y[1], b)
	z[2], b = bits.Sub64(x[2], y[2], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], q0, 0)
		z[1], c = bits.Add64(z[1], q1, c)
		z[2], _ = bits.Add64(z[2], q2, c)
	}
	return z
}

// Neg z = q - x
func (z *e_nocarry_0189) Neg(x *e_nocarry_0189) *e_nocarry_0189 {
	if x.IsZero() {
		z.SetZero()
		return z
	}
	var borrow uint64
	z[0], borrow = bits.Sub64(q0, x[0], 0)
	z[1], borrow = bits.Sub64(q1, x[1], borrow)
	z[2], _ = bits.Sub64(q2, x[2], borrow)
	return z
}

// Select is a constant-time conditional move.
// If c=0, z = x0. Else z = x1
func (z *e_nocarry_0189) Select(c int, x0 *e_nocarry_0189, x1 *e_nocarry_0189) *e_nocarry_0189 {
	cC := uint64((int64(c) | -int64(c)) >> 63) // "canonicized" into: 0 if c=0, -1 otherwise
	z[0] = x0[0] ^ cC&(x0[0]^x1[0])
	z[1] = x0[1] ^ cC&(x0[1]^x1[1])
	z[2] = x0[2] ^ cC&(x0[2]^x1[2])
	return z
}

func _mulGeneric(z, x, y *e_nocarry_0189) {
	// see Mul for algorithm documentation

	var t [3]uint64
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
		t[2], t[1] = madd3(m, q2, c[0], c[2], c[1])
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
		t[2], t[1] = madd3(m, q2, c[0], c[2], c[1])
	}
	{
		// round 2
		v := x[2]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], z[0] = madd2(m, q1, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		z[2], z[1] = madd3(m, q2, c[0], c[2], c[1])
	}

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}

}

func _fromMontGeneric(z *e_nocarry_0189) {
	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	// see Mul for algorithm documentation
	{
		// m = z[0]n'[0] mod W
		m := z[0] * qInvNeg
		C := madd0(m, q0, z[0])
		C, z[0] = madd2(m, q1, z[1], C)
		C, z[1] = madd2(m, q2, z[2], C)
		z[2] = C
	}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * qInvNeg
		C := madd0(m, q0, z[0])
		C, z[0] = madd2(m, q1, z[1], C)
		C, z[1] = madd2(m, q2, z[2], C)
		z[2] = C
	}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * qInvNeg
		C := madd0(m, q0, z[0])
		C, z[0] = madd2(m, q1, z[1], C)
		C, z[1] = madd2(m, q2, z[2], C)
		z[2] = C
	}

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}
}

func _reduceGeneric(z *e_nocarry_0189) {

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}
}

// BatchInvert returns a new slice with every element inverted.
// Uses Montgomery batch inversion trick
func BatchInvert(a []e_nocarry_0189) []e_nocarry_0189 {
	res := make([]e_nocarry_0189, len(a))
	if len(a) == 0 {
		return res
	}

	zeroes := make([]bool, len(a))
	accumulator := One()

	for i := 0; i < len(a); i++ {
		if a[i].IsZero() {
			zeroes[i] = true
			continue
		}
		res[i] = accumulator
		accumulator.Mul(&accumulator, &a[i])
	}

	accumulator.Inverse(&accumulator)

	for i := len(a) - 1; i >= 0; i-- {
		if zeroes[i] {
			continue
		}
		res[i].Mul(&res[i], &accumulator)
		accumulator.Mul(&accumulator, &a[i])
	}

	return res
}

func _butterflyGeneric(a, b *e_nocarry_0189) {
	t := *a
	a.Add(a, b)
	b.Sub(&t, b)
}

// BitLen returns the minimum number of bits needed to represent z
// returns 0 if z == 0
func (z *e_nocarry_0189) BitLen() int {
	if z[2] != 0 {
		return 128 + bits.Len64(z[2])
	}
	if z[1] != 0 {
		return 64 + bits.Len64(z[1])
	}
	return bits.Len64(z[0])
}

// Exp z = xᵏ (mod q)
func (z *e_nocarry_0189) Exp(x e_nocarry_0189, k *big.Int) *e_nocarry_0189 {
	if k.IsUint64() && k.Uint64() == 0 {
		return z.SetOne()
	}

	e := k
	if k.Sign() == -1 {
		// negative k, we invert
		// if k < 0: xᵏ (mod q) == (x⁻¹)ᵏ (mod q)
		x.Inverse(&x)

		// we negate k in a temp big.Int since
		// Int.Bit(_) of k and -k is different
		e = bigIntPool.Get().(*big.Int)
		defer bigIntPool.Put(e)
		e.Neg(k)
	}

	z.Set(&x)

	for i := e.BitLen() - 2; i >= 0; i-- {
		z.Square(z)
		if e.Bit(i) == 1 {
			z.Mul(z, &x)
		}
	}

	return z
}

// rSquare where r is the Montgommery constant
// see section 2.3.2 of Tolga Acar's thesis
// https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf
var rSquare = e_nocarry_0189{
	13016996726018978220,
	15774878995660659205,
	136514380579669223,
}

// ToMont converts z to Montgomery form
// sets and returns z = z * r²
func (z *e_nocarry_0189) ToMont() *e_nocarry_0189 {
	return z.Mul(z, &rSquare)
}

// ToRegular returns z in regular form (doesn't mutate z)
func (z e_nocarry_0189) ToRegular() e_nocarry_0189 {
	return *z.FromMont()
}

// String returns the decimal representation of z as generated by
// z.Text(10).
func (z *e_nocarry_0189) String() string {
	return z.Text(10)
}

// Text returns the string representation of z in the given base.
// Base must be between 2 and 36, inclusive. The result uses the
// lower-case letters 'a' to 'z' for digit values 10 to 35.
// No prefix (such as "0x") is added to the string. If z is a nil
// pointer it returns "<nil>".
// If base == 10 and -z fits in a uint16 prefix "-" is added to the string.
func (z *e_nocarry_0189) Text(base int) string {
	if base < 2 || base > 36 {
		panic("invalid base")
	}
	if z == nil {
		return "<nil>"
	}

	const maxUint16 = 65535
	if base == 10 {
		var zzNeg e_nocarry_0189
		zzNeg.Neg(z)
		zzNeg.FromMont()
		if zzNeg.FitsOnOneWord() && zzNeg[0] <= maxUint16 && zzNeg[0] != 0 {
			return "-" + strconv.FormatUint(zzNeg[0], base)
		}
	}
	zz := *z
	zz.FromMont()
	if zz.FitsOnOneWord() {
		return strconv.FormatUint(zz[0], base)
	}
	vv := bigIntPool.Get().(*big.Int)
	r := zz.ToBigInt(vv).Text(base)
	bigIntPool.Put(vv)
	return r
}

// ToBigInt returns z as a big.Int in Montgomery form
func (z *e_nocarry_0189) ToBigInt(res *big.Int) *big.Int {
	var b [Limbs * 8]byte
	binary.BigEndian.PutUint64(b[16:24], z[0])
	binary.BigEndian.PutUint64(b[8:16], z[1])
	binary.BigEndian.PutUint64(b[0:8], z[2])

	return res.SetBytes(b[:])
}

// ToBigIntRegular returns z as a big.Int in regular form
func (z e_nocarry_0189) ToBigIntRegular(res *big.Int) *big.Int {
	z.FromMont()
	return z.ToBigInt(res)
}

// Bytes returns the value of z as a big-endian byte array
func (z *e_nocarry_0189) Bytes() (res [Limbs * 8]byte) {
	_z := z.ToRegular()
	binary.BigEndian.PutUint64(res[16:24], _z[0])
	binary.BigEndian.PutUint64(res[8:16], _z[1])
	binary.BigEndian.PutUint64(res[0:8], _z[2])

	return
}

// Marshal returns the value of z as a big-endian byte slice
func (z *e_nocarry_0189) Marshal() []byte {
	b := z.Bytes()
	return b[:]
}

// SetBytes interprets e as the bytes of a big-endian unsigned integer,
// sets z to that value, and returns z.
func (z *e_nocarry_0189) SetBytes(e []byte) *e_nocarry_0189 {
	// get a big int from our pool
	vv := bigIntPool.Get().(*big.Int)
	vv.SetBytes(e)

	// set big int
	z.SetBigInt(vv)

	// put temporary object back in pool
	bigIntPool.Put(vv)

	return z
}

// SetBigInt sets z to v and returns z
func (z *e_nocarry_0189) SetBigInt(v *big.Int) *e_nocarry_0189 {
	z.SetZero()

	var zero big.Int

	// fast path
	c := v.Cmp(&_modulus)
	if c == 0 {
		// v == 0
		return z
	} else if c != 1 && v.Cmp(&zero) != -1 {
		// 0 < v < q
		return z.setBigInt(v)
	}

	// get temporary big int from the pool
	vv := bigIntPool.Get().(*big.Int)

	// copy input + modular reduction
	vv.Set(v)
	vv.Mod(v, &_modulus)

	// set big int byte value
	z.setBigInt(vv)

	// release object into pool
	bigIntPool.Put(vv)
	return z
}

// setBigInt assumes 0 ⩽ v < q
func (z *e_nocarry_0189) setBigInt(v *big.Int) *e_nocarry_0189 {
	vBits := v.Bits()

	if bits.UintSize == 64 {
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
	} else {
		for i := 0; i < len(vBits); i++ {
			if i%2 == 0 {
				z[i/2] = uint64(vBits[i])
			} else {
				z[i/2] |= uint64(vBits[i]) << 32
			}
		}
	}

	return z.ToMont()
}

// SetString creates a big.Int with number and calls SetBigInt on z
//
// The number prefix determines the actual base: A prefix of
// ”0b” or ”0B” selects base 2, ”0”, ”0o” or ”0O” selects base 8,
// and ”0x” or ”0X” selects base 16. Otherwise, the selected base is 10
// and no prefix is accepted.
//
// For base 16, lower and upper case letters are considered the same:
// The letters 'a' to 'f' and 'A' to 'F' represent digit values 10 to 15.
//
// An underscore character ”_” may appear between a base
// prefix and an adjacent digit, and between successive digits; such
// underscores do not change the value of the number.
// Incorrect placement of underscores is reported as a panic if there
// are no other errors.
//
// If the number is invalid this method leaves z unchanged and returns nil, error.
func (z *e_nocarry_0189) SetString(number string) (*e_nocarry_0189, error) {
	// get temporary big int from the pool
	vv := bigIntPool.Get().(*big.Int)

	if _, ok := vv.SetString(number, 0); !ok {
		return nil, errors.New("e_nocarry_0189.SetString failed -> can't parse number into a big.Int " + number)
	}

	z.SetBigInt(vv)

	// release object into pool
	bigIntPool.Put(vv)

	return z, nil
}

// MarshalJSON returns json encoding of z (z.Text(10))
// If z == nil, returns null
func (z *e_nocarry_0189) MarshalJSON() ([]byte, error) {
	if z == nil {
		return []byte("null"), nil
	}
	const maxSafeBound = 15 // we encode it as number if it's small
	s := z.Text(10)
	if len(s) <= maxSafeBound {
		return []byte(s), nil
	}
	var sbb strings.Builder
	sbb.WriteByte('"')
	sbb.WriteString(s)
	sbb.WriteByte('"')
	return []byte(sbb.String()), nil
}

// UnmarshalJSON accepts numbers and strings as input
// See e_nocarry_0189.SetString for valid prefixes (0x, 0b, ...)
func (z *e_nocarry_0189) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) > Bits*3 {
		return errors.New("value too large (max = e_nocarry_0189.Bits * 3)")
	}

	// we accept numbers and strings, remove leading and trailing quotes if any
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	// get temporary big int from the pool
	vv := bigIntPool.Get().(*big.Int)

	if _, ok := vv.SetString(s, 0); !ok {
		return errors.New("can't parse into a big.Int: " + s)
	}

	z.SetBigInt(vv)

	// release object into pool
	bigIntPool.Put(vv)
	return nil
}

var (
	_bLegendreExponente_nocarry_0189 *big.Int
	_bSqrtExponente_nocarry_0189     *big.Int
)

func init() {
	_bLegendreExponente_nocarry_0189, _ = new(big.Int).SetString("ec0cee11613bbc8b95d928dfe2976a1b704348eec680567", 16)
	const sqrtExponente_nocarry_0189 = "76067708b09dde45caec946ff14bb50db821a47763402b4"
	_bSqrtExponente_nocarry_0189, _ = new(big.Int).SetString(sqrtExponente_nocarry_0189, 16)
}

// Legendre returns the Legendre symbol of z (either +1, -1, or 0.)
func (z *e_nocarry_0189) Legendre() int {
	var l e_nocarry_0189
	// z^((q-1)/2)
	l.Exp(*z, _bLegendreExponente_nocarry_0189)

	if l.IsZero() {
		return 0
	}

	// if l == 1
	if l.IsOne() {
		return 1
	}
	return -1
}

// Sqrt z = √x (mod q)
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *e_nocarry_0189) Sqrt(x *e_nocarry_0189) *e_nocarry_0189 {
	// q ≡ 3 (mod 4)
	// using  z ≡ ± x^((p+1)/4) (mod q)
	var y, square e_nocarry_0189
	y.Exp(*x, _bSqrtExponente_nocarry_0189)
	// as we didn't compute the legendre symbol, ensure we found y such that y * y = x
	square.Square(&y)
	if square.Equal(x) {
		return z.Set(&y)
	}
	return nil
}

const (
	k               = 32 // word size / 2
	signBitSelector = uint64(1) << 63
	approxLowBitsN  = k - 1
	approxHighBitsN = k + 1
)

const (
	inversionCorrectionFactorWord0 = 13189686649545512035
	inversionCorrectionFactorWord1 = 12553978508602752614
	inversionCorrectionFactorWord2 = 103218372913920363
	invIterationsN                 = 14
)

// Inverse z = x⁻¹ (mod q)
//
// if x == 0, sets and returns z = x
func (z *e_nocarry_0189) Inverse(x *e_nocarry_0189) *e_nocarry_0189 {
	// Implements "Optimized Binary GCD for Modular Inversion"
	// https://github.com/pornin/bingcd/blob/main/doc/bingcd.pdf

	a := *x
	b := e_nocarry_0189{
		q0,
		q1,
		q2,
	} // b := q

	u := e_nocarry_0189{1}

	// Update factors: we get [u; v] ← [f₀ g₀; f₁ g₁] [u; v]
	// cᵢ = fᵢ + 2³¹ - 1 + 2³² * (gᵢ + 2³¹ - 1)
	var c0, c1 int64

	// Saved update factors to reduce the number of field multiplications
	var pf0, pf1, pg0, pg1 int64

	var i uint

	var v, s e_nocarry_0189

	// Since u,v are updated every other iteration, we must make sure we terminate after evenly many iterations
	// This also lets us get away with half as many updates to u,v
	// To make this constant-time-ish, replace the condition with i < invIterationsN
	for i = 0; i&1 == 1 || !a.IsZero(); i++ {
		n := max(a.BitLen(), b.BitLen())
		aApprox, bApprox := approximate(&a, n), approximate(&b, n)

		// f₀, g₀, f₁, g₁ = 1, 0, 0, 1
		c0, c1 = updateFactorIdentityMatrixRow0, updateFactorIdentityMatrixRow1

		for j := 0; j < approxLowBitsN; j++ {

			// -2ʲ < f₀, f₁ ≤ 2ʲ
			// |f₀| + |f₁| < 2ʲ⁺¹

			if aApprox&1 == 0 {
				aApprox /= 2
			} else {
				s, borrow := bits.Sub64(aApprox, bApprox, 0)
				if borrow == 1 {
					s = bApprox - aApprox
					bApprox = aApprox
					c0, c1 = c1, c0
					// invariants unchanged
				}

				aApprox = s / 2
				c0 = c0 - c1

				// Now |f₀| < 2ʲ⁺¹ ≤ 2ʲ⁺¹ (only the weaker inequality is needed, strictly speaking)
				// Started with f₀ > -2ʲ and f₁ ≤ 2ʲ, so f₀ - f₁ > -2ʲ⁺¹
				// Invariants unchanged for f₁
			}

			c1 *= 2
			// -2ʲ⁺¹ < f₁ ≤ 2ʲ⁺¹
			// So now |f₀| + |f₁| < 2ʲ⁺²
		}

		s = a

		var g0 int64
		// from this point on c0 aliases for f0
		c0, g0 = updateFactorsDecompose(c0)
		aHi := a.linearCombNonModular(&s, c0, &b, g0)
		if aHi&signBitSelector != 0 {
			// if aHi < 0
			c0, g0 = -c0, -g0
			aHi = negL(&a, aHi)
		}
		// right-shift a by k-1 bits
		a[0] = (a[0] >> approxLowBitsN) | ((a[1]) << approxHighBitsN)
		a[1] = (a[1] >> approxLowBitsN) | ((a[2]) << approxHighBitsN)
		a[2] = (a[2] >> approxLowBitsN) | (aHi << approxHighBitsN)

		var f1 int64
		// from this point on c1 aliases for g0
		f1, c1 = updateFactorsDecompose(c1)
		bHi := b.linearCombNonModular(&s, f1, &b, c1)
		if bHi&signBitSelector != 0 {
			// if bHi < 0
			f1, c1 = -f1, -c1
			bHi = negL(&b, bHi)
		}
		// right-shift b by k-1 bits
		b[0] = (b[0] >> approxLowBitsN) | ((b[1]) << approxHighBitsN)
		b[1] = (b[1] >> approxLowBitsN) | ((b[2]) << approxHighBitsN)
		b[2] = (b[2] >> approxLowBitsN) | (bHi << approxHighBitsN)

		if i&1 == 1 {
			// Combine current update factors with previously stored ones
			// [F₀, G₀; F₁, G₁] ← [f₀, g₀; f₁, g₁] [pf₀, pg₀; pf₁, pg₁], with capital letters denoting new combined values
			// We get |F₀| = | f₀pf₀ + g₀pf₁ | ≤ |f₀pf₀| + |g₀pf₁| = |f₀| |pf₀| + |g₀| |pf₁| ≤ 2ᵏ⁻¹|pf₀| + 2ᵏ⁻¹|pf₁|
			// = 2ᵏ⁻¹ (|pf₀| + |pf₁|) < 2ᵏ⁻¹ 2ᵏ = 2²ᵏ⁻¹
			// So |F₀| < 2²ᵏ⁻¹ meaning it fits in a 2k-bit signed register

			// c₀ aliases f₀, c₁ aliases g₁
			c0, g0, f1, c1 = c0*pf0+g0*pf1,
				c0*pg0+g0*pg1,
				f1*pf0+c1*pf1,
				f1*pg0+c1*pg1

			s = u

			// 0 ≤ u, v < 2²⁵⁵
			// |F₀|, |G₀| < 2⁶³
			u.linearComb(&u, c0, &v, g0)
			// |F₁|, |G₁| < 2⁶³
			v.linearComb(&s, f1, &v, c1)

		} else {
			// Save update factors
			pf0, pg0, pf1, pg1 = c0, g0, f1, c1
		}
	}

	// For every iteration that we miss, v is not being multiplied by 2ᵏ⁻²
	const pSq uint64 = 1 << (2 * (k - 1))
	a = e_nocarry_0189{pSq}
	// If the function is constant-time ish, this loop will not run (no need to take it out explicitly)
	for ; i < invIterationsN; i += 2 {
		// could optimize further with mul by word routine or by pre-computing a table since with k=26,
		// we would multiply by pSq up to 13times;
		// on x86, the assembly routine outperforms generic code for mul by word
		// on arm64, we may loose up to ~5% for 6 limbs
		mul(&v, &v, &a)
	}

	u.Set(x) // for correctness check

	z.Mul(&v, &e_nocarry_0189{
		inversionCorrectionFactorWord0,
		inversionCorrectionFactorWord1,
		inversionCorrectionFactorWord2,
	})

	// correctness check
	v.Mul(&u, z)
	if !v.IsOne() && !u.IsZero() {
		return z.inverseExp(&u)
	}

	return z
}

// inverseExp computes z = x⁻¹ (mod q) = x**(q-2) (mod q)
func (z *e_nocarry_0189) inverseExp(x *e_nocarry_0189) *e_nocarry_0189 {
	qMinusTwo := Modulus()
	qMinusTwo.Sub(qMinusTwo, big.NewInt(2))
	return z.Exp(*x, qMinusTwo)
}

// approximate a big number x into a single 64 bit word using its uppermost and lowermost bits
// if x fits in a word as is, no approximation necessary
func approximate(x *e_nocarry_0189, nBits int) uint64 {

	if nBits <= 64 {
		return x[0]
	}

	const mask = (uint64(1) << (k - 1)) - 1 // k-1 ones
	lo := mask & x[0]

	hiWordIndex := (nBits - 1) / 64

	hiWordBitsAvailable := nBits - hiWordIndex*64
	hiWordBitsUsed := min(hiWordBitsAvailable, approxHighBitsN)

	mask_ := uint64(^((1 << (hiWordBitsAvailable - hiWordBitsUsed)) - 1))
	hi := (x[hiWordIndex] & mask_) << (64 - hiWordBitsAvailable)

	mask_ = ^(1<<(approxLowBitsN+hiWordBitsUsed) - 1)
	mid := (mask_ & x[hiWordIndex-1]) >> hiWordBitsUsed

	return lo | mid | hi
}

// linearComb z = xC * x + yC * y;
// 0 ≤ x, y < 2¹⁸⁹
// |xC|, |yC| < 2⁶³
func (z *e_nocarry_0189) linearComb(x *e_nocarry_0189, xC int64, y *e_nocarry_0189, yC int64) {
	// | (hi, z) | < 2 * 2⁶³ * 2¹⁸⁹ = 2²⁵³
	// therefore | hi | < 2⁶¹ ≤ 2⁶³
	hi := z.linearCombNonModular(x, xC, y, yC)
	z.montReduceSigned(z, hi)
}

// montReduceSigned z = (xHi * r + x) * r⁻¹ using the SOS algorithm
// Requires |xHi| < 2⁶³. Most significant bit of xHi is the sign bit.
func (z *e_nocarry_0189) montReduceSigned(x *e_nocarry_0189, xHi uint64) {
	const signBitRemover = ^signBitSelector
	mustNeg := xHi&signBitSelector != 0
	// the SOS implementation requires that most significant bit is 0
	// Let X be xHi*r + x
	// If X is negative we would have initially stored it as 2⁶⁴ r + X (à la 2's complement)
	xHi &= signBitRemover
	// with this a negative X is now represented as 2⁶³ r + X

	var t [2*Limbs - 1]uint64
	var C uint64

	m := x[0] * qInvNeg

	C = madd0(m, q0, x[0])
	C, t[1] = madd2(m, q1, x[1], C)
	C, t[2] = madd2(m, q2, x[2], C)

	// m * qElement[2] ≤ (2⁶⁴ - 1) * (2⁶³ - 1) = 2¹²⁷ - 2⁶⁴ - 2⁶³ + 1
	// x[2] + C ≤ 2*(2⁶⁴ - 1) = 2⁶⁵ - 2
	// On LHS, (C, t[2]) ≤ 2¹²⁷ - 2⁶⁴ - 2⁶³ + 1 + 2⁶⁵ - 2 = 2¹²⁷ + 2⁶³ - 1
	// So on LHS, C ≤ 2⁶³
	t[3] = xHi + C
	// xHi + C < 2⁶³ + 2⁶³ = 2⁶⁴

	// <standard SOS>
	{
		const i = 1
		m = t[i] * qInvNeg

		C = madd0(m, q0, t[i+0])
		C, t[i+1] = madd2(m, q1, t[i+1], C)
		C, t[i+2] = madd2(m, q2, t[i+2], C)

		t[i+Limbs] += C
	}
	{
		const i = 2
		m := t[i] * qInvNeg

		C = madd0(m, q0, t[i+0])
		C, z[0] = madd2(m, q1, t[i+1], C)
		z[2], z[1] = madd2(m, q2, t[i+2], C)
	}

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], _ = bits.Sub64(z[2], q2, b)
	}
	// </standard SOS>

	if mustNeg {
		// We have computed ( 2⁶³ r + X ) r⁻¹ = 2⁶³ + X r⁻¹ instead
		var b uint64
		z[0], b = bits.Sub64(z[0], signBitSelector, 0)
		z[1], b = bits.Sub64(z[1], 0, b)
		z[2], b = bits.Sub64(z[2], 0, b)

		// Occurs iff x == 0 && xHi < 0, i.e. X = rX' for -2⁶³ ≤ X' < 0

		if b != 0 {
			// z[2] = -1
			// negative: add q
			const neg1 = 0xFFFFFFFFFFFFFFFF

			var carry uint64

			z[0], carry = bits.Add64(z[0], q0, 0)
			z[1], carry = bits.Add64(z[1], q1, carry)
			z[2], _ = bits.Add64(neg1, q2, carry)
		}
	}
}

const (
	updateFactorsConversionBias    int64 = 0x7fffffff7fffffff // (2³¹ - 1)(2³² + 1)
	updateFactorIdentityMatrixRow0       = 1
	updateFactorIdentityMatrixRow1       = 1 << 32
)

func updateFactorsDecompose(c int64) (int64, int64) {
	c += updateFactorsConversionBias
	const low32BitsFilter int64 = 0xFFFFFFFF
	f := c&low32BitsFilter - 0x7FFFFFFF
	g := c>>32&low32BitsFilter - 0x7FFFFFFF
	return f, g
}

// negL negates in place [x | xHi] and return the new most significant word xHi
func negL(x *e_nocarry_0189, xHi uint64) uint64 {
	var b uint64

	x[0], b = bits.Sub64(0, x[0], 0)
	x[1], b = bits.Sub64(0, x[1], b)
	x[2], b = bits.Sub64(0, x[2], b)
	xHi, _ = bits.Sub64(0, xHi, b)

	return xHi
}

// mulWNonModular multiplies by one word in non-montgomery, without reducing
func (z *e_nocarry_0189) mulWNonModular(x *e_nocarry_0189, y int64) uint64 {

	// w := abs(y)
	m := y >> 63
	w := uint64((y ^ m) - m)

	var c uint64
	c, z[0] = bits.Mul64(x[0], w)
	c, z[1] = madd1(x[1], w, c)
	c, z[2] = madd1(x[2], w, c)

	if y < 0 {
		c = negL(z, c)
	}

	return c
}

// linearCombNonModular computes a linear combination without modular reduction
func (z *e_nocarry_0189) linearCombNonModular(x *e_nocarry_0189, xC int64, y *e_nocarry_0189, yC int64) uint64 {
	var yTimes e_nocarry_0189

	yHi := yTimes.mulWNonModular(y, yC)
	xHi := z.mulWNonModular(x, xC)

	var carry uint64
	z[0], carry = bits.Add64(z[0], yTimes[0], 0)
	z[1], carry = bits.Add64(z[1], yTimes[1], carry)
	z[2], carry = bits.Add64(z[2], yTimes[2], carry)

	yHi, _ = bits.Add64(xHi, yHi, carry)

	return yHi
}
