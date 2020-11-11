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

// Code generated by gurvy DO NOT EDIT

package bls377

import (
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/consensys/gurvy/bls377/fp"
	"github.com/consensys/gurvy/bls377/fr"
	"github.com/consensys/gurvy/utils"
	"github.com/consensys/gurvy/utils/parallel"
)

// G2Jac is a point with e2 coordinates
type G2Jac struct {
	X, Y, Z e2
}

// g2Proj point in projective coordinates
type g2Proj struct {
	x, y, z e2
}

// G2Affine point in affine coordinates
type G2Affine struct {
	X, Y e2
}

// AddAssign point addition in montgomery form
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#addition-add-2007-bl
func (p *G2Jac) AddAssign(a *G2Jac) *G2Jac {

	// p is infinity, return a
	if p.Z.IsZero() {
		p.Set(a)
		return p
	}

	// a is infinity, return p
	if a.Z.IsZero() {
		return p
	}

	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V e2
	Z1Z1.Square(&a.Z)
	Z2Z2.Square(&p.Z)
	U1.Mul(&a.X, &Z2Z2)
	U2.Mul(&p.X, &Z1Z1)
	S1.Mul(&a.Y, &p.Z).
		Mul(&S1, &Z2Z2)
	S2.Mul(&p.Y, &a.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U1.Equal(&U2) && S1.Equal(&S2) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &U1)
	I.Double(&H).
		Square(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &S1).Double(&r)
	V.Mul(&U1, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	S1.Mul(&S1, &J).Double(&S1)
	p.Y.Sub(&p.Y, &S1)
	p.Z.Add(&p.Z, &a.Z)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &Z2Z2).
		Mul(&p.Z, &H)

	return p
}

// AddMixed point addition
// http://www.hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-0.html#addition-madd-2007-bl
func (p *G2Jac) AddMixed(a *G2Affine) *G2Jac {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.Z.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Z.SetOne()
		return p
	}

	// get some Element from our pool
	var Z1Z1, U2, S2, H, HH, I, J, r, V e2
	Z1Z1.Square(&p.Z)
	U2.Mul(&a.X, &Z1Z1)
	S2.Mul(&a.Y, &p.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U2.Equal(&p.X) && S2.Equal(&p.Y) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &p.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &p.Y).Double(&r)
	V.Mul(&p.X, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	J.Mul(&J, &p.Y).Double(&J)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	p.Y.Sub(&p.Y, &J)
	p.Z.Add(&p.Z, &H)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &HH)

	return p
}

// Double doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) Double(q *G2Jac) *G2Jac {
	p.Set(q)
	p.DoubleAssign()
	return p
}

// DoubleAssign doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g2p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) DoubleAssign() *G2Jac {

	// get some Element from our pool
	var XX, YY, YYYY, ZZ, S, M, T e2

	XX.Square(&p.X)
	YY.Square(&p.Y)
	YYYY.Square(&YY)
	ZZ.Square(&p.Z)
	S.Add(&p.X, &YY)
	S.Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).Add(&M, &XX)
	p.Z.Add(&p.Z, &p.Y).
		Square(&p.Z).
		Sub(&p.Z, &YY).
		Sub(&p.Z, &ZZ)
	T.Square(&M)
	p.X = T
	T.Double(&S)
	p.X.Sub(&p.X, &T)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).Double(&YYYY).Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)

	return p
}

// ScalarMultiplication computes and returns p = a*s
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G2Jac) ScalarMultiplication(a *G2Jac, s *big.Int) *G2Jac {
	return p.mulGLV(a, s)
}

// ScalarMultiplication computes and returns p = a*s
func (p *G2Affine) ScalarMultiplication(a *G2Affine, s *big.Int) *G2Affine {
	var _p G2Jac
	_p.FromAffine(a)
	_p.mulGLV(&_p, s)
	p.FromJacobian(&_p)
	return p
}

// Set set p to the provided point
func (p *G2Jac) Set(a *G2Jac) *G2Jac {
	p.X, p.Y, p.Z = a.X, a.Y, a.Z
	return p
}

