package element

// see MulDoc for general documentation on the Multiplication algorithm.
// gnark-crypto will generate couple of variation for this (important) algorithm.
//
// If len(modulus) == 1 word
// 		the algorithm is straightforward (mul_one_limb) and is the same for all targets,
//		and all moduli.
// Else (see https://hackmd.io/@gnark/modular_multiplication for more details)
// 		by default, gnark-crypto generates the CIOS multiplication "mul_cios", as defined in
// 		https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf section 2.3.2
// 		let boundMultiply = 2**63 - 2 (highest bit and lowest bit unset)
// 		let boundSquare = 2**62 - 1 (2 highest bits unset)
// 		if modulus[lastWord] <= boundMultiply
// 			generate optimized CIOS: "mul_no_carry"
//			on x86 architecture, if BMI,ADX instruction are available (most recent machine) this
// 			algorithm is generated in assembly
//			on arm64 architecture, the same algorithm such as the Go compiler produces better code for ARM. ("mul_no_carry_arm64")
//		if modulus[lastWord] <= boundSquare
//			if target == arm64
//				generate optimized Squaring CIOS: "square_no_carry_arm64" (in all other cases, Square calls the mul impl)

const MulDoc = `
{{define "mul_doc noCarry"}}
// Implements CIOS multiplication -- section 2.3.2 of Tolga Acar's thesis
// https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf
// 
// The algorithm:
// 
// for i=0 to N-1
// 		C := 0
// 		for j=0 to N-1
// 			(C,t[j]) := t[j] + x[j]*y[i] + C
// 		(t[N+1],t[N]) := t[N] + C
//		
// 		C := 0
// 		m := t[0]*q'[0] mod D
// 		(C,_) := t[0] + m*q[0]
// 		for j=1 to N-1
// 			(C,t[j-1]) := t[j] + m*q[j] + C
//		
// 		(C,t[N-1]) := t[N] + C
// 		t[N] := t[N+1] + C
//
// → N is the number of machine words needed to store the modulus q
// → D is the word size. For example, on a 64-bit architecture D is 2	64
// → x[i], y[i], q[i] is the ith word of the numbers x,y,q
// → q'[0] is the lowest word of the number -q⁻¹ mod r. This quantity is pre-computed, as it does not depend on the inputs.
// → t is a temporary array of size N+2 
// → C, S are machine words. A pair (C,S) refers to (hi-bits, lo-bits) of a two-word number
{{- if .noCarry}}
// 
// As described here https://hackmd.io/@gnark/modular_multiplication we can get rid of one carry chain and simplify:
//
// for i=0 to N-1
// 		(A,t[0]) := t[0] + x[0]*y[i] 
// 		m := t[0]*q'[0] mod W
// 		C,_ := t[0] + m*q[0]
// 		for j=1 to N-1
// 			(A,t[j])  := t[j] + x[j]*y[i] + A
// 			(C,t[j-1]) := t[j] + m*q[j] + C
//		
// 		t[N-1] = C + A
//
// This optimization saves 5N + 2 additions in the algorithm, and can be used whenever the highest bit 
// of the modulus is zero (and not all of the remaining bits are set).
{{- end}}
{{ end }}
`

const Mul = `

import (
	"math/bits"
)

// Mul z = x * y (mod q)
{{- if $.NoCarry}}
//
// x and y must be strictly inferior to q
{{- end }}
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	{{ mul_doc $.NoCarry }}

	{{- if eq $.NbWords 1}}
		{{ template "mul_cios_one_limb" dict "all" . "V1" "x" "V2" "y" }}
	{{- else }}
		{{- if .NoCarry}}
			{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
			{{ template "reduce"  . }}
		{{- else }}
			{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" }}
			{{ template "reduce"  . }}
		{{- end }}
	{{- end }}
	return z
}

// Square z = x * x (mod q)
{{- if $.NoCarry}}
//
// x must be strictly inferior to q
{{- end }}
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	// see Mul for algorithm documentation
	{{- if eq $.NbWords 1}}
		{{ template "mul_cios_one_limb" dict "all" . "V1" "x" "V2" "x" }}
	{{- else }}
		{{- if .NoCarrySquare}}
			{{ template "square_nocarry" dict "all" . "V1" "x" "V2" "x"}}
			{{ template "reduce"  . }}
		{{- else }}
			{{- if .NoCarry}}
				{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x" }}
			{{- else }}
				{{ template "mul_cios" dict "all" . "V1" "x" "V2" "x" }}
			{{- end}}
			{{ template "reduce"  . }}
		{{- end }}
	{{- end }}
	return z
}

`
