package fft

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func BenchmarkExperiment1(b *testing.B) {
	var exp = NewExp[fr.Element]()

	var a, bb fr.Element
	for i := 0; i < b.N; i++ {
		exp._butterfly(&a, &bb)
	}

}

func BenchmarkExperiment2(b *testing.B) {
	var exp = NewExp[fr.Element]()

	var a, bb fr.Element
	for i := 0; i < b.N; i++ {
		exp._butterflyGeneric(&a, &bb)
	}

}