// Equal tests if two points (in Jacobian coordinates) are equal
func (p *G2Jac) Equal(a *G2Jac) bool {

	if p.Z.IsZero() && a.Z.IsZero() {
		return true
	}
	_p := G2Affine{}
	_p.FromJacobian(p)

	_a := G2Affine{}
	_a.FromJacobian(a)

	return _p.X.Equal(&_a.X) && _p.Y.Equal(&_a.Y)
}

// Equal tests if two points (in Affine coordinates) are equal
func (p *G2Affine) Equal(a *G2Affine) bool {
	return p.X.Equal(&a.X) && p.Y.Equal(&a.Y)
}

// Neg computes -G
func (p *G2Jac) Neg(a *G2Jac) *G2Jac {
	*p = *a
	p.Y.Neg(&a.Y)
	return p
}

// Neg computes -G
func (p *G2Affine) Neg(a *G2Affine) *G2Affine {
	p.X = a.X
	p.Y.Neg(&a.Y)
	return p
}

// SubAssign substracts two points on the curve
func (p *G2Jac) SubAssign(a *G2Jac) *G2Jac {
	var tmp G2Jac
	tmp.Set(a)
	tmp.Y.Neg(&tmp.Y)
	p.AddAssign(&tmp)
	return p
}

// FromJacobian rescale a point in Jacobian coord in z=1 plane
func (p *G2Affine) FromJacobian(p1 *G2Jac) *G2Affine {

	var a, b e2

	if p1.Z.IsZero() {
		p.X.SetZero()
		p.Y.SetZero()
		return p
	}

	a.Inverse(&p1.Z)
	b.Square(&a)
	p.X.Mul(&p1.X, &b)
	p.Y.Mul(&p1.Y, &b).Mul(&p.Y, &a)

	return p
}

// FromJacobian converts a point from Jacobian to projective coordinates
func (p *g2Proj) FromJacobian(Q *G2Jac) *g2Proj {
	// memalloc
	var buf e2
	buf.Square(&Q.Z)

	p.x.Mul(&Q.X, &Q.Z)
	p.y.Set(&Q.Y)
	p.z.Mul(&Q.Z, &buf)

	return p
}

func (p *G2Jac) String() string {
	if p.Z.IsZero() {
		return "O"
	}
	_p := G2Affine{}
	_p.FromJacobian(p)
	return "E([" + _p.X.String() + "," + _p.Y.String() + "]),"
}

// FromAffine sets p = Q, p in Jacboian, Q in affine
func (p *G2Jac) FromAffine(Q *G2Affine) *G2Jac {
	if Q.X.IsZero() && Q.Y.IsZero() {
		p.Z.SetZero()
		p.X.SetOne()
		p.Y.SetOne()
		return p
	}
	p.Z.SetOne()
	p.X.Set(&Q.X)
	p.Y.Set(&Q.Y)
	return p
}

func (p *G2Affine) String() string {
	var x, y e2
	x.Set(&p.X)
	y.Set(&p.Y)
	return "E([" + x.String() + "," + y.String() + "]),"
}

// IsInfinity checks if the point is infinity (in affine, it's encoded as (0,0))
func (p *G2Affine) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// IsOnCurve returns true if p in on the curve
func (p *G2Jac) IsOnCurve() bool {
	var left, right, tmp e2
	left.Square(&p.Y)
	right.Square(&p.X).Mul(&right, &p.X)
	tmp.Square(&p.Z).
		Square(&tmp).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &bTwistCurveCoeff)
	right.Add(&right, &tmp)
	return left.Equal(&right)
}

// IsOnCurve returns true if p in on the curve
func (p *G2Affine) IsOnCurve() bool {
	var point G2Jac
	point.FromAffine(p)
	return point.IsOnCurve() // call this function to handle infinity point
}

// IsInSubGroup returns true if p is in the correct subgroup, false otherwise
func (p *G2Affine) IsInSubGroup() bool {
	var _p G2Jac
	_p.FromAffine(p)
	return _p.IsOnCurve() && _p.IsInSubGroup()
}

