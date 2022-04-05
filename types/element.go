package types

import "math/big"

type Element[T any] interface {
	SetInt64(int64) *T
	SetUint64(uint64) *T
	SetOne() *T
	SetString(string) *T
	SetInterface(i1 interface{}) (*T, error)
	Exp(T, *big.Int) *T
	Inverse(*T) *T
	Neg(*T) *T
	Double(*T) *T
	Mul(*T, *T) *T
	Add(*T, *T) *T
	Sub(*T, *T) *T
	Div(*T, *T) *T
	Butterfly(*T)
	BitLen() int
	FromMont() *T
	Bit(i uint64) uint64
	Marshal() []byte
	IsUint64() bool
	Uint64() uint64

	ToBigIntRegular(res *big.Int) *big.Int

	IsZero() bool
	IsOne() bool

	Equal(*T) bool
	String() string

	*T
}
