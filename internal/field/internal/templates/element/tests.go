package element

const Test = `


import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"math/bits"
	"fmt"
	{{if .UsingP20Inverse}} 
	"github.com/consensys/gnark-crypto/internal/field"
	mrand "math/rand" 
	{{end}}
	"testing"
	
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
	ggen "github.com/leanovate/gopter/gen"

	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchRes{{.ElementName}} {{.ElementName}}

func Benchmark{{toTitle .ElementName}}Select(b *testing.B) {
	var x, y {{.ElementName}}
	x.SetRandom()
	y.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Select(i%3, &x, &y)
	}
}

func Benchmark{{toTitle .ElementName}}SetRandom(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = x.SetRandom()
	}
}

func Benchmark{{toTitle .ElementName}}SetBytes(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	bb := x.Bytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.SetBytes(bb[:])
	}

}

func Benchmark{{toTitle .ElementName}}MulByConstants(b *testing.B) {
	b.Run("mulBy3", func(b *testing.B){
		benchRes{{.ElementName}}.SetRandom()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MulBy3(&benchRes{{.ElementName}})
		}
	})
	b.Run("mulBy5", func(b *testing.B){
		benchRes{{.ElementName}}.SetRandom()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MulBy5(&benchRes{{.ElementName}})
		}
	})
	b.Run("mulBy13", func(b *testing.B){
		benchRes{{.ElementName}}.SetRandom()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MulBy13(&benchRes{{.ElementName}})
		}
	})
}

func Benchmark{{toTitle .ElementName}}Inverse(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Inverse(&x)
	}

}

func Benchmark{{toTitle .ElementName}}Butterfly(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Butterfly(&x, &benchRes{{.ElementName}})
	}
}


func Benchmark{{toTitle .ElementName}}Exp(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b1, _ := rand.Int(rand.Reader, Modulus())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Exp(x, b1)
	}
}


func Benchmark{{toTitle .ElementName}}Double(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Double(&benchRes{{.ElementName}})
	}
}


func Benchmark{{toTitle .ElementName}}Add(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Add(&x, &benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Sub(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sub(&x, &benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Neg(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Neg(&benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Div(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Div(&x, &benchRes{{.ElementName}})
	}
}


func Benchmark{{toTitle .ElementName}}FromMont(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.FromMont()
	}
}

func Benchmark{{toTitle .ElementName}}ToMont(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.ToMont()
	}
}
func Benchmark{{toTitle .ElementName}}Square(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Square(&benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Sqrt(b *testing.B) {
	var a {{.ElementName}}
	a.SetUint64(4)
	a.Neg(&a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sqrt(&a)
	}
}

func Benchmark{{toTitle .ElementName}}Mul(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Mul(&benchRes{{.ElementName}}, &x)
	}
}

func Benchmark{{toTitle .ElementName}}Cmp(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}} = x 
	benchRes{{.ElementName}}[0] = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Cmp(&x)
	}
}

func Test{{toTitle .ElementName}}Cmp(t *testing.T) {
	var x, y {{.ElementName}}
	
	if x.Cmp(&y) != 0 {
		t.Fatal("x == y")
	}

	one := One()
	y.Sub(&y, &one)

	if x.Cmp(&y) != -1 {
		t.Fatal("x < y")
	}
	if y.Cmp(&x) != 1 {
		t.Fatal("x < y")
	}

	x = y 
	if x.Cmp(&y) != 0 {
		t.Fatal("x == y")
	}

	x.Sub(&x, &one)
	if x.Cmp(&y) != -1 {
		t.Fatal("x < y")
	}
	if y.Cmp(&x) != 1 {
		t.Fatal("x < y")
	}
}

{{- if gt .NbWords 1}}
func Test{{toTitle .ElementName}}IsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y {{.ElementName}}
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

func Test{{toTitle .ElementName}}IsUint64(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)


	properties.Property("reduce should output a result smaller than modulus", prop.ForAll(
		func(v uint64) bool {
			var e {{.ElementName}}
			e.SetUint64(v)

			if !e.IsUint64() {
				return false
			}

			return e.Uint64() == v
		},
		ggen.UInt64(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

{{- end}}

func Test{{toTitle .ElementName}}NegZero(t *testing.T) {
	var a, b {{.ElementName}}
	b.SetZero()
	for a.IsZero() {
		a.SetRandom()
	}
	a.Neg(&b)
	if !a.IsZero() {
		t.Fatal("neg(0) != 0")
	}
}

// -------------------------------------------------------------------------------------------------
// Gopter tests
// most of them are generated with a template

{{ if gt .NbWords 6}}
const (
	nbFuzzShort = 20
	nbFuzz = 100
)
{{else}}
const (
	nbFuzzShort = 200
	nbFuzz = 1000
)
{{end}}

// special values to be used in tests
var staticTestValues []{{.ElementName}}

func init() {
	staticTestValues = append(staticTestValues, {{.ElementName}}{}) // zero
	staticTestValues = append(staticTestValues, One()) 				// one
	staticTestValues = append(staticTestValues, rSquare) 			// r²
	var e, one {{.ElementName}}
	one.SetOne()
	e.Sub(&q{{.ElementName}}, &one)
	staticTestValues = append(staticTestValues, e) 	// q - 1
	e.Double(&one)
	staticTestValues = append(staticTestValues, e) 	// 2 


	{
		a := q{{.ElementName}}
		a[0]--
		staticTestValues = append(staticTestValues, a)
	}

	{{- $qi := index $.Q $.NbWordsLastIndex}}
	{{- range $i := iterate 0 3}}
		staticTestValues = append(staticTestValues, {{$.ElementName}}{ {{$i}} })
		{{- if gt $.NbWords 1}}
			{{- if le $i $qi}}
			staticTestValues = append(staticTestValues, {{$.ElementName}}{0, {{$i}} })
			{{- end}}
		{{- end}}
	{{- end}}

	{
		a := q{{.ElementName}}
		a[{{.NbWordsLastIndex}}]--
		staticTestValues = append(staticTestValues, a)
	}

	{{- if ne .NbWords 1}}
	{
		a := q{{.ElementName}}
		a[{{.NbWordsLastIndex}}]--
		a[0]++
		staticTestValues = append(staticTestValues, a)
	}
	{{- end}}

	{
		a := q{{.ElementName}}
		a[{{.NbWordsLastIndex}}] = 0
		staticTestValues = append(staticTestValues, a)
	}

	for _, sv := range staticTestValues {
		if !sv.smallerThanModulus() {
			panic("invalid static test value") // shouldn't happen.
		}
	}

}

func Test{{toTitle .ElementName}}Reduce(t *testing.T) {
	testValues := make([]{{.ElementName}}, len(staticTestValues))
	copy(testValues, staticTestValues)

	for _, s := range testValues {
		expected := s
		reduce(&s)
		_reduceGeneric(&expected)
		if !s.Equal(&expected) {
			t.Fatal("reduce failed: asm and generic impl don't match")
		}
	}


	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := genFull()

	properties.Property("reduce should output a result smaller than modulus", prop.ForAll(
		func(a {{.ElementName}}) bool {
			b := a
			reduce(&a)
			_reduceGeneric(&b)
			return a.smallerThanModulus()  && a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	
}

func Test{{toTitle .ElementName}}Equal(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("x.Equal(&y) iff x == y; likely false for random pairs", prop.ForAll(
		func(a testPair{{.ElementName}}, b testPair{{.ElementName}}) bool {
			return a.element.Equal(&b.element) == (a.element == b.element)
		},
		genA,
		genB,
	))

	properties.Property("x.Equal(&y) if x == y", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			b := a.element
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func Test{{toTitle .ElementName}}Bytes(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("SetBytes(Bytes()) should stay constant", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var b {{.ElementName}}
			bytes := a.element.Bytes()
			b.SetBytes(bytes[:])
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func Test{{toTitle .ElementName}}InverseExp(t *testing.T) {
	// inverse must be equal to exp^-2
	exp := Modulus()
	exp.Sub(exp, new(big.Int).SetUint64(2))

	invMatchExp := func(a testPair{{.ElementName}}) bool {
		var b {{.ElementName}}
		b.Set(&a.element)
		a.element.Inverse(&a.element)
		b.Exp(b, exp)

		return a.element.Equal(&b)
	}

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}
	properties := gopter.NewProperties(parameters)
	genA := gen()
	properties.Property("inv == exp^-2", prop.ForAll(invMatchExp, genA))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

	parameters.MinSuccessfulTests = 1
	properties = gopter.NewProperties(parameters)
	properties.Property("inv(0) == 0", prop.ForAll(invMatchExp, ggen.OneConstOf(testPair{{.ElementName}}{})))
	properties.TestingRun(t, gopter.ConsoleReporter(false))


}


func mulByConstant(z *{{.ElementName}}, c uint8) {
	var y {{.ElementName}}
	y.SetUint64(uint64(c))
	z.Mul(z, &y)
}


func Test{{toTitle .ElementName}}MulByConstants(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	implemented := []uint8{0,1,2,3,5,13}
	properties.Property("mulByConstant", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			for _, c := range implemented {
				var constant {{.ElementName}}
				constant.SetUint64(uint64(c))
	
				b := a.element 
				b.Mul(&b, &constant)

				aa := a.element
				mulByConstant(&aa, c)

				if !aa.Equal(&b) {
					return false
				}
			}
			
			return true
		},
		genA,
	))
	

	properties.Property("MulBy3(x) == Mul(x, 3)", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var constant {{.ElementName}}
			constant.SetUint64(3)

			b := a.element 
			b.Mul(&b, &constant)

			MulBy3(&a.element)

			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("MulBy5(x) == Mul(x, 5)", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var constant {{.ElementName}}
			constant.SetUint64(5)

			b := a.element 
			b.Mul(&b, &constant)

			MulBy5(&a.element)

			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("MulBy13(x) == Mul(x, 13)", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var constant {{.ElementName}}
			constant.SetUint64(13)

			b := a.element 
			b.Mul(&b, &constant)

			MulBy13(&a.element)

			return a.element.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	
}

func Test{{toTitle .ElementName}}Legendre(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("legendre should output same result than big.Int.Jacobi", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			return a.element.Legendre() == big.Jacobi(&a.bigint, Modulus()) 
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	
}

func Test{{toTitle .ElementName}}BitLen(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("BitLen should output same result than big.Int.BitLen", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			return a.element.FromMont().BitLen() ==  a.bigint.BitLen()
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	
}



func Test{{toTitle .ElementName}}Butterflies(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("butterfly0 == a -b; a +b", prop.ForAll(
		func(a,b testPair{{.ElementName}}) bool {
			a0, b0 := a.element, b.element 
			
			_butterflyGeneric(&a.element, &b.element)
			Butterfly(&a0, &b0)

			return a.element.Equal(&a0) && b.element.Equal(&b0)
		},
		genA,
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))


}

func Test{{toTitle .ElementName}}LexicographicallyLargest(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("element.Cmp should match LexicographicallyLargest output", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var negA {{.ElementName}}
			negA.Neg(&a.element)

			cmpResult := a.element.Cmp(&negA)
			lResult := a.element.LexicographicallyLargest()

			if lResult && cmpResult == 1 {
				return true 
			}
			if !lResult && cmpResult !=1 {
				return true
			}
			return false
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	
}


{{template "testBinaryOp" dict "all" . "Op" "Add"}}
{{template "testBinaryOp" dict "all" . "Op" "Sub"}}
{{template "testBinaryOp" dict "all" . "Op" "Mul" "GenericOp" "_mulGeneric"}}
{{template "testBinaryOp" dict "all" . "Op" "Div"}}
{{template "testBinaryOp" dict "all" . "Op" "Exp"}}

{{template "testUnaryOp" dict "all" . "Op" "Square" }}
{{template "testUnaryOp" dict "all" . "Op" "Inverse"}}
{{template "testUnaryOp" dict "all" . "Op" "Sqrt"}}
{{template "testUnaryOp" dict "all" . "Op" "Double"}}
{{template "testUnaryOp" dict "all" . "Op" "Neg" }}

{{ define "testBinaryOp" }}

func Test{{toTitle .all.ElementName}}{{.Op}}(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}
	

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("{{.Op}}: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c, d {{.all.ElementName}}
			d.Set(&a.element)
			{{if eq .Op "Exp"}}
				c.{{.Op}}(a.element, &b.bigint)
				a.element.{{.Op}}(a.element, &b.bigint)
				b.element.{{.Op}}(d, &b.bigint)
			{{else}}
				c.{{.Op}}(&a.element, &b.element)
				a.element.{{.Op}}(&a.element, &b.element)
				b.element.{{.Op}}(&d, &b.element)
			{{end}}
			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("{{.Op}}: operation result must match big.Int result", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			{
				var c {{.all.ElementName}}
				{{if eq .Op "Exp"}}
					c.{{.Op}}(a.element, &b.bigint)
				{{else}}
					c.{{.Op}}(&a.element, &b.element)
				{{end}}
				var d, e big.Int 
				{{- if eq .Op "Div"}}
					d.ModInverse(&b.bigint, Modulus())
					d.Mul(&d, &a.bigint).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					d.Exp(&a.bigint, &b.bigint, Modulus())
				{{- else}}
					d.{{.Op}}(&a.bigint, &b.bigint).Mod(&d, Modulus())
				{{- end }}


				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				} 
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]{{.all.ElementName}}, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int 
				r.ToBigIntRegular(&rb) 

				var c {{.all.ElementName}}
				{{- if eq .Op "Div"}}
					c.{{.Op}}(&a.element, &r)
					d.ModInverse(&rb, Modulus())
					d.Mul(&d, &a.bigint).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					c.{{.Op}}(a.element, &rb)
					d.Exp(&a.bigint, &rb, Modulus())
				{{- else}}
					c.{{.Op}}(&a.element, &r)
					d.{{.Op}}(&a.bigint, &rb).Mod(&d, Modulus())
				{{- end }}

				{{if .GenericOp}}
					// checking generic impl against asm path
					var cGeneric {{.all.ElementName}}
					cGeneric.{{.GenericOp}}( &a.element, &r)
					if !cGeneric.Equal(&c) {
						// need to give context to failing error.
						return false
					}
				{{end}}

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				} 
			}
			return true 
		},
		genA,
		genB,
	))

	properties.Property("{{.Op}}: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			{{if eq .Op "Exp"}}
				c.{{.Op}}(a.element, &b.bigint)
			{{else}}
				c.{{.Op}}(&a.element, &b.element)
			{{end}}
			return c.smallerThanModulus()
		},
		genA,
		genB,
	))

	{{if .GenericOp}}
	properties.Property("{{.Op}}: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c,d {{.all.ElementName}}
			c.{{.Op}}(&a.element, &b.element)
			d.{{.GenericOp}}( &a.element, &b.element)
			return c.Equal(&d)
		},
		genA,
		genB,
	))

	{{end}}


	specialValueTest := func() {
		// test special values against special values
		testValues := make([]{{.all.ElementName}}, len(staticTestValues))
		copy(testValues, staticTestValues)
	
		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for _, b := range testValues {

				var bBig, d, e big.Int 
				b.ToBigIntRegular(&bBig)

				var c {{.all.ElementName}}
				


				{{- if eq .Op "Div"}}
					c.{{.Op}}(&a, &b)
					d.ModInverse(&bBig, Modulus())
					d.Mul(&d, &aBig).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					c.{{.Op}}(a, &bBig)
					d.Exp(&aBig, &bBig, Modulus())
				{{- else}}
					c.{{.Op}}(&a, &b)
					d.{{.Op}}(&aBig, &bBig).Mod(&d, Modulus())
				{{- end }}
	
				{{if .GenericOp}}
					// checking asm against generic impl
					var cGeneric {{.all.ElementName}}
					cGeneric.{{.GenericOp}}( &a, &b)
					if !cGeneric.Equal(&c) {
						t.Fatal("{{.Op}} failed special test values: asm and generic impl don't match")
					}
				{{end}}
				

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("{{.Op}} failed special test values")
				} 
			}
		}
	}


	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

{{ end }}


{{ define "testUnaryOp" }}

func Test{{toTitle .all.ElementName}}{{.Op}}(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("{{.Op}}: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			{{if eq .Op "Sqrt"}}
			b := a.element
			{{else}}
			var b {{.all.ElementName}}
			{{end}}
			b.{{.Op}}(&a.element)
			a.element.{{.Op}}(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("{{.Op}}: operation result must match big.Int result", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			c.{{.Op}}(&a.element)

			var d, e big.Int 
			{{- if eq .Op "Square"}}
				d.Mul(&a.bigint, &a.bigint).Mod(&d, Modulus())
			{{- else if eq .Op "Inverse"}}
				d.ModInverse(&a.bigint, Modulus())
			{{- else if eq .Op "Sqrt"}}
				d.ModSqrt(&a.bigint, Modulus())
			{{- else if eq .Op "Double"}}
				d.Lsh(&a.bigint, 1).Mod(&d, Modulus())
			{{- else if eq .Op "Neg"}}
				d.Neg(&a.bigint).Mod(&d, Modulus())
			{{- end }}


			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("{{.Op}}: operation result must be smaller than modulus", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			c.{{.Op}}(&a.element)
			return c.smallerThanModulus()
		},
		genA,
	))

	{{if .GenericOp}}
	properties.Property("{{.Op}}: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c,d {{.all.ElementName}}
			c.{{.Op}}(&a.element)
			d.{{.GenericOp}}( &a.element)
			return c.Equal(&d)
		},
		genA,
	))

	{{end}}


	specialValueTest := func() {
		// test special values
		testValues := make([]{{.all.ElementName}}, len(staticTestValues))
		copy(testValues, staticTestValues)
	
		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c {{.all.ElementName}}
			c.{{.Op}}(&a)

			var  d, e big.Int 
			{{- if eq .Op "Square"}}
				d.Mul(&aBig, &aBig).Mod(&d, Modulus())
			{{- else if eq .Op "Inverse"}}
				d.ModInverse(&aBig, Modulus())
			{{- else if eq .Op "Sqrt"}}
				d.ModSqrt(&aBig, Modulus())
			{{- else if eq .Op "Double"}}
				d.Lsh(&aBig, 1).Mod(&d, Modulus())
			{{- else if eq .Op "Neg"}}
				d.Neg(&aBig).Mod(&d, Modulus())
			{{- end }}

			{{if .GenericOp}}
				// checking asm against generic impl
				var cGeneric {{.all.ElementName}}
				cGeneric.{{.GenericOp}}( &a)
				if !cGeneric.Equal(&c) {
					t.Fatal("{{.Op}} failed special test values: asm and generic impl don't match")
				}
			{{end}}
			

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("{{.Op}} failed special test values")
			} 
		}
	}


	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

{{ end }}

{{ if .UseAddChain}}
func Test{{toTitle .ElementName}}FixedExp(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	var (
		_bLegendreExponent{{.ElementName}} *big.Int
		_bSqrtExponent{{.ElementName}} *big.Int
	)

	_bLegendreExponent{{.ElementName}}, _ = new(big.Int).SetString("{{.LegendreExponent}}", 16)
	{{- if .SqrtQ3Mod4}}
		const sqrtExponent{{.ElementName}} = "{{.SqrtQ3Mod4Exponent}}"
	{{- else if .SqrtAtkin}}
		const sqrtExponent{{.ElementName}} = "{{.SqrtAtkinExponent}}"
	{{- else if .SqrtTonelliShanks}}
		const sqrtExponent{{.ElementName}} = "{{.SqrtSMinusOneOver2}}"
	{{- end }}
	_bSqrtExponent{{.ElementName}}, _ = new(big.Int).SetString(sqrtExponent{{.ElementName}}, 16)

	genA := gen()

	properties.Property(fmt.Sprintf("expBySqrtExp must match Exp(%s)", sqrtExponent{{.ElementName}}), prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			d := a.element
			c.expBySqrtExp(c)
			d.Exp(d, _bSqrtExponent{{.ElementName}})
			return c.Equal(&d)
		},
		genA,
	))

	properties.Property("expByLegendreExp must match Exp({{.LegendreExponent}})", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			d := a.element
			c.expByLegendreExp(c)
			d.Exp(d, _bLegendreExponent{{.ElementName}})
			return c.Equal(&d)
		},
		genA,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

{{ end }}





func Test{{toTitle .ElementName}}Halve(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	var twoInv {{.ElementName}}
	twoInv.SetUint64(2)
	twoInv.Inverse(&twoInv)

	properties.Property("z.Halve must match z / 2", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			d := a.element
			c.Halve()
			d.Mul(&d, &twoInv)
			return c.Equal(&d)
		},
		genA,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func combineSelectionArguments(c int64, z int8) int {
	if z%3 == 0 {
		return 0
	}
	return int(c)
}

func Test{{toTitle .ElementName}}Select(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := genFull()
	genB := genFull()
	genC := ggen.Int64()	//the condition
	genZ := ggen.Int8()	//to make zeros artificially more likely

	properties.Property("Select: must select correctly", prop.ForAll(
		func(a, b {{.ElementName}}, cond int64, z int8) bool {
			condC := combineSelectionArguments(cond, z)

			var c {{.ElementName}}
			c.Select(condC, &a, &b)
			
			if condC == 0 {
				return c.Equal(&a)
			}
			return c.Equal(&b)
		},
		genA,
		genB,
		genC,
		genZ,
	))

	properties.Property("Select: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b {{.ElementName}}, cond int64, z int8) bool {
			condC := combineSelectionArguments(cond, z)
			
			var c, d {{.ElementName}}
			d.Set(&a)
			c.Select(condC, &a, &b)
			a.Select(condC, &a, &b)
			b.Select(condC, &d, &b)
			return a.Equal(&b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
		genC,
		genZ,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func Test{{toTitle .ElementName}}SetInt64(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("z.SetInt64 must match z.SetString", prop.ForAll(
		func(a testPair{{.ElementName}}, v int64) bool {
			c := a.element
			d := a.element

			c.SetInt64(v)
			d.SetString(fmt.Sprintf("%v",v))

			return c.Equal(&d)
		},
		genA, ggen.Int64(),
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}


func Test{{toTitle .ElementName}}SetInterface(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genInt := ggen.Int
	genInt8 := ggen.Int8
	genInt16 := ggen.Int16
	genInt32 := ggen.Int32
	genInt64 := ggen.Int64

	genUint := ggen.UInt
	genUint8 := ggen.UInt8
	genUint16 := ggen.UInt16
	genUint32 := ggen.UInt32
	genUint64 := ggen.UInt64

	{{setInterface .ElementName "int8"}}
	{{setInterface .ElementName "int16"}}
	{{setInterface .ElementName "int32"}}
	{{setInterface .ElementName "int64"}}
	{{setInterface .ElementName "int"}}

	{{setInterface .ElementName "uint8"}}
	{{setInterface .ElementName "uint16"}}
	{{setInterface .ElementName "uint32"}}
	{{setInterface .ElementName "uint64"}}
	{{setInterface .ElementName "uint"}}


	properties.TestingRun(t, gopter.ConsoleReporter(false))

	{
		assert := require.New(t)
		var e {{.ElementName}}
		r, err := e.SetInterface(nil)
		assert.Nil(r)
		assert.Error(err)

		var ptE *{{.ElementName}}
		var ptB *big.Int

		r, err = e.SetInterface(ptE)
		assert.Nil(r)
		assert.Error(err)
		ptE = new({{.ElementName}}).SetOne()
		r, err = e.SetInterface(ptE)
		assert.NoError(err)
		assert.True(r.IsOne())

		r, err = e.SetInterface(ptB)
		assert.Nil(r)
		assert.Error(err)

	}
}


{{define "setInterface eName tName"}}

properties.Property("z.SetInterface must match z.SetString with {{.tName}}", prop.ForAll(
	func(a testPair{{.eName}}, v {{.tName}}) bool {
		c := a.element
		d := a.element

		c.SetInterface(v)
		d.SetString(fmt.Sprintf("%v",v))

		return c.Equal(&d)
	},
	genA, gen{{toTitle .tName}}(),
))

{{end}}

func Test{{toTitle .ElementName}}NegativeExp(t *testing.T) {
	t.Parallel()

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	
	properties.Property("x⁻ᵏ == 1/xᵏ", prop.ForAll(
		func(a,b testPair{{.ElementName}}) bool {

			var nb, d, e big.Int 
			nb.Neg(&b.bigint)

			var c {{.ElementName}}
			c.Exp(a.element, &nb)

			d.Exp(&a.bigint, &nb, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA, genA,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func Test{{toTitle .ElementName}}New{{.ElementName}}(t *testing.T) {
	assert := require.New(t)

	t.Parallel()

	e := New{{.ElementName}}(1)
	assert.True(e.IsOne())

	e = New{{.ElementName}}(0)
	assert.True(e.IsZero())
}


func Test{{toTitle .ElementName}}BatchInvert(t *testing.T) {
	assert := require.New(t)

	t.Parallel()

	// ensure batchInvert([x]) == invert(x)
	for i:=int64(-1); i <=2; i++ {
		var e, eInv {{.ElementName}}
		e.SetInt64(i)
		eInv.Inverse(&e)

		a := []{{.ElementName}}{e}
		aInv := BatchInvert(a)

		assert.True(aInv[0].Equal(&eInv), "batchInvert != invert")

	}

	// test x * x⁻¹ == 1
	tData := [][]int64 {
		[]int64{-1,1,2,3},
		[]int64{0, -1,1,2,3, 0},
		[]int64{0, -1,1,0, 2,3, 0},
		[]int64{-1,1,0, 2,3},
		[]int64{0,0,1},
		[]int64{1,0,0},
		[]int64{0,0,0},
	}

	for _, t := range tData {
		a := make([]{{.ElementName}}, len(t))
		for i:=0; i <len(a);i++ {
			a[i].SetInt64(t[i])
		}

		aInv := BatchInvert(a)

		assert.True(len(aInv) == len(a))

		for i:=0; i <len(a);i++ {
			if a[i].IsZero() {
				assert.True(aInv[i].IsZero(), "0⁻¹ != 0")
			} else {
				assert.True(a[i].Mul(&a[i], &aInv[i]).IsOne(), "x * x⁻¹ != 1")
			}
		}
	}


	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("batchInvert --> x * x⁻¹ == 1", prop.ForAll(
		func(tp testPair{{.ElementName}}, r uint8) bool {

			a := make([]{{.ElementName}}, r)
			if r != 0 {
				a[0] = tp.element

			}
			one := One()
			for i:=1; i <len(a);i++ {
				a[i].Add(&a[i-1], &one)
			}
	
			aInv := BatchInvert(a)
	
			assert.True(len(aInv) == len(a))
	
			for i:=0; i <len(a);i++ {
				if a[i].IsZero() {
					if !aInv[i].IsZero() {
						return false 
					}
				} else {
					if !a[i].Mul(&a[i], &aInv[i]).IsOne() {
						return false
					}
				}
			}
			return true
		},
		genA,ggen.UInt8(),
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func Test{{toTitle .ElementName}}FromMont(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Assembly implementation must be consistent with generic one", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			d := a.element
			c.FromMont()
			_fromMontGeneric(&d)
			return c.Equal(&d)
		},
		genA,
	))

	properties.Property("x.FromMont().ToMont() == x", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			c.FromMont().ToMont()
			return c.Equal(&a.element)
		},
		genA,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}



func Test{{toTitle .ElementName}}JSON(t *testing.T) {
	assert := require.New(t)

	type S struct {
		A {{.ElementName}}
		B [3]{{.ElementName}}
		C *{{.ElementName}}
		D *{{.ElementName}}
	}

	// encode to JSON
	var s S
	s.A.SetString("-1")
	s.B[2].SetUint64(42)
	s.D = new({{.ElementName}}).SetUint64(8000)

	encoded, err := json.Marshal(&s)
	assert.NoError(err)
	{{- if eq $.NbWords 1}}
	// since our modulus is on 1 word, we may need to adjust "42" and "8000" values;
	formatValue := func(v int64) string {
		const maxUint16 = 65535
		var a, aNeg big.Int 
		a.SetInt64(v)
		a.Mod(&a, Modulus())
		aNeg.Neg(&a).Mod(&aNeg, Modulus())
		fmt.Println("aNeg", aNeg.Text(10))
		if aNeg.Uint64() != 0 && aNeg.Uint64() <= maxUint16 {
			return "-"+aNeg.Text(10)
		}
		return a.Text(10)
	} 
	expected := fmt.Sprintf("{\"A\":-1,\"B\":[0,0,%s],\"C\":null,\"D\":%s}", formatValue(42), formatValue(8000))
	{{- else}}
	const expected = "{\"A\":-1,\"B\":[0,0,42],\"C\":null,\"D\":8000}"
	{{- end}}
	assert.Equal(expected, string(encoded))

	// decode valid
	var decoded S
	err = json.Unmarshal([]byte(expected), &decoded)
	assert.NoError(err)

	assert.Equal(s, decoded, "element -> json -> element round trip failed")

	// decode hex and string values
	withHexValues := "{\"A\":\"-1\",\"B\":[0,\"0x00000\",\"0x2A\"],\"C\":null,\"D\":\"8000\"}"

	var decodedS S
	err = json.Unmarshal([]byte(withHexValues), &decodedS)
	assert.NoError(err)

	assert.Equal(s, decodedS, " json with strings  -> element  failed")

}

type testPair{{.ElementName}} struct {
	element {{.ElementName}}
	bigint       big.Int
}


func gen() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var g testPair{{.ElementName}}

		g.element = {{.ElementName}}{
			{{- range $i := .NbWordsIndexesFull}}
			genParams.NextUint64(),{{end}}
		}
		if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
			g.element[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
		}
		

		for !g.element.smallerThanModulus() {
			g.element = {{.ElementName}}{
				{{- range $i := .NbWordsIndexesFull}}
				genParams.NextUint64(),{{end}}
			}
			if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
				g.element[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
			}
		}

		g.element.ToBigIntRegular(&g.bigint)
		genResult := gopter.NewGenResult(g, gopter.NoShrinker)
		return genResult
	}
}


func genFull() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {

		genRandomFq := func() {{.ElementName}} {
			var g {{.ElementName}}

			g = {{.ElementName}}{
				{{- range $i := .NbWordsIndexesFull}}
				genParams.NextUint64(),{{end}}
			}

			if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
				g[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
			}

			for !g.smallerThanModulus() {
				g = {{.ElementName}}{
					{{- range $i := .NbWordsIndexesFull}}
					genParams.NextUint64(),{{end}}
				}
				if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
					g[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
				}
			}

			return g 
		}
		a := genRandomFq()

		var carry uint64
		{{- range $i := .NbWordsIndexesFull}}
			{{- if eq $i $.NbWordsLastIndex}}
			a[{{$i}}], _ = bits.Add64(a[{{$i}}], q{{$.ElementName}}[{{$i}}], carry)
			{{- else}}
			a[{{$i}}], carry = bits.Add64(a[{{$i}}], q{{$.ElementName}}[{{$i}}], carry)
			{{- end}}
		{{- end}}
		
		genResult := gopter.NewGenResult(a, gopter.NoShrinker)
		return genResult
	}
}
{{if $.UsingP20Inverse}}
func (z *{{.ElementName}}) matchVeryBigInt(aHi uint64, aInt *big.Int) error {
	var modulus big.Int
	var aIntMod big.Int
	modulus.SetInt64(1)
	modulus.Lsh(&modulus, (Limbs+1)*64)
	aIntMod.Mod(aInt, &modulus)

	slice := append(z[:], aHi)

	return field.BigIntMatchUint64Slice(&aIntMod, slice)
}

//TODO: Phase out in favor of property based testing
func (z *{{.ElementName}}) assertMatchVeryBigInt(t *testing.T, aHi uint64, aInt *big.Int) {

	if err := z.matchVeryBigInt(aHi, aInt); err != nil {
		t.Error(err)
	}
}
{{- end}}



// TEMPORARY SECTION FOR BENCHMARKS -- 


{{ define "SQUARE_NOCARRY" }}
var {{range $i := .all.NbWordsIndexesFull}}t{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
var {{range $i := $.all.NbWordsIndexesFull}}u{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
var {{range $i := interval 0 (add $.all.NbWordsLastIndex 1)}}lo{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64

// note that if hi, _ = bits.Mul64() didn't generate
// UMULH and MUL, (but just UMULH) we could use same pattern
// as in mulRaw and reduce the stack space of this function (no need for lo..)

{{- range $i := .all.NbWordsIndexesFull}}
{

	{{$jStart := add $i 1}}
	{{$jEnd := add $.all.NbWordsLastIndex 1}}

	var c0, c2 uint64


	// for j=i+1 to N-1
	//     p,C,t[j] = 2*a[j]*a[i] + t[j] + (p,C)
	// A = C

	{{- if eq $i 0}}
		u{{$i}}, lo1 = bits.Mul64(x[{{$i}}], x[{{$i}}])
		{{- range $j := interval $jStart $jEnd}}
			u{{$j}}, t{{$j}} = bits.Mul64(x[{{$j}}], x[{{$i}}])
		{{- end}}

		// propagate lo, from t[j] to end, twice.
		{{- range $j := interval $jStart $jEnd}}
			{{- if eq $j $jStart}}
				t{{$j}}, c0 = bits.Add64(t{{$j}}, t{{$j}}, 0)
			{{- else }}
				t{{$j}}, c0 = bits.Add64(t{{$j}}, t{{$j}}, c0)
			{{- end}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				c2, _ = bits.Add64(c2, 0, c0)
			{{- end}}
		{{- end}}

		t{{$i}}, c0 = bits.Add64( lo1,t{{$i}}, 0)
	{{- else}}
		{{- range $j := interval (sub $jStart 1) $jEnd}}
			u{{$j}}, lo{{$j}} = bits.Mul64(x[{{$j}}], x[{{$i}}])
		{{- end}}

		// propagate lo, from t[j] to end, twice.
		{{- range $j := interval $jStart $jEnd}}
			{{- if eq $j $jStart}}
				lo{{$j}}, c0 = bits.Add64(lo{{$j}}, lo{{$j}}, 0)
			{{- else }}
				lo{{$j}}, c0 = bits.Add64(lo{{$j}}, lo{{$j}}, c0)
			{{- end}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				c2, _ = bits.Add64(c2, 0, c0)
			{{- end}}
		{{- end}}
		{{- range $j := interval $jStart $jEnd}}
			{{- if eq $j $jStart}}
				t{{$j}}, c0 = bits.Add64(lo{{$j}}, t{{$j}}, 0)
			{{- else }}
				t{{$j}}, c0 = bits.Add64(lo{{$j}}, t{{$j}}, c0)
			{{- end}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				c2, _ = bits.Add64(c2, 0, c0)
			{{- end}}
		{{- end}}

		t{{$i}}, c0 = bits.Add64( lo{{$i}},t{{$i}}, 0)
	{{- end}}


	// propagate u{{$i}} + hi
	{{- range $j := interval $jStart $jEnd}}
		t{{$j}}, c0 = bits.Add64(u{{sub $j 1}}, t{{$j}}, c0)
	{{- end}}
	c2, _ = bits.Add64(u{{$.all.NbWordsLastIndex}}, c2, c0)

	// hi again
	{{- range $j := interval $jStart $jEnd}}
		{{- if eq $j $.all.NbWordsLastIndex}}
		c2, _ = bits.Add64(c2, u{{$j}}, {{- if eq $j $jStart}} 0 {{- else}}c0{{- end}})
		{{- else if eq $j $jStart}}
			t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, 0)
		{{- else }}
			t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, c0)
		{{- end}}
	{{- end}}

	{{- $k := $.all.NbWordsLastIndex}}

	// this part is unchanged.
	m := qInvNeg * t0
	{{- range $j := $.all.NbWordsIndexesFull}}
		u{{$j}}, lo{{$j}} = bits.Mul64(m, q{{$j}})
	{{- end}}
	{{- range $j := $.all.NbWordsIndexesFull}}
	{{- if ne $j 0}}
		{{- if eq $j 1}}
			_, c0 = bits.Add64(t0, lo{{sub $j 1}}, 0)
		{{- else}}
			t{{sub $j 2}}, c0 = bits.Add64(t{{sub $j 1}}, lo{{sub $j 1}}, c0)
		{{- end}}
	{{- end}}
	{{- end}}
	t{{sub $.all.NbWordsLastIndex 1}}, c0 = bits.Add64(0, lo{{$.all.NbWordsLastIndex}}, c0) 
	u{{$k}}, _ = bits.Add64(u{{$k}}, 0, c0)

	{{- range $j := $.all.NbWordsIndexesFull}}
		{{- if eq $j 0}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, 0)
		{{- else if eq $j $.all.NbWordsLastIndex}}
			c2, _ = bits.Add64(c2, 0, c0)
		{{- else}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, c0)
		{{- end}}
	{{- end}}

	{{- $k := sub $.all.NbWordsLastIndex 0}}
	{{- $l := sub $.all.NbWordsLastIndex 1}}
	t{{$l}}, c0 = bits.Add64(t{{$k}}, t{{$l}}, 0)
	t{{$k}}, _ = bits.Add64(u{{$k}}, c2, c0)
}
{{- end}}


{{- range $i := $.all.NbWordsIndexesFull}}
z[{{$i}}] = t{{$i}}
{{- end}}

{{ end }}


{{ define "MUL_NOCARRY" }}
var t [{{.all.NbWords}}]uint64
var c [3]uint64
{{- range $j := .all.NbWordsIndexesFull}}
{
	// round {{$j}}
	v := {{$.V1}}[{{$j}}]
	{{- if eq $j $.all.NbWordsLastIndex}}
		c[1], c[0] = madd1(v, {{$.V2}}[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		{{- if eq $.all.NbWords 1}}
			z[0], _ = madd3(m, q0, c[0], c[2], c[1])
		{{- else}}
			{{- range $i := $.all.NbWordsIndexesNoZero}}
				c[1], c[0] = madd2(v, {{$.V2}}[{{$i}}],  c[1], t[{{$i}}])
				{{- if eq $i $.all.NbWordsLastIndex}}
					z[{{sub $.all.NbWords 1}}], z[{{sub $i 1}}] = madd3(m, q{{$i}}, c[0], c[2], c[1])
				{{- else}}
					c[2], z[{{sub $i 1}}] = madd2(m, q{{$i}},  c[2], c[0])
				{{- end}}
			{{- end}}
		{{- end}}
	{{- else if eq $j 0}}
		c[1], c[0] = bits.Mul64(v, {{$.V2}}[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		{{- range $i := $.all.NbWordsIndexesNoZero}}
			c[1], c[0] = madd1(v, {{$.V2}}[{{$i}}], c[1])
			{{- if eq $i $.all.NbWordsLastIndex}}
				t[{{sub $.all.NbWords 1}}], t[{{sub $i 1}}]  = madd3(m, q{{$i}}, c[0], c[2], c[1])
			{{- else}}
				c[2], t[{{sub $i 1}}] = madd2(m, q{{$i}}, c[2], c[0])
			{{- end}}
		{{- end}}
	{{- else}}
		c[1], c[0] = madd1(v, {{$.V2}}[0], t[0])
		m := c[0] * qInvNeg
		c[2] = madd0(m, q0, c[0])
		{{- range $i := $.all.NbWordsIndexesNoZero}}
			c[1], c[0] = madd2(v, {{$.V2}}[{{$i}}], c[1], t[{{$i}}])
			{{- if eq $i $.all.NbWordsLastIndex}}
				t[{{sub $.all.NbWords 1}}], t[{{sub $i 1}}] = madd3(m, q{{$i}}, c[0], c[2], c[1])
			{{- else}}
				c[2], t[{{sub $i 1}}] = madd2(m, q{{$i}}, c[2], c[0])
			{{- end}}
		{{- end}}
	{{-  end }}
}
{{- end}}
{{ end }}


{{ define "MUL_CIOS" }}
	var t [{{add .all.NbWords 1}}]uint64
	var D uint64
	var m, C uint64

	{{- range $j := .all.NbWordsIndexesFull}}
		// -----------------------------------
		// First loop
		{{ if eq $j 0}}
			C, t[0] = bits.Mul64({{$.V2}}[{{$j}}], {{$.V1}}[0])
			{{- range $i := $.all.NbWordsIndexesNoZero}}
				C, t[{{$i}}] = madd1({{$.V2}}[{{$j}}], {{$.V1}}[{{$i}}], C)
			{{- end}}
		{{ else }}
			C, t[0] = madd1({{$.V2}}[{{$j}}], {{$.V1}}[0], t[0])
			{{- range $i := $.all.NbWordsIndexesNoZero}}
				C, t[{{$i}}] = madd2({{$.V2}}[{{$j}}], {{$.V1}}[{{$i}}], t[{{$i}}], C)
			{{- end}}
		{{ end }}
		t[{{$.all.NbWords}}], D = bits.Add64(t[{{$.all.NbWords}}], C, 0)

		// m = t[0]n'[0] mod W
		m = t[0] * qInvNeg

		// -----------------------------------
		// Second loop
		C = madd0(m, q0, t[0])
		{{- range $i := $.all.NbWordsIndexesNoZero}}
				C, t[{{sub $i 1}}] = madd2(m, q{{$i}}, t[{{$i}}], C)
		{{- end}}

		 t[{{sub $.all.NbWords 1}}], C = bits.Add64(t[{{$.all.NbWords}}], C, 0)
		 t[{{$.all.NbWords}}], _ = bits.Add64(0, D, C)
	{{- end}}


	if t[{{$.all.NbWords}}] != 0 {
		// we need to reduce, we have a result on {{add 1 $.all.NbWords}} words
		{{- if gt $.all.NbWords 1}}
		var b uint64
		{{- end}}
		z[0], {{- if gt $.all.NbWords 1}}b{{- else}}_{{- end}} = bits.Sub64(t[0], q0, 0)
		{{- range $i := .all.NbWordsIndexesNoZero}}
			{{-  if eq $i $.all.NbWordsLastIndex}}
				z[{{$i}}], _ = bits.Sub64(t[{{$i}}], q{{$i}}, b)
			{{-  else  }}
				z[{{$i}}], b = bits.Sub64(t[{{$i}}], q{{$i}}, b)
			{{- end}}
		{{- end}}
		return 
	}

	// copy t into z 
	{{- range $i := $.all.NbWordsIndexesFull}}
		z[{{$i}}] = t[{{$i}}]
	{{- end}}

{{ end }}


{{ define "MUL_ARM_NOCARRY_1" }}
var {{range $i := .all.NbWordsIndexesFull}}t{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
var {{range $i := .all.NbWordsIndexesFull}}u{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
var {{range $i := .all.NbWordsIndexesFull}}v{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64

{{- range $j := $.all.NbWordsIndexesFull}}
		v{{$j}} = {{$.V2}}[{{$j}}]
	{{- end}}

{{- range $i := .all.NbWordsIndexesFull}}
{
	var c0, c1, c2 uint64
	v := {{$.V1}}[{{$i}}]
	{{- if eq $i 0}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			u{{$j}}, t{{$j}} = bits.Mul64(v, v{{$j}})
		{{- end}}
	{{- else}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			u{{$j}}, c1 = bits.Mul64(v, v{{$j}})
			{{- if eq $j 0}}
				t{{$j}}, c0 = bits.Add64(c1, t{{$j}}, 0)
			{{- else }}
				t{{$j}}, c0 = bits.Add64(c1, t{{$j}}, c0)
			{{- end}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				{{/* yes, we're tempted to write c2 = c0, but that slow the whole MUL by 20% */}}
				c2, _ = bits.Add64(0, 0, c0)
			{{- end}}
		{{- end}}
	{{- end}}

	{{- range $j := $.all.NbWordsIndexesFull}}
	{{- if eq $j 0}}
		t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, 0)
	{{- else if eq $j $.all.NbWordsLastIndex}}
		{{- if eq $i 0}}
			c2, _ = bits.Add64(u{{$j}}, 0, c0)
		{{- else}}
			c2, _ = bits.Add64(u{{$j}},c2, c0)
		{{- end}}
	{{- else }}
		t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, c0)
	{{- end}}
	{{- end}}
	
	{{- $k := $.all.NbWordsLastIndex}}

	m := qInvNeg * t0

	u0, c1 = bits.Mul64(m, q0)
	{{- range $j := $.all.NbWordsIndexesFull}}
	{{- if ne $j 0}}
		{{- if eq $j 1}}
			_, c0 = bits.Add64(t0, c1, 0)
		{{- else}}
			t{{sub $j 2}}, c0 = bits.Add64(t{{sub $j 1}}, c1, c0)
		{{- end}}
		u{{$j}}, c1 = bits.Mul64(m, q{{$j}})
	{{- end}}
	{{- end}}
	{{/* TODO @gbotrel it seems this can create a carry (c0) -- study the bounds */}}
	t{{sub $.all.NbWordsLastIndex 1}}, c0 = bits.Add64(0, c1, c0) 
	u{{$k}}, _ = bits.Add64(u{{$k}}, 0, c0)

	{{- range $j := $.all.NbWordsIndexesFull}}
		{{- if eq $j 0}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, 0)
		{{- else if eq $j $.all.NbWordsLastIndex}}
			c2, _ = bits.Add64(c2, 0, c0)
		{{- else}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, c0)
		{{- end}}
	{{- end}}

	{{- $l := sub $.all.NbWordsLastIndex 1}}
	t{{$l}}, c0 = bits.Add64(t{{$k}}, t{{$l}}, 0)
	t{{$k}}, _ = bits.Add64(u{{$k}}, c2, c0)

}
{{- end}}


{{- range $i := $.all.NbWordsIndexesFull}}
z[{{$i}}] = t{{$i}}
{{- end}}

{{ end }}


{{ define "MUL_ARM_NOCARRY_2" }}
var {{range $i := .all.NbWordsIndexesFull}}t{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
var {{range $i := .all.NbWordsIndexesFull}}u{{$i}}{{- if ne $i $.all.NbWordsLastIndex}},{{- end}}{{- end}} uint64
{{- range $i := .all.NbWordsIndexesFull}}
{
	var c0, c1, c2 uint64
	v := {{$.V1}}[{{$i}}]
	{{- if eq $i 0}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			u{{$j}}, t{{$j}} = bits.Mul64(v, {{$.V2}}[{{$j}}])
		{{- end}}
	{{- else}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			u{{$j}}, c1 = bits.Mul64(v, {{$.V2}}[{{$j}}])
			{{- if eq $j 0}}
				t{{$j}}, c0 = bits.Add64(c1, t{{$j}}, 0)
			{{- else }}
				t{{$j}}, c0 = bits.Add64(c1, t{{$j}}, c0)
			{{- end}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				{{/* yes, we're tempted to write c2 = c0, but that slow the whole MUL by 20% */}}
				c2, _ = bits.Add64(0, 0, c0)
			{{- end}}
		{{- end}}
	{{- end}}

	{{- range $j := $.all.NbWordsIndexesFull}}
	{{- if eq $j 0}}
		t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, 0)
	{{- else if eq $j $.all.NbWordsLastIndex}}
		{{- if eq $i 0}}
			c2, _ = bits.Add64(u{{$j}}, 0, c0)
		{{- else}}
			c2, _ = bits.Add64(u{{$j}},c2, c0)
		{{- end}}
	{{- else }}
		t{{add $j 1}}, c0 = bits.Add64(u{{$j}}, t{{add $j 1}}, c0)
	{{- end}}
	{{- end}}
	
	{{- $k := $.all.NbWordsLastIndex}}

	m := qInvNeg * t0

	u0, c1 = bits.Mul64(m, q0)
	{{- range $j := $.all.NbWordsIndexesFull}}
	{{- if ne $j 0}}
		{{- if eq $j 1}}
			_, c0 = bits.Add64(t0, c1, 0)
		{{- else}}
			t{{sub $j 2}}, c0 = bits.Add64(t{{sub $j 1}}, c1, c0)
		{{- end}}
		u{{$j}}, c1 = bits.Mul64(m, q{{$j}})
	{{- end}}
	{{- end}}
	{{/* TODO @gbotrel it seems this can create a carry (c0) -- study the bounds */}}
	t{{sub $.all.NbWordsLastIndex 1}}, c0 = bits.Add64(0, c1, c0) 
	u{{$k}}, _ = bits.Add64(u{{$k}}, 0, c0)

	{{- range $j := $.all.NbWordsIndexesFull}}
		{{- if eq $j 0}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, 0)
		{{- else if eq $j $.all.NbWordsLastIndex}}
			c2, _ = bits.Add64(c2, 0, c0)
		{{- else}}
			t{{$j}}, c0 = bits.Add64(u{{$j}}, t{{$j}}, c0)
		{{- end}}
	{{- end}}

	{{- $l := sub $.all.NbWordsLastIndex 1}}
	t{{$l}}, c0 = bits.Add64(t{{$k}}, t{{$l}}, 0)
	t{{$k}}, _ = bits.Add64(u{{$k}}, c2, c0)

}
{{- end}}


{{- range $i := $.all.NbWordsIndexesFull}}
z[{{$i}}] = t{{$i}}
{{- end}}

{{ end }}

{{- if ne $.NbWords 1}}

func BBB_SQUARE_NOCARRY(z,x *Element) {
	{{ template "SQUARE_NOCARRY" dict "all" . "V1" "x" "V2" "x"}}
	{{ template "reduce"  . }}
}


func BenchmarkBBB_SQUARE_NOCARRY(b *testing.B) {
	var r {{.ElementName}}
	r.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_SQUARE_NOCARRY(&r, &r)
	}
	benchRes{{.ElementName}} = r 
}


func BBB_MUL_NOCARRY(z,x,y *Element) {
	{{ template "MUL_NOCARRY" dict "all" . "V1" "x" "V2" "y"}}
	{{ template "reduce"  . }}
}

func BBB_MUL_CIOS(z,x,y *Element) {
	{{ template "MUL_CIOS" dict "all" . "V1" "x" "V2" "y"}}
	{{ template "reduce"  . }}
}

func BBB_MUL_ARM_NOCARRY_1(z,x,y *Element) {
	{{ template "MUL_ARM_NOCARRY_1" dict "all" . "V1" "x" "V2" "y"}}
}

func BBB_MUL_ARM_NOCARRY_2(z,x,y *Element) {
	{{ template "MUL_ARM_NOCARRY_2" dict "all" . "V1" "x" "V2" "y"}}
}


func BenchmarkBBB_MUL_NOCARRY(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_NOCARRY(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BenchmarkBBB_MUL_STD(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Mul(&x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BenchmarkBBB_MUL_CIOS(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_CIOS(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BenchmarkBBB_MUL_ARM_NOCARRY_1(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_ARM_NOCARRY_1(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BenchmarkBBB_MUL_ARM_NOCARRY_2(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_ARM_NOCARRY_2(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BenchmarkBBB_MUL_SOS_ROLLED(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_SOS_ROLLED(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}



func BenchmarkBBB_MUL_CIOS_ROLLED(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	var y {{.ElementName}}
	y.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BBB_MUL_CIOS_ROLLED(&x, &x, &y)
	}
	benchRes{{.ElementName}} = x 
}

func BBB_MUL_SOS_ROLLED(z, x, y *Element) {
	var t [Limbs*2 + 2]uint64

	for i := 0; i < Limbs; i++ {
		C := uint64(0)
		for j := 0; j < Limbs; j++ {
			C, t[i+j] = madd2(x[j], y[i], t[i+j], C)
		}
		t[i+Limbs] = C
	}

	for i := 0; i < Limbs; i++ {
		C := uint64(0)
		m := t[i] * qInvNeg
		for j := 0; j < Limbs; j++ {
			C, t[i+j] = madd2(qElement[j], m, t[i+j], C)
		}
		// TODO propagate carry ADD
		k := i + Limbs
		var c uint64
		t[k], c = bits.Add64(t[k], C, 0)
		k++
		for ; k < len(t); k++ {
			t[k], c = bits.Add64(t[k], 0, c)
		}
	}
	copy(z[:], t[Limbs:])
	// *z = t
	{{ template "reduce"  . }}
}

func BBB_MUL_CIOS_ROLLED(z, x, y *Element) {
	var A, C uint64
	var t Element
	for i := 0; i < Limbs; i++ {
		A, t[0] = madd1(x[0], y[i], t[0])
		m := t[0] * qInvNeg
		C = madd0(m, q0, t[0])
		for j := 1; j < Limbs; j++ {
			A, t[j] = madd2(x[j], y[i], t[j], A)
			C, t[j-1] = madd2(qElement[j], m, t[j], C)
		}
		t[Limbs-1] = C + A
	}

	*z = t
	// if z ⩾ q → z -= q
	{{ template "reduce"  . }}
}


{{- end}}

`