// IsInSubGroup returns true if p is on the r-torsion, false otherwise.
// Z[r,0]+Z[-lambdaG2, 1] is the kernel
// of (u,v)->u+lambdaG2v mod r. Expressing r, lambdaG2 as
// polynomials in x, a short vector of this Zmodule is
// 1, x**2. So we check that p+x**2*phi(p)
// is the infinity.
func (p *G2Jac) IsInSubGroup() bool {

	var res G2Jac
	res.phi(p).
		ScalarMultiplication(&res, &xGen).
		ScalarMultiplication(&res, &xGen).
		AddAssign(p)

	return res.IsOnCurve() && res.Z.IsZero()

}

// mulWindowed 2-bits windowed exponentiation
func (p *G2Jac) mulWindowed(a *G2Jac, s *big.Int) *G2Jac {

	var res G2Jac
	var ops [3]G2Jac

	res.Set(&g2Infinity)
	ops[0].Set(a)
	ops[1].Double(&ops[0])
	ops[2].Set(&ops[0]).AddAssign(&ops[1])

	b := s.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0xc0)
		for j := 0; j < 4; j++ {
			res.DoubleAssign().DoubleAssign()
			c := (w & mask) >> (6 - 2*j)
			if c != 0 {
				res.AddAssign(&ops[c-1])
			}
			mask = mask >> 2
		}
	}
	p.Set(&res)

	return p

}

// psi(p) = u o frob o u**-1 where u:E'->E iso from the twist to E
func (p *G2Jac) psi(a *G2Jac) *G2Jac {
	p.Set(a)
	p.X.Conjugate(&p.X).Mul(&p.X, &endo.u)
	p.Y.Conjugate(&p.Y).Mul(&p.Y, &endo.v)
	p.Z.Conjugate(&p.Z)
	return p
}

// phi assigns p to phi(a) where phi: (x,y)->(ux,y), and returns p
func (p *G2Jac) phi(a *G2Jac) *G2Jac {
	p.Set(a)

	p.X.MulByElement(&p.X, &thirdRootOneG2)

	return p
}

