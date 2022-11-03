package element

const MulNoCarryARM64 = `
{{ define "mul_nocarry" }}
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

`
