// Copyright 2022 ConsenSys Software Inc.
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

package arm64

import (
	"fmt"
	"github.com/consensys/bavard"
	"github.com/consensys/bavard/arm64"
	"strconv"
	"strings"
)

// Registers used: 2 * nbWords
func (f *FFArm64) generateAdd() {
	f.Comment("add(res, x, y *Element)")
	//stackSize := f.StackSize(f.NbWords*2, 0, 0)
	registers := f.FnHeader("add", 0, 24)
	defer f.AssertCleanStack(0, 0)

	// registers
	z := registers.PopN(f.NbWords)
	xPtr := registers.Pop()
	yPtr := registers.Pop()
	ops := registers.PopN(2)

	f.LDP("x+8(FP)", xPtr, yPtr)
	f.Comment("load operands and add mod 2^r")

	op0 := f.ADDS
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.RegisterOffset(xPtr, 8*i), z[i], ops[0])
		f.LDP(f.RegisterOffset(yPtr, 8*i), z[i+1], ops[1])

		op0(z[i], z[i+1], z[i])
		op0 = f.ADCS

		f.ADCS(ops[0], ops[1], z[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.RegisterOffset(xPtr, 8*i), z[i], "can't import these in pairs")
		f.MOVD(f.RegisterOffset(yPtr, 8*i), ops[0])
		op0(z[i], ops[0], z[i])
	}
	registers.Push(xPtr, yPtr)
	registers.Push(ops...)

	t := registers.PopN(f.NbWords)
	f.reduce(z, t)
	registers.Push(t...)

	f.Comment("store")
	zPtr := registers.Pop()
	f.MOVD("res+0(FP)", zPtr)
	f.storeVector(z, zPtr)

	f.RET()

}

func (f *FFArm64) generateDouble() {
	f.Comment("double(res, x *Element)")
	registers := f.FnHeader("double", 0, 16)
	defer f.AssertCleanStack(0, 0)

	// registers
	z := registers.PopN(f.NbWords)
	xPtr := registers.Pop()
	zPtr := registers.Pop()
	//ops := registers.PopN(2)

	f.LDP("res+0(FP)", zPtr, xPtr)
	f.Comment("load operands and add mod 2^r")

	op0 := f.ADDS
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.RegisterOffset(xPtr, 8*i), z[i], z[i+1])

		op0(z[i], z[i], z[i])
		op0 = f.ADCS

		f.ADCS(z[i+1], z[i+1], z[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.RegisterOffset(xPtr, 8*i), z[i])
		op0(z[i], z[i], z[i])
	}
	registers.Push(xPtr)

	t := registers.PopN(f.NbWords)
	f.reduce(z, t)
	registers.Push(t...)

	f.Comment("store")
	f.storeVector(z, zPtr)

	f.RET()

}

// generateSub NO LONGER uses one more register than generateAdd, but that's okay since we have 29 registers available.
func (f *FFArm64) generateSub() {
	f.Comment("sub(res, x, y *Element)")

	registers := f.FnHeader("sub", 0, 24)
	defer f.AssertCleanStack(0, 0)

	// registers
	z := registers.PopN(f.NbWords)
	xPtr := registers.Pop()
	yPtr := registers.Pop()
	ops := registers.PopN(2)

	f.LDP("x+8(FP)", xPtr, yPtr)
	f.Comment("load operands and subtract mod 2^r")

	op0 := f.SUBS
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.RegisterOffset(xPtr, 8*i), z[i], ops[0])
		f.LDP(f.RegisterOffset(yPtr, 8*i), z[i+1], ops[1])

		op0(z[i+1], z[i], z[i])
		op0 = f.SBCS

		f.SBCS(ops[1], ops[0], z[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.RegisterOffset(xPtr, 8*i), z[i], "can't import these in pairs")
		f.MOVD(f.RegisterOffset(yPtr, 8*i), ops[0])
		op0(ops[0], z[i], z[i])
	}
	registers.Push(xPtr, yPtr)
	registers.Push(ops...)

	f.Comment("load modulus and select")

	t := registers.PopN(f.NbWords)
	zero := registers.Pop()
	f.MOVD(0, zero)

	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.GlobalOffset("q", 8*i), t[i], t[i+1])

		f.CSEL("CS", zero, t[i], t[i])
		f.CSEL("CS", zero, t[i+1], t[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.GlobalOffset("q", 8*i), t[i])

		f.CSEL("CS", zero, t[i], t[i])
	}

	registers.Push(zero)

	f.Comment("augment (or not)")

	op0 = f.ADDS
	for i := 0; i < f.NbWords; i++ {
		op0(z[i], t[i], z[i])
		op0 = f.ADCS
	}

	registers.Push(t...)

	f.Comment("store")
	zPtr := registers.Pop()
	f.MOVD("res+0(FP)", zPtr)
	f.storeVector(z, zPtr)

	f.RET()

}