// mulGLV performs scalar multiplication using GLV
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G2Jac) mulGLV(a *G2Jac, s *big.Int) *G2Jac {

	var table [15]G2Jac
	var zero big.Int
	var res G2Jac
	var k1, k2 fr.Element

	res.Set(&g2Infinity)

	// table[b3b2b1b0-1] = b3b2*phi(a) + b1b0*a
	table[0].Set(a)
	table[3].phi(a)

	// split the scalar, modifies +-a, phi(a) accordingly
	k := utils.SplitScalar(s, &glvBasis)

	if k[0].Cmp(&zero) == -1 {
		k[0].Neg(&k[0])
		table[0].Neg(&table[0])
	}
	if k[1].Cmp(&zero) == -1 {
		k[1].Neg(&k[1])
		table[3].Neg(&table[3])
	}

	// precompute table (2 bits sliding window)
	// table[b3b2b1b0-1] = b3b2*phi(a) + b1b0*a if b3b2b1b0 != 0
	table[1].Double(&table[0])
	table[2].Set(&table[1]).AddAssign(&table[0])
	table[4].Set(&table[3]).AddAssign(&table[0])
	table[5].Set(&table[3]).AddAssign(&table[1])
	table[6].Set(&table[3]).AddAssign(&table[2])
	table[7].Double(&table[3])
	table[8].Set(&table[7]).AddAssign(&table[0])
	table[9].Set(&table[7]).AddAssign(&table[1])
	table[10].Set(&table[7]).AddAssign(&table[2])
	table[11].Set(&table[7]).AddAssign(&table[3])
	table[12].Set(&table[11]).AddAssign(&table[0])
	table[13].Set(&table[11]).AddAssign(&table[1])
	table[14].Set(&table[11]).AddAssign(&table[2])

	// bounds on the lattice base vectors guarantee that k1, k2 are len(r)/2 bits long max
	k1.SetBigInt(&k[0]).FromMont()
	k2.SetBigInt(&k[1]).FromMont()

	// loop starts from len(k1)/2 due to the bounds
	for i := len(k1)/2 - 1; i >= 0; i-- {
		mask := uint64(3) << 62
		for j := 0; j < 32; j++ {
			res.Double(&res).Double(&res)
			b1 := (k1[i] & mask) >> (62 - 2*j)
			b2 := (k2[i] & mask) >> (62 - 2*j)
			if b1|b2 != 0 {
				s := (b2<<2 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 2
		}
	}

	p.Set(&res)
	return p
}

// BatchScalarMultiplicationG2 multiplies the same base (generator) by all scalars
// and return resulting points in affine coordinates
// uses a simple windowed-NAF like exponentiation algorithm
func BatchScalarMultiplicationG2(base *G2Affine, scalars []fr.Element) []G2Affine {

	// approximate cost in group ops is
	// cost = 2^{c-1} + n(scalar.nbBits+nbChunks)

	nbPoints := uint64(len(scalars))
	min := ^uint64(0)
	bestC := 0
	for c := 2; c < 18; c++ {
		cost := uint64(1 << (c - 1))
		nbChunks := uint64(fr.Limbs * 64 / c)
		if (fr.Limbs*64)%c != 0 {
			nbChunks++
		}
		cost += nbPoints * ((fr.Limbs * 64) + nbChunks)
		if cost < min {
			min = cost
			bestC = c
		}
	}
	c := uint64(bestC) // window size
	nbChunks := int(fr.Limbs * 64 / c)
	if (fr.Limbs*64)%c != 0 {
		nbChunks++
	}
	mask := uint64((1 << c) - 1) // low c bits are 1
	msbWindow := uint64(1 << (c - 1))

	// precompute all powers of base for our window
	// note here that if performance is critical, we can implement as in the msmX methods
	// this allocation to be on the stack
	baseTable := make([]G2Jac, (1 << (c - 1)))
	baseTable[0].Set(&g2Infinity)
	baseTable[0].AddMixed(base)
	for i := 1; i < len(baseTable); i++ {
		baseTable[i] = baseTable[i-1]
		baseTable[i].AddMixed(base)
	}

	pScalars := partitionScalars(scalars, c)

	// compute offset and word selector / shift to select the right bits of our windows
	selectors := make([]selector, nbChunks)
	for chunk := 0; chunk < nbChunks; chunk++ {
		jc := uint64(uint64(chunk) * c)
		d := selector{}
		d.index = jc / 64
		d.shift = jc - (d.index * 64)
		d.mask = mask << d.shift
		d.multiWordSelect = (64%c) != 0 && d.shift > (64-c) && d.index < (fr.Limbs-1)
		if d.multiWordSelect {
			nbBitsHigh := d.shift - uint64(64-c)
			d.maskHigh = (1 << nbBitsHigh) - 1
			d.shiftHigh = (c - nbBitsHigh)
		}
		selectors[chunk] = d
	}

	toReturn := make([]G2Affine, len(scalars))

	// for each digit, take value in the base table, double it c time, voila.
	parallel.Execute(len(pScalars), func(start, end int) {
		var p G2Jac
		for i := start; i < end; i++ {
			p.Set(&g2Infinity)
			for chunk := nbChunks - 1; chunk >= 0; chunk-- {
				s := selectors[chunk]
				if chunk != nbChunks-1 {
					for j := uint64(0); j < c; j++ {
						p.DoubleAssign()
					}
				}

				bits := (pScalars[i][s.index] & s.mask) >> s.shift
				if s.multiWordSelect {
					bits += (pScalars[i][s.index+1] & s.maskHigh) << s.shiftHigh
				}

				if bits == 0 {
					continue
				}

				// if msbWindow bit is set, we need to substract
				if bits&msbWindow == 0 {
					// add

					p.AddAssign(&baseTable[bits-1])

				} else {
					// sub

					t := baseTable[bits & ^msbWindow]
					t.Neg(&t)
					p.AddAssign(&t)

				}
			}

			// set our result point

			toReturn[i].FromJacobian(&p)

		}
	})

	return toReturn

}

// SizeOfG2Compressed represents the size in bytes that a G2Affine need in binary form, compressed
const SizeOfG2Compressed = 48 * 2

// SizeOfG2Uncompressed represents the size in bytes that a G2Affine need in binary form, uncompressed
const SizeOfG2Uncompressed = SizeOfG2Compressed * 2

// Bytes returns binary representation of p
// will store X coordinate in regular form and a parity bit
// we follow the BLS381 style encoding as specified in ZCash and now IETF
// The most significant bit, when set, indicates that the point is in compressed form. Otherwise, the point is in uncompressed form.
// The second-most significant bit indicates that the point is at infinity. If this bit is set, the remaining bits of the group element's encoding should be set to zero.
// The third-most significant bit is set if (and only if) this point is in compressed form and it is not the point at infinity and its y-coordinate is the lexicographically largest of the two associated with the encoded x-coordinate.
func (p *G2Affine) Bytes() (res [SizeOfG2Compressed]byte) {

	// check if p is infinity point
	if p.X.IsZero() && p.Y.IsZero() {
		res[0] = mCompressedInfinity
		return
	}

	// tmp is used to convert from montgomery representation to regular
	var tmp fp.Element

	msbMask := mCompressedSmallest
	// compressed, we need to know if Y is lexicographically bigger than -Y
	// if p.Y ">" -p.Y
	if p.Y.LexicographicallyLargest() {
		msbMask = mCompressedLargest
	}

	// we store X  and mask the most significant word with our metadata mask
	// p.X.A1 | p.X.A0
	tmp = p.X.A0
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[88:96], tmp[0])
	binary.BigEndian.PutUint64(res[80:88], tmp[1])
	binary.BigEndian.PutUint64(res[72:80], tmp[2])
	binary.BigEndian.PutUint64(res[64:72], tmp[3])
	binary.BigEndian.PutUint64(res[56:64], tmp[4])
	binary.BigEndian.PutUint64(res[48:56], tmp[5])

	tmp = p.X.A1
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[40:48], tmp[0])
	binary.BigEndian.PutUint64(res[32:40], tmp[1])
	binary.BigEndian.PutUint64(res[24:32], tmp[2])
	binary.BigEndian.PutUint64(res[16:24], tmp[3])
	binary.BigEndian.PutUint64(res[8:16], tmp[4])
	binary.BigEndian.PutUint64(res[0:8], tmp[5])

	res[0] |= msbMask

	return
}

