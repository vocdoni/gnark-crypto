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

package bls12378

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-378/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-378/internal/fptower"
)

// mapToCurve2 implements the Shallue and van de Woestijne method, applicable to any elliptic curve in Weierstrass form
// No cofactor clearing or isogeny
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-hash-to-curve-14#appendix-F.1
func mapToCurve2(u *fptower.E2) G2Affine {
	var tv1, tv2, tv3, tv4 fptower.E2
	var x1, x2, x3, gx1, gx2, gx, x, y fptower.E2
	var one fptower.E2
	var gx1NotSquare, gx1SquareOrGx2Not int

	//constants
	//c1 = g(Z)
	//c2 = -Z / 2
	//c3 = sqrt(-g(Z) * (3 * Z² + 4 * A))     # sgn0(c3) MUST equal 0
	//c4 = -4 * g(Z) / (3 * Z² + 4 * A)

	//Z  = 1
	//c1 = 1 + x
	//c2 = 302624103037653085866624240790900480369923845885462456876760372017370467951700652388141901174418655585487141470208
	//c3 = 103101875904467895810933718736441438872588282074344974385644588750905787716875163227630481701084529630905410327018 + 532069566337298640900882650392420558794316228756974744467711245052783922718532692426384643917730651223755366771180*x
	//c4 = 201749402025102057244416160527266986913282563923641637917840248011580311967800434925427934116279103723658094313471 + 201749402025102057244416160527266986913282563923641637917840248011580311967800434925427934116279103723658094313471*x

	//TODO: Move outside function?
	Z := fptower.E2{
		A0: fp.Element{1481365419032838079, 10045892448872562649, 7242180086616818316, 8832319421896135475, 13356930855120736188, 28498675542444634},
		A1: fp.Element{0},
	}
	c1 := fptower.E2{
		A0: fp.Element{1481365419032838079, 10045892448872562649, 7242180086616818316, 8832319421896135475, 13356930855120736188, 28498675542444634},
		A1: fp.Element{1481365419032838079, 10045892448872562649, 7242180086616818316, 8832319421896135475, 13356930855120736188, 28498675542444634},
	}
	c2 := fptower.E2{
		A0: fp.Element{14005317430843277345, 11643745377477984275, 11080596138069871993, 340432435853690873, 14787289713583717363, 127429472983909274},
		A1: fp.Element{0},
	}
	c3 := fptower.E2{
		A0: fp.Element{835112817875701876, 5875037160571017383, 16327530373670413863, 6452190755966272550, 1307167413516723115, 195508019089309651},
		A1: fp.Element{6791728344683345066, 16970244505711969686, 6187856097244537483, 15212952921759718760, 16780173455074099752, 22826623786084456},
	}
	c4 := fptower.E2{
		A0: fp.Element{17686179628435811074, 8827732204055603934, 3797093435112099908, 714611657777348053, 10811765714697799025, 150906846950249276},
		A1: fp.Element{17686179628435811074, 8827732204055603934, 3797093435112099908, 714611657777348053, 10811765714697799025, 150906846950249276},
	}

	one.SetOne()

	tv1.Square(u)       //    1.  tv1 = u²
	tv1.Mul(&tv1, &c1)  //    2.  tv1 = tv1 * c1
	tv2.Add(&one, &tv1) //    3.  tv2 = 1 + tv1
	tv1.Sub(&one, &tv1) //    4.  tv1 = 1 - tv1
	tv3.Mul(&tv1, &tv2) //    5.  tv3 = tv1 * tv2

	tv3.Inverse(&tv3)   //    6.  tv3 = inv0(tv3)
	tv4.Mul(u, &tv1)    //    7.  tv4 = u * tv1
	tv4.Mul(&tv4, &tv3) //    8.  tv4 = tv4 * tv3
	tv4.Mul(&tv4, &c3)  //    9.  tv4 = tv4 * c3
	x1.Sub(&c2, &tv4)   //    10.  x1 = c2 - tv4

	gx1.Square(&x1) //    11. gx1 = x1²
	//TODO: Beware A ≠ 0
	//12. gx1 = gx1 + A
	gx1.Mul(&gx1, &x1)                 //    13. gx1 = gx1 * x1
	gx1.Add(&gx1, &bTwistCurveCoeff)   //    14. gx1 = gx1 + B
	gx1NotSquare = gx1.Legendre() >> 1 //    15.  e1 = is_square(gx1)
	// gx1NotSquare = 0 if gx1 is a square, -1 otherwise

	x2.Add(&c2, &tv4) //    16.  x2 = c2 + tv4
	gx2.Square(&x2)   //    17. gx2 = x2²
	//    18. gx2 = gx2 + A
	gx2.Mul(&gx2, &x2)               //    19. gx2 = gx2 * x2
	gx2.Add(&gx2, &bTwistCurveCoeff) //    20. gx2 = gx2 + B

	{
		gx2NotSquare := gx2.Legendre() >> 1              // gx2Square = 0 if gx2 is a square, -1 otherwise
		gx1SquareOrGx2Not = gx2NotSquare | ^gx1NotSquare //    21.  e2 = is_square(gx2) AND NOT e1   # Avoid short-circuit logic ops
	}

	x3.Square(&tv2)   //    22.  x3 = tv2²
	x3.Mul(&x3, &tv3) //    23.  x3 = x3 * tv3
	x3.Square(&x3)    //    24.  x3 = x3²
	x3.Mul(&x3, &c4)  //    25.  x3 = x3 * c4

	x3.Add(&x3, &Z)                  //    26.  x3 = x3 + Z
	x.Select(gx1NotSquare, &x1, &x3) //    27.   x = CMOV(x3, x1, e1)   # x = x1 if gx1 is square, else x = x3
	// Select x1 iff gx1 is square iff gx1NotSquare = 0
	x.Select(gx1SquareOrGx2Not, &x2, &x) //    28.   x = CMOV(x, x2, e2)    # x = x2 if gx2 is square and gx1 is not
	// Select x2 iff gx2 is square and gx1 is not, iff gx1SquareOrGx2Not = 0
	gx.Square(&x) //    29.  gx = x²
	//    30.  gx = gx + A

	gx.Mul(&gx, &x)                //    31.  gx = gx * x
	gx.Add(&gx, &bTwistCurveCoeff) //    32.  gx = gx + B

	y.Sqrt(&gx)                             //    33.   y = sqrt(gx)
	signsNotEqual := g2Sgn0(u) ^ g2Sgn0(&y) //    34.  e3 = sgn0(u) == sgn0(y)

	tv1.Neg(&y)
	y.Select(int(signsNotEqual), &y, &tv1) //    35.   y = CMOV(-y, y, e3)       # Select correct sign of y
	return G2Affine{x, y}
}

