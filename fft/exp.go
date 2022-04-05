package fft

import (
	"github.com/consensys/gnark-crypto/types"
)

type Exp[E any, _ types.Element[E]] struct {
}

func (e *Exp[E, ptE]) _butterflyGeneric(a, b *E) {
	t := *a
	ptE(a).Add(a, b)
	ptE(b).Sub(&t, b)
}

func (e *Exp[E, ptE]) _butterfly(a, b *E) {
	ptE(a).Butterfly(b)
}

func NewExp[E any, ptE types.Element[E]]() *Exp[E, ptE] {

	var exp Exp[E, ptE]
	return &exp
}