// RawBytes returns binary representation of p (stores X and Y coordinate)
// see Bytes() for a compressed representation
func (p *G2Affine) RawBytes() (res [SizeOfG2Uncompressed]byte) {

	// check if p is infinity point
	if p.X.IsZero() && p.Y.IsZero() {

		res[0] = mUncompressedInfinity

		return
	}

	// tmp is used to convert from montgomery representation to regular
	var tmp fp.Element

	// not compressed
	// we store the Y coordinate
	// p.Y.A1 | p.Y.A0
	tmp = p.Y.A0
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[184:192], tmp[0])
	binary.BigEndian.PutUint64(res[176:184], tmp[1])
	binary.BigEndian.PutUint64(res[168:176], tmp[2])
	binary.BigEndian.PutUint64(res[160:168], tmp[3])
	binary.BigEndian.PutUint64(res[152:160], tmp[4])
	binary.BigEndian.PutUint64(res[144:152], tmp[5])

	tmp = p.Y.A1
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[136:144], tmp[0])
	binary.BigEndian.PutUint64(res[128:136], tmp[1])
	binary.BigEndian.PutUint64(res[120:128], tmp[2])
	binary.BigEndian.PutUint64(res[112:120], tmp[3])
	binary.BigEndian.PutUint64(res[104:112], tmp[4])
	binary.BigEndian.PutUint64(res[96:104], tmp[5])

	// we store X  and mask the most significant word with our metadata mask
	// p.X.A1 | p.X.A0
	tmp = p.X.A0
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[88:96], tmp[0])
	binary.BigEndian.PutUint64(res[80:88], tmp[1])
	binary.BigEndian.PutUint64(res[72:80], tmp[2])
	binary.BigEndian.PutUint64(res[64:72], tmp[3])
	binary.BigEndian.PutUint64(res[56:64], tmp[4])
	binary.BigEndian.PutUint64(res[48:56], tmp[5])

	tmp = p.X.A1
	tmp.FromMont()
	binary.BigEndian.PutUint64(res[40:48], tmp[0])
	binary.BigEndian.PutUint64(res[32:40], tmp[1])
	binary.BigEndian.PutUint64(res[24:32], tmp[2])
	binary.BigEndian.PutUint64(res[16:24], tmp[3])
	binary.BigEndian.PutUint64(res[8:16], tmp[4])
	binary.BigEndian.PutUint64(res[0:8], tmp[5])

	res[0] |= mUncompressed

	return
}

