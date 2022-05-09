package arm64

import (
	"github.com/consensys/bavard/arm64"
	"strconv"
)

// Variables Variable manager
type Variables struct {
	names     map[string]arm64.Register
	registers *arm64.Registers
	f         *FFArm64
}

func newVariables(f *FFArm64, registers *arm64.Registers) Variables {
	return Variables{f: f, registers: registers, names: make(map[string]arm64.Register)}
}

func (v *Variables) def(name string) string {
	r := v.registers.Pop()
	if _, existing := v.names[name]; existing {
		panic("variable already exists")
	}
	v.names[name] = r
	v.f.Write("#define ")
	v.f.Write(name)
	v.f.Write(" ")
	v.f.WriteLn(r.Name())

	return name
}

func (v *Variables) defN(name string, count int) []string {
	res := make([]string, count)

	for i := 0; i < count; i++ {
		res[i] = v.def(name + strconv.Itoa(i))
	}

	return res
}

func (v *Variables) undef(names ...string) {
	for _, name := range names {
		r, exists := v.names[name]
		if !exists {
			panic("no such variable")
		}
		v.registers.Push(r)
		delete(v.names, name)
		v.f.WriteLn("#undef " + name)
	}

	//v.f.WriteLn("#undef " + strings.Join(names, ", "))
}

func keys(m map[string]arm64.Register) []string {
	_keys := make([]string, 0, len(m))
	for key := range m {
		_keys = append(_keys, key)
	}
	return _keys
}

func (v *Variables) undefAll() {
	v.undef(keys(v.names)...)
}