func (f *FFArm64) generateNeg() {
	f.Comment("neg(res, x *Element)")
	registers := f.FnHeader("neg", 0, 16)
	defer f.AssertCleanStack(0, 0)

	// registers
	z := registers.PopN(f.NbWords)
	xPtr := registers.Pop()
	zPtr := registers.Pop()
	ops := registers.PopN(2)
	xNotZero := registers.Pop()

	f.LDP("res+0(FP)", zPtr, xPtr)
	f.Comment("load operands and subtract")

	f.MOVD(0, xNotZero)
	op0 := f.SUBS
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.RegisterOffset(xPtr, 8*i), z[i], z[i+1])
		f.LDP(f.GlobalOffset("q", 8*i), ops[0], ops[1])

		f.ORR(z[i], xNotZero, xNotZero, "has x been 0 so far?")
		f.ORR(z[i+1], xNotZero, xNotZero)

		op0(z[i], ops[0], z[i])
		op0 = f.SBCS

		f.SBCS(z[i+1], ops[1], z[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.RegisterOffset(xPtr, 8*i), z[i], "can't import these in pairs")
		f.MOVD(f.GlobalOffset("q", 8*i), ops[0])

		f.ORR(z[i], xNotZero, xNotZero)

		op0(z[i], ops[0], z[i])
	}

	registers.Push(xPtr)
	registers.Push(ops...)

	f.TST(-1, xNotZero)
	for i := 0; i < f.NbWords; i++ {
		f.CSEL("EQ", xNotZero, z[i], z[i])
	}

	f.Comment("store")
	f.storeVector(z, zPtr)

	f.RET()

}

// Needs 3 * nbWords + 8 registers
func (f *FFArm64) generateMul() {
	f.generateMadd()
	f.generateLoadVector()
	defer f.AssertCleanStack(0, 0)

	registers := f.FnHeader("mul", 0, 24)
	variables := newVariables(f, &registers)
	defer variables.undefAll()

	f.Comment("mul(res, x, y)")

	x := registers.PopN(2) //TODO: Always fetch for the next iteration? Will need 4
	xPtr := registers.Pop()
	yPtr := registers.Pop()
	qInv0 := variables.def("_qInv0") //registers.Pop()
	c0 := variables.def("c0")        //registers.Pop()
	c1 := variables.def("c1")        //registers.Pop()
	c2 := variables.def("c2")        //registers.Pop()
	//c3 := registers.Pop() //Need an extra one for madd0
	m := variables.def("m") //registers.Pop()

	q := variables.defN("q", f.NbWords) //registers.PopN(f.NbWords)
	y := variables.defN("y", f.NbWords) //registers.PopN(f.NbWords)

	f.Comment("Load all of y")
	f.LDP("x+8(FP)", xPtr, yPtr)
	f.loadVector(yPtr, y)
	registers.Push(yPtr)
	z := variables.defN("z", f.NbWords) //registers.PopN(f.NbWords)

	f.MOVD(f.GlobalOffset("qInv0", 0), qInv0, "Load qInv0")
	f.Comment("Load q")
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.GlobalOffset("q", 8*i), q[i], q[i+1] /*, fmt.Sprintf("%s, %s = q[%d], q[%d]", q[i].Name(), q[i+1].Name(), i, i+1)*/)
	}
	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.GlobalOffset("q", 8*i), q[i] /*, fmt.Sprintf("%s = q[%d]", q[i].Name(), i)*/)
	}

	for j := 0; j < f.NbWords; j++ {

		if j%2 == 0 {
			if j+1 < f.NbWords {
				f.LDP(f.RegisterOffset(xPtr, 8*j), x[0], x[1] /*, fmt.Sprintf("%s, %s = x[%d], x[%d]", x[0].Name(), x[1].Name(), j, j+1)*/)
			} else {
				f.MOVD(f.RegisterOffset(xPtr, 8*j), x[0] /*, fmt.Sprintf("%s = x[%d]", x[0].Name(), j)*/)
			}
		}

		v := x[j%2]

		f.Comment("Round " + strconv.Itoa(j))

		if j == 0 {
			f.MUL(v, y[0], c0)
			f.UMULH(v, y[0], c1)
		} else {
			f.madd1(c1, c0, v, y[0], z[0])
		}

		f.MUL(qInv0, c0, m)
		f.madd0(c2, m, q[0], c0)

		for i := 1; i < f.NbWords; i++ {
			if j == 0 {
				f.madd1(c1, c0, v, y[i], c1)
			} else {
				f.madd2(c1, c0, v, y[i], c1, z[i])
			}

			if i+1 == f.NbWords {
				f.madd3(z[i], z[i-1], m, q[i], c0, c2, c1)
			} else {
				f.madd2(c2, z[i], m, q[i], c2, c0)
			}
		}
	}

	f.Comment("Reduce if necessary")
	f.SUBS(q[0], z[0], y[0])
	for i := 1; i < f.NbWords; i++ {
		f.SBCS(q[i], z[i], y[i])
	}
	for i := 0; i < f.NbWords; i++ {
		f.CSEL("CS", y[i], z[i], z[i])
	}

	registers.Push(xPtr)
	zPtr := registers.Pop()
	f.MOVD("res+0(FP)", zPtr, "zPtr")
	f.storeVector(z, zPtr)

	f.RET()
}