// SetBytes sets p from binary representation in buf and returns number of consumed bytes
// bytes in buf must match either RawBytes() or Bytes() output
// if buf is too short io.ErrShortBuffer is returned
// if buf contains compressed representation (output from Bytes()) and we're unable to compute
// the Y coordinate (i.e the square root doesn't exist) this function retunrs an error
// note that this doesn't check if the resulting point is on the curve or in the correct subgroup
func (p *G2Affine) SetBytes(buf []byte) (int, error) {
	if len(buf) < SizeOfG2Compressed {
		return 0, io.ErrShortBuffer
	}

	// most significant byte
	mData := buf[0] & mMask
	buf[0] &= ^mMask // clear meta data

	// check buffer size
	if (mData == mUncompressed) || (mData == mUncompressedInfinity) {
		if len(buf) < SizeOfG2Uncompressed {
			return 0, io.ErrShortBuffer
		}
	}

	// if infinity is encoded in the metadata, we don't need to read the buffer
	if mData == mCompressedInfinity {
		p.X.SetZero()
		p.Y.SetZero()
		return SizeOfG2Compressed, nil
	}
	if mData == mUncompressedInfinity {
		p.X.SetZero()
		p.Y.SetZero()
		return SizeOfG2Uncompressed, nil
	}

	// tmp is used to convert to montgomery representation
	var tmp fp.Element

	// read X coordinate
	// p.X.A1 | p.X.A0
	tmp[0] = binary.BigEndian.Uint64(buf[88:96])
	tmp[1] = binary.BigEndian.Uint64(buf[80:88])
	tmp[2] = binary.BigEndian.Uint64(buf[72:80])
	tmp[3] = binary.BigEndian.Uint64(buf[64:72])
	tmp[4] = binary.BigEndian.Uint64(buf[56:64])
	tmp[5] = binary.BigEndian.Uint64(buf[48:56])
	tmp.ToMont()
	p.X.A0.Set(&tmp)

	tmp[0] = binary.BigEndian.Uint64(buf[40:48])
	tmp[1] = binary.BigEndian.Uint64(buf[32:40])
	tmp[2] = binary.BigEndian.Uint64(buf[24:32])
	tmp[3] = binary.BigEndian.Uint64(buf[16:24])
	tmp[4] = binary.BigEndian.Uint64(buf[8:16])
	tmp[5] = binary.BigEndian.Uint64(buf[0:8])
	tmp.ToMont()
	p.X.A1.Set(&tmp)

	// uncompressed point
	if mData == mUncompressed {
		// read Y coordinate
		// p.Y.A1 | p.Y.A0
		tmp[0] = binary.BigEndian.Uint64(buf[184:192])
		tmp[1] = binary.BigEndian.Uint64(buf[176:184])
		tmp[2] = binary.BigEndian.Uint64(buf[168:176])
		tmp[3] = binary.BigEndian.Uint64(buf[160:168])
		tmp[4] = binary.BigEndian.Uint64(buf[152:160])
		tmp[5] = binary.BigEndian.Uint64(buf[144:152])
		tmp.ToMont()
		p.Y.A0.Set(&tmp)

		tmp[0] = binary.BigEndian.Uint64(buf[136:144])
		tmp[1] = binary.BigEndian.Uint64(buf[128:136])
		tmp[2] = binary.BigEndian.Uint64(buf[120:128])
		tmp[3] = binary.BigEndian.Uint64(buf[112:120])
		tmp[4] = binary.BigEndian.Uint64(buf[104:112])
		tmp[5] = binary.BigEndian.Uint64(buf[96:104])
		tmp.ToMont()
		p.Y.A1.Set(&tmp)

		return SizeOfG2Uncompressed, nil
	}

	// we have a compressed coordinate, we need to solve the curve equation to compute Y
	var YSquared, Y e2

	YSquared.Square(&p.X).Mul(&YSquared, &p.X)
	YSquared.Add(&YSquared, &bTwistCurveCoeff)
	if YSquared.Legendre() == -1 {
		return 0, errors.New("invalid compressed coordinate: square root doesn't exist")
	}
	Y.Sqrt(&YSquared)

	if Y.LexicographicallyLargest() {
		// Y ">" -Y
		if mData == mCompressedSmallest {
			Y.Neg(&Y)
		}
	} else {
		// Y "<=" -Y
		if mData == mCompressedLargest {
			Y.Neg(&Y)
		}
	}

	p.Y = Y

	return SizeOfG2Compressed, nil
}

