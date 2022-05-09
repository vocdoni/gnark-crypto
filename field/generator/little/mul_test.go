package little

import (
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
	"math/big"
	"testing"
)

//go:noescape
func madd2Asm(res *[2]uint64, a, b, c, d uint64)

func TestMadd2(t *testing.T) {

	v := uint64(10)
	y1 := uint64(10)
	c1 := uint64(8)
	t1 := uint64(11)

	//x1, x0 := madd2(v, y1, c1, t1)

	var x [2]uint64
	madd2Asm(&x, v, y1, c1, t1)
}

func TestElementMulVerySpecial(t *testing.T) {
	a := Element{14831352495647895042, 10}
	var c Element
	_mulGeneric(&c, &a, &a)

	// t after round 0 =
	//0 = {uint64} 6431973904983132939
	//1 = {uint64} 11			CORRECT

	// after madd0 in round 1
	// c2 = 3219202104254117638
	// m = 3565083475063441699

	// after madd2 in round 1
	// c1, c0 = 0, 119

	// z after round 1 =
	//0 = {uint64} 14831352495647895042
	//1 = {uint64} 10			INCORRECT

	c.Mul(&a, &a)
}

func TestElementMulSpecial(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	specialValueTest := func() {
		// test special values against special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for i, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for j, b := range testValues {

				var bBig, d big.Int
				b.ToBigIntRegular(&bBig)

				var c Element
				c.Mul(&a, &b)
				d.Mul(&aBig, &bBig).Mod(&d, Modulus())

				// checking asm against generic impl
				var cGeneric Element
				_mulGeneric(&cGeneric, &a, &b)
				if !cGeneric.Equal(&c) {
					t.Log(a, b)
					t.Log("Got", c, "expected", cGeneric)
					t.Fatal("Mul failed special test values: asm and generic impl don't match\nFailed at (", i, ",", j, ")")
				}

				t.Log("Passed", a, b)

				/*if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("Mul failed special test values")
				}*/
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

	genA := gen()
	genB := gen()

	properties.Property("Mul: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			d.Set(&a.element)

			c.Mul(&a.element, &b.element)
			a.element.Mul(&a.element, &b.element)
			b.element.Mul(&d, &b.element)

			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Mul: operation result must match big.Int result", prop.ForAll(
		func(a, b testPairElement) bool {
			{
				var c Element

				c.Mul(&a.element, &b.element)

				var d, e big.Int
				d.Mul(&a.bigint, &b.bigint).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]Element, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int
				r.ToBigIntRegular(&rb)

				var c Element
				c.Mul(&a.element, &r)
				d.Mul(&a.bigint, &rb).Mod(&d, Modulus())

				// checking generic impl against asm path
				var cGeneric Element
				_mulGeneric(&cGeneric, &a.element, &r)
				if !cGeneric.Equal(&c) {
					// need to give context to failing error.
					return false
				}

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}
			return true
		},
		genA,
		genB,
	))

	properties.Property("Mul: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPairElement) bool {
			var c Element

			c.Mul(&a.element, &b.element)

			return !c.biggerOrEqualModulus()
		},
		genA,
		genB,
	))

	properties.Property("Mul: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			c.Mul(&a.element, &b.element)
			_mulGeneric(&d, &a.element, &b.element)
			return c.Equal(&d)
		},
		genA,
		genB,
	))

	// if we have ADX instruction enabled, test both path in assembly
	if supportAdx {
		t.Log("disabling ADX")
		supportAdx = false
		properties.TestingRun(t, gopter.ConsoleReporter(false))
		specialValueTest()
		supportAdx = true
	}
}
