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

// Package integration contains field arithmetic operations for modulus = 0xdd916e...25da05.
//
// The API is similar to math/big (big.Int), but the operations are significantly faster (up to 20x for the modular multiplication on amd64, see also https://hackmd.io/@gnark/modular_multiplication)
//
// The modulus is hardcoded in all the operations.
//
// Field elements are represented as an array, and assumed to be in Montgomery form in all methods:
//
//	type e_cios_0384 [6]uint64
//
// # Usage
//
// Example API signature:
//
//	// Mul z = x * y (mod q)
//	func (z *Element) Mul(x, y *Element) *Element
//
// and can be used like so:
//
//	var a, b Element
//	a.SetUint64(2)
//	b.SetString("984896738")
//	a.Mul(a, b)
//	a.Sub(a, a)
//	 .Add(a, b)
//	 .Inv(a)
//	b.Exp(b, new(big.Int).SetUint64(42))
//
// Modulus q =
//
//	q[base10] = 34102450673962811389969639094425662171713053853416332565622981543138899915666475256934880798495867643877776618281477
//	q[base16] = 0xdd916e8aa3a09686d24b97d018911a6ec60c72702b7002a792ca0222e7257f22b0005a0dc340da2919faebe56725da05
//
// # Warning
//
// This code has not been audited and is provided as-is. In particular, there is no security guarantees such as constant time implementation or side-channel attack resistance.
package integration