// unsafeComputeY called by Decoder when processing slices of compressed point in parallel (step 2)
// it computes the Y coordinate from the already set X coordinate and is compute intensive
func (p *G2Affine) unsafeComputeY() error {
	// stored in unsafeSetCompressedBytes

	mData := byte(p.Y.A0[0])

	// we have a compressed coordinate, we need to solve the curve equation to compute Y
	var YSquared, Y e2

	YSquared.Square(&p.X).Mul(&YSquared, &p.X)
	YSquared.Add(&YSquared, &bTwistCurveCoeff)
	if YSquared.Legendre() == -1 {
		return errors.New("invalid compressed coordinate: square root doesn't exist")
	}
	Y.Sqrt(&YSquared)

	if Y.LexicographicallyLargest() {
		// Y ">" -Y
		if mData == mCompressedSmallest {
			Y.Neg(&Y)
		}
	} else {
		// Y "<=" -Y
		if mData == mCompressedLargest {
			Y.Neg(&Y)
		}
	}

	p.Y = Y

	return nil
}

// unsafeSetCompressedBytes is called by Decoder when processing slices of compressed point in parallel (step 1)
// assumes buf[:8] mask is set to compressed
// returns true if point is infinity and need no further processing
// it sets X coordinate and uses Y for scratch space to store decompression metadata
func (p *G2Affine) unsafeSetCompressedBytes(buf []byte) (isInfinity bool) {

	// read the most significant byte
	mData := buf[0] & mMask
	buf[0] &= ^mMask

	if mData == mCompressedInfinity {
		p.X.SetZero()
		p.Y.SetZero()
		isInfinity = true
		return
	}

	// read X

	// tmp is used to convert to montgomery representation
	var tmp fp.Element

	// read X coordinate
	// p.X.A1 | p.X.A0
	tmp[0] = binary.BigEndian.Uint64(buf[88:96])
	tmp[1] = binary.BigEndian.Uint64(buf[80:88])
	tmp[2] = binary.BigEndian.Uint64(buf[72:80])
	tmp[3] = binary.BigEndian.Uint64(buf[64:72])
	tmp[4] = binary.BigEndian.Uint64(buf[56:64])
	tmp[5] = binary.BigEndian.Uint64(buf[48:56])
	tmp.ToMont()
	p.X.A0.Set(&tmp)

	tmp[0] = binary.BigEndian.Uint64(buf[40:48])
	tmp[1] = binary.BigEndian.Uint64(buf[32:40])
	tmp[2] = binary.BigEndian.Uint64(buf[24:32])
	tmp[3] = binary.BigEndian.Uint64(buf[16:24])
	tmp[4] = binary.BigEndian.Uint64(buf[8:16])
	tmp[5] = binary.BigEndian.Uint64(buf[0:8])
	tmp.ToMont()
	p.X.A1.Set(&tmp)

	// store mData in p.Y.A0[0]
	p.Y.A0[0] = uint64(mData)

	// recomputing Y will be done asynchronously
	return
}
