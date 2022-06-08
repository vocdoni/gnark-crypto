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

package bls24315

import (
	"github.com/consensys/gnark-crypto/ecc/bls24-315/fp"
	"github.com/consensys/gnark-crypto/ecc/bls24-315/internal/fptower"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
	"math/rand"
	"strings"
	"testing"
)

func TestHashToFpG2(t *testing.T) {
	for _, c := range encodeToG2Vector.cases {
		elems, err := hashToFp([]byte(c.msg), encodeToG2Vector.dst, 4)
		if err != nil {
			t.Error(err)
		}
		g2TestMatchCoord(t, "u", c.msg, c.u, g2CoordAt(elems, 0))
	}

	for _, c := range hashToG2Vector.cases {
		elems, err := hashToFp([]byte(c.msg), hashToG2Vector.dst, 2*4)
		if err != nil {
			t.Error(err)
		}
		g2TestMatchCoord(t, "u0", c.msg, c.u0, g2CoordAt(elems, 0))
		g2TestMatchCoord(t, "u1", c.msg, c.u1, g2CoordAt(elems, 1))
	}
}

func TestMapToCurve2(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[G2] mapping output must be on curve", prop.ForAll(
		func(a fptower.E4) bool {

			g := mapToCurve2(&a)

			if !g.IsOnCurve() {
				t.Log("SVDW output not on curve")
				return false
			}

			return true
		},
		GenE4(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	for _, c := range encodeToG2Vector.cases {
		var u fptower.E4
		g2CoordSetString(&u, c.u)
		q := mapToCurve2(&u)
		g2TestMatchPoint(t, "Q", c.msg, c.Q, &q)
	}

	for _, c := range hashToG2Vector.cases {
		var u fptower.E4
		g2CoordSetString(&u, c.u0)
		q := mapToCurve2(&u)
		g2TestMatchPoint(t, "Q0", c.msg, c.Q0, &q)

		g2CoordSetString(&u, c.u1)
		q = mapToCurve2(&u)
		g2TestMatchPoint(t, "Q1", c.msg, c.Q1, &q)
	}
}

func TestMapToG2(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[G2] mapping to curve should output point on the curve", prop.ForAll(
		func(a fptower.E4) bool {
			g := MapToG2(a)
			return g.IsInSubGroup()
		},
		GenE4(),
	))

	properties.Property("[G2] mapping to curve should be deterministic", prop.ForAll(
		func(a fptower.E4) bool {
			g1 := MapToG2(a)
			g2 := MapToG2(a)
			return g1.Equal(&g2)
		},
		GenE4(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestEncodeToG2(t *testing.T) {
	t.Parallel()
	for _, c := range encodeToG2Vector.cases {
		p, err := EncodeToG2([]byte(c.msg), encodeToG2Vector.dst)
		if err != nil {
			t.Fatal(err)
		}
		g2TestMatchPoint(t, "P", c.msg, c.P, &p)
	}
}

func TestHashToG2(t *testing.T) {
	t.Parallel()
	for _, c := range hashToG2Vector.cases {
		p, err := HashToG2([]byte(c.msg), hashToG2Vector.dst)
		if err != nil {
			t.Fatal(err)
		}
		g2TestMatchPoint(t, "P", c.msg, c.P, &p)
	}
}

func BenchmarkEncodeToG2(b *testing.B) {
	const size = 54
	bytes := make([]byte, size)
	dst := encodeToG2Vector.dst
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		bytes[rand.Int()%size] = byte(rand.Int())

		if _, err := EncodeToG2(bytes, dst); err != nil {
			b.Fail()
		}
	}
}

func BenchmarkHashToG2(b *testing.B) {
	const size = 54
	bytes := make([]byte, size)
	dst := hashToG2Vector.dst
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		bytes[rand.Int()%size] = byte(rand.Int())

		if _, err := HashToG2(bytes, dst); err != nil {
			b.Fail()
		}
	}
}

//Only works on simple extensions (two-story towers)
func g2CoordSetString(z *fptower.E4, s string) {
	ssplit := strings.Split(s, ",")
	if len(ssplit) != 4 {
		panic("not equal to tower size")
	}
	z.SetString(
		ssplit[0],
		ssplit[1],
		ssplit[2],
		ssplit[3],
	)
}

func g2CoordAt(slice []fp.Element, i int) fptower.E4 {
	return fptower.E4{
		B0: slice[i*4+0],
		B1: slice[i*4+1],
		B2: slice[i*4+2],
		B3: slice[i*4+3],
	}
}

func g2TestMatchCoord(t *testing.T, coordName string, msg string, expectedStr string, seen fptower.E4) {
	var expected fptower.E4

	g2CoordSetString(&expected, expectedStr)

	if !expected.Equal(&seen) {
		t.Errorf("mismatch on \"%s\", %s:\n\texpected %s\n\tsaw      %s", msg, coordName, expected.String(), &seen)
	}
}

func g2TestMatchPoint(t *testing.T, pointName string, msg string, expected point, seen *G2Affine) {
	g2TestMatchCoord(t, pointName+".x", msg, expected.x, seen.X)
	g2TestMatchCoord(t, pointName+".y", msg, expected.y, seen.Y)
}

var encodeToG2Vector encodeTestVector
var hashToG2Vector hashTestVector
