package fft

import (
	"fmt"
	"io"
	"math/big"
	"math/bits"
	"runtime"
	"sync"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/types"
)

// Domain with a power of 2 cardinality
// compute a field element of order 2x and store it in FinerGenerator
// all other values can be derived from x, GeneratorSqrt
type Domain[E any, _ types.Element[E]] struct {
	Cardinality            uint64
	CardinalityInv         E
	Generator              E
	GeneratorInv           E
	FrMultiplicativeGen    E // generator of Fr*
	FrMultiplicativeGenInv E

	// the following slices are not serialized and are (re)computed through domain.preComputeTwiddles()

	// Twiddles factor for the FFT using Generator for each stage of the recursive FFT
	Twiddles [][]E

	// Twiddles factor for the FFT using GeneratorInv for each stage of the recursive FFT
	TwiddlesInv [][]E

	// we precompute these mostly to avoid the memory intensive bit reverse permutation in the groth16.Prover

	// CosetTable u*<1,g,..,g^(n-1)>
	CosetTable         []E
	CosetTableReversed []E // optional, this is computed on demand at the creation of the domain

	// CosetTable[i][j] = domain.Generator(i-th)SqrtInv ^ j
	CosetTableInv         []E
	CosetTableInvReversed []E // optional, this is computed on demand at the creation of the domain
}

// NewDomain returns a subgroup with a power of 2 cardinality
// cardinality >= m
func NewDomain[E any, ptE types.Element[E]](m uint64) *Domain[E, ptE] {

	var domain Domain[E, ptE]
	x := ecc.NextPowerOfTwo(m)
	domain.Cardinality = uint64(x)

	// generator of the largest 2-adic subgroup
	var rootOfUnity E

	ptE(&rootOfUnity).SetString("19103219067921713944291392827692070036145651957329286315305642004821462161904")
	const maxOrderRoot uint64 = 28
	ptE(&domain.FrMultiplicativeGen).SetUint64(5)

	ptE(&domain.FrMultiplicativeGenInv).Inverse(&domain.FrMultiplicativeGen)

	// find generator for Z/2^(log(m))Z
	logx := uint64(bits.TrailingZeros64(x))
	if logx > maxOrderRoot {
		panic(fmt.Sprintf("m (%d) is too big: the required root of unity does not exist", m))
	}

	// Generator = FinerGenerator^2 has order x
	expo := uint64(1 << (maxOrderRoot - logx))
	ptE(&domain.Generator).Exp(rootOfUnity, big.NewInt(int64(expo))) // order x
	ptE(&domain.GeneratorInv).Inverse(&domain.Generator)
	ptE(&domain.CardinalityInv).SetUint64(uint64(x))
	ptE(&domain.CardinalityInv).Inverse(&domain.CardinalityInv)

	// twiddle factors
	domain.preComputeTwiddles()

	// store the bit reversed coset tables
	domain.reverseCosetTables()

	return &domain
}

func (d *Domain[E, ptE]) reverseCosetTables() {
	d.CosetTableReversed = make([]E, d.Cardinality)
	d.CosetTableInvReversed = make([]E, d.Cardinality)
	copy(d.CosetTableReversed, d.CosetTable)
	copy(d.CosetTableInvReversed, d.CosetTableInv)
	BitReverse(d.CosetTableReversed)
	BitReverse(d.CosetTableInvReversed)
}

func (d *Domain[E, ptE]) preComputeTwiddles() {
	// nb fft stages
	nbStages := uint64(bits.TrailingZeros64(d.Cardinality))

	d.Twiddles = make([][]E, nbStages)
	d.TwiddlesInv = make([][]E, nbStages)
	d.CosetTable = make([]E, d.Cardinality)
	d.CosetTableInv = make([]E, d.Cardinality)

	var wg sync.WaitGroup

	// for each fft stage, we pre compute the twiddle factors
	twiddles := func(t [][]E, omega E) {
		for i := uint64(0); i < nbStages; i++ {
			t[i] = make([]E, 1+(1<<(nbStages-i-1)))
			var w E
			if i == 0 {
				w = omega
			} else {
				w = t[i-1][2]
			}
			ptE(&t[i][0]).SetOne()
			t[i][1] = w
			for j := 2; j < len(t[i]); j++ {
				ptE(&t[i][j]).Mul(&t[i][j-1], &w)
			}
		}
		wg.Done()
	}

	expTable := func(sqrt E, t []E) {
		ptE(&t[0]).SetOne()
		precomputeExpTable[E, ptE](sqrt, t)
		wg.Done()
	}

	wg.Add(4)
	go twiddles(d.Twiddles, d.Generator)
	go twiddles(d.TwiddlesInv, d.GeneratorInv)
	go expTable(d.FrMultiplicativeGen, d.CosetTable)
	go expTable(d.FrMultiplicativeGenInv, d.CosetTableInv)

	wg.Wait()

}

func precomputeExpTable[E any, ptE types.Element[E]](w E, table []E) {
	n := len(table)

	// see if it makes sense to parallelize exp tables pre-computation
	interval := 0
	if runtime.NumCPU() >= 4 {
		interval = (n - 1) / (runtime.NumCPU() / 4)
	}

	// this ratio roughly correspond to the number of multiplication one can do in place of a Exp operation
	const ratioExpMul = 6000 / 17

	if interval < ratioExpMul {
		precomputeExpTableChunk[E, ptE](w, 1, table[1:])
		return
	}

	// we parallelize
	var wg sync.WaitGroup
	for i := 1; i < n; i += interval {
		start := i
		end := i + interval
		if end > n {
			end = n
		}
		wg.Add(1)
		go func() {
			precomputeExpTableChunk[E, ptE](w, uint64(start), table[start:end])
			wg.Done()
		}()
	}
	wg.Wait()
}

func precomputeExpTableChunk[E any, ptE types.Element[E]](w E, power uint64, table []E) {

	// this condition ensures that creating a domain of size 1 with cosets don't fail
	if len(table) > 0 {
		ptE(&table[0]).Exp(w, new(big.Int).SetUint64(power))
		for i := 1; i < len(table); i++ {
			ptE(&table[i]).Mul(&table[i-1], &w)
		}
	}
}

// WriteTo writes a binary representation of the domain (without the precomputed twiddle factors)
// to the provided writer
func (d *Domain[E, ptE]) WriteTo(w io.Writer) (int64, error) {

	panic("not implemented")
}

// ReadFrom attempts to decode a domain from Reader
func (d *Domain[E, ptE]) ReadFrom(r io.Reader) (int64, error) {
	panic("not implemented")
}

// BitReverse applies the bit-reversal permutation to a.
// len(a) must be a power of 2 (as in every single function in this file)
func BitReverse[E any](a []E) {
	n := uint64(len(a))
	nn := uint64(64 - bits.TrailingZeros64(n))

	for i := uint64(0); i < n; i++ {
		irev := bits.Reverse64(i) >> nn
		if irev > i {
			a[i], a[irev] = a[irev], a[i]
		}
	}
}