// g2Sgn0 is an algebraic substitute for the notion of sign in ordered fields
// Namely, every non-zero quadratic residue in a finite field of characteristic =/= 2 has exactly two square roots, one of each sign
// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/ section 4.1
// The sign of an element is not obviously related to that of its Montgomery form
func g2Sgn0(z *fptower.E2) uint64 {

	nonMont := *z
	nonMont.FromMont()

	sign := uint64(0)
	zero := uint64(1)
	var signI uint64
	var zeroI uint64

	signI = nonMont.A0[0] % 2
	sign = sign | (zero & signI)

	zeroI = g1NotZero(&nonMont.A0)
	zeroI = 1 ^ (zeroI|-zeroI)>>63
	zero = zero & zeroI

	signI = nonMont.A1[0] % 2
	sign = sign | (zero & signI)

	return sign

}

// MapToG2 invokes the SVDW map, and guarantees that the result is in g2
func MapToG2(u fptower.E2) G2Affine {
	res := mapToCurve2(&u)
	res.ClearCofactor(&res)
	return res
}

// EncodeToG2 hashes a message to a point on the G2 curve using the SVDW map.
// It is faster than HashToG2, but the result is not uniformly distributed. Unsuitable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
//https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/#section-6.6.3
func EncodeToG2(msg, dst []byte) (G2Affine, error) {

	var res G2Affine
	u, err := hashToFp(msg, dst, 2)
	if err != nil {
		return res, err
	}

	res = mapToCurve2(&fptower.E2{
		A0: u[0],
		A1: u[1],
	})

	res.ClearCofactor(&res)
	return res, nil
}

// HashToG2 hashes a message to a point on the G2 curve using the SVDW map.
// Slower than EncodeToG2, but usable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#section-3
func HashToG2(msg, dst []byte) (G2Affine, error) {
	u, err := hashToFp(msg, dst, 2*2)
	if err != nil {
		return G2Affine{}, err
	}

	Q0 := mapToCurve2(&fptower.E2{
		A0: u[0],
		A1: u[1],
	})
	Q1 := mapToCurve2(&fptower.E2{
		A0: u[2+0],
		A1: u[2+1],
	})

	var _Q0, _Q1 G2Jac
	_Q0.FromAffine(&Q0)
	_Q1.FromAffine(&Q1).AddAssign(&_Q0)

	_Q1.ClearCofactor(&_Q1)

	Q1.FromJacobian(&_Q1)
	return Q1, nil
}

func g2NotZero(x *fptower.E2) uint64 {
	//Assuming G1 is over Fp and that if hashing is available for G2, it also is for G1
	return g1NotZero(&x.A0) | g1NotZero(&x.A1)

}
