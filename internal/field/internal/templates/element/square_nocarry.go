package element

// on x86 SquareNoCarry == mul
const SquareNoCarry = `
{{ define "square_nocarry" }}
	{{ template "mul_nocarry" dict "all" $.all "V1" $.V1 "V2" $.V2}}
{{ end}}
`