// MACROS?
func (f *FFArm64) _generateLoadOrStoreVector(name string, instruction string, addrAndRegToSourceAndDst func(addr string, reg string) (string, string)) {
	f.Write("#define ")
	f.Write(name)
	f.Write("(ePtr, ")

	var i int
	names := make([]string, f.NbWords)
	for i = 0; i < f.NbWords; i++ {
		names[i] = "e" + strconv.Itoa(i)
	}
	f.Write(strings.Join(names, ", "))
	f.WriteLn(")\\")

	for i = 0; i < f.NbWords-1; i += 2 {
		f.Write("\t")
		f.Write(instruction)
		f.Write(" ")

		addr := fmt.Sprintf("%d(ePtr)", 8*i)
		regs := fmt.Sprintf("(%s, %s)", names[i], names[i+1])

		src, dst := addrAndRegToSourceAndDst(addr, regs)

		f.Write(src)
		f.Write(", ")
		f.Write(dst)
		f.WriteLn("\\")
	}

	if f.NbWords%2 == 1 {
		i = f.NbWords - 1
		src, dst := addrAndRegToSourceAndDst(fmt.Sprintf("%d(ePtr)", 8*i), names[i])
		f.Write("\tMOVD ")
		f.Write(src)
		f.Write(", ")
		f.Write(dst)
		f.WriteLn("\\")
	}
}
func (f *FFArm64) generateLoadVector() {

	f._generateLoadOrStoreVector("loadVector", "LDP", func(addr string, reg string) (string, string) {
		return addr, reg
	})

	/*	f.WriteLn("#define loadVector(ePtr, ")

		var i int
		names := make([]string, f.NbWords)
		for i = 0; i < f.NbWords; i++ {
			names[i] = "e" + strconv.Itoa(i)
			f.Write(names[i])
		}
		f.Write(strings.Join(names, ", "))
		f.WriteLn(")\\")

		const ePtrComma = "(ePtr), "
		for i = 0; i < f.NbWords; i += 2 {
			f.Write("\tLDP ")
			f.Write(strconv.Itoa(8 * i))
			f.Write(ePtrComma)
			f.Write(names[i])
			f.Write(", ")
			f.Write(names[i+1])
			f.WriteLn("\\")
		}

		if f.NbWords%2 == 1 {
			i = f.NbWords - 1
			f.Write("\tMOVD ")
			f.Write(strconv.Itoa(8 * i))
			f.Write(ePtrComma)
			f.Write(names[i])
			f.WriteLn("\\")
		}*/
}

func (f *FFArm64) generateStoreVector() {
	f._generateLoadOrStoreVector("storeVector", "STP", func(addr string, reg string) (string, string) {
		return reg, addr
	})
}

func (f *FFArm64) generateMadd() {
	f.WriteLn(
		`
// (hi, -) = a*b + c
#define madd0(hi, a, b, c) \
    madd1(hi, hi, a, b, c)\

// (hi, lo) = a*b + c
// hi can be the same register as any other operand, including lo
// lo can't be the same register as any of the input
#define madd1(hi, lo, a, b, c) \
    MUL a, b, lo\
	ADDS c, lo, lo\
    UMULH a, b, hi\
    ADC $0, hi, hi\

// madd2 (hi, lo) = a*b + c + d
#define madd2(hi, lo, a, b, c, d) \
    madd3(a, b, c, d, $0, hi, lo)\

//madd3 (hi, lo) = a*b + c + d + (e,0)
#define madd3(hi, lo, a, b, c, d, e) \
    MUL a, b, lo\
    UMULH a, b, hi\
    ADDS c, lo, lo\
    ADC $0, hi, hi\
    ADDS d, lo, lo\
    ADC e, hi, hi\
`)
}

// madd0 (hi, -) = a*b + c
func (f *FFArm64) madd0(hi, a, b, c interface{}) {
	f.callTemplate("madd0", hi, a, b, c)
}

// madd1 (hi, lo) = a*b + c
func (f *FFArm64) madd1(hi, lo, a, b, c interface{}) {
	f.callTemplate("madd1", hi, lo, a, b, c)
}

// madd2 (hi, lo) = a*b + c + d
func (f *FFArm64) madd2(hi, lo, a, b, c, d interface{}) {
	f.callTemplate("madd2", hi, lo, a, b, c, d)
}

// madd3 (hi, lo) = a*b + c + d + (e,0)
func (f *FFArm64) madd3(hi, lo, a, b, c, d, e interface{}) {
	f.callTemplate("madd3", hi, lo, a, b, c, d, e)
}

func toInterfaceSlice(first interface{}, rest interface{}) []interface{} {
	restSlice, err := bavard.AssertSlice(rest)

	if err != nil {
		panic("not a slice")
	}

	res := make([]interface{}, restSlice.Len()+1)
	res[0] = first
	for i := 0; i < restSlice.Len(); i++ {
		res[i+1] = restSlice.Index(i).Interface()
	}
	return res
}

func (f *FFArm64) loadVector(vectorHeadPtr interface{}, vector interface{}) {
	f.callTemplate("loadVector", toInterfaceSlice(vectorHeadPtr, vector)...)
}

func (f *FFArm64) storeVector(vector interface{}, baseAddress arm64.Register) {
	f.callTemplate("storeVector", toInterfaceSlice(baseAddress, vector)...)
	/*	for i := 0; i < f.NbWords-1; i += 2 {
			f.STP(vector[i], vector[i+1], f.RegisterOffset(baseAddress, 8*i))
		}

		if f.NbWords%2 == 1 {
			i := f.NbWords - 1
			f.MOVD(vector[i], f.RegisterOffset(baseAddress, 8*i))
		}*/
}

func (f *FFArm64) callTemplate(templateName string, ops ...interface{}) {
	f.Write(templateName)
	f.Write("(")
	for i := 0; i < len(ops); i++ {
		f.Write(arm64.Operand(ops[i]))
		if i+1 < len(ops) {
			f.Write(", ")
		}
	}
	f.WriteLn(")")
}

//TODO: Put it in a macro
func (f *FFArm64) reduce(z, t []arm64.Register) {

	if len(z) != f.NbWords || len(t) != f.NbWords {
		panic("need 2*nbWords registers")
	}

	f.Comment("load modulus and subtract")

	op0 := f.SUBS
	for i := 0; i < f.NbWords-1; i += 2 {
		f.LDP(f.GlobalOffset("q", 8*i), t[i], t[i+1])

		op0(t[i], z[i], t[i])
		op0 = f.SBCS

		f.SBCS(t[i+1], z[i+1], t[i+1])
	}

	if f.NbWords%2 == 1 {
		i := f.NbWords - 1
		f.MOVD(f.GlobalOffset("q", 8*i), t[i])

		op0(t[i], z[i], t[i])
	}

	f.Comment("reduce if necessary")

	for i := 0; i < f.NbWords; i++ {
		f.CSEL("CS", t[i], z[i], z[i])
	}
}

// </macros>
