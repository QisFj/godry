package main

import "testing"

//go:generate inplacegen
func Test_example2(t *testing.T) {
	// ;INPLACEGEN_(test)_FROM
	// ;2
	// ;[["Int", "int"], ["Float64", "float64"]]
	// ;*[["Add", "+"], ["Sub", "-"], ["Mul", "*"]]
	// ;	{{$1}}Func := map[string]func(v1, v2 {{$[1.2]}}) {{$[1.2]}} {
	// ;	{{range .ex1}}	"{{.v1}}": func(v1, v2 {{$[1.2]}}) {{$[1.2]}} { return v1 {{.v2}} v2},
	// ;	{{end -}}
	// ;	}
	// ;	_ = {{$1}}Func
	IntFunc := map[string]func(v1, v2 int) int{
		"Add": func(v1, v2 int) int { return v1 + v2 },
		"Sub": func(v1, v2 int) int { return v1 - v2 },
		"Mul": func(v1, v2 int) int { return v1 * v2 },
	}
	_ = IntFunc
	Float64Func := map[string]func(v1, v2 float64) float64{
		"Add": func(v1, v2 float64) float64 { return v1 + v2 },
		"Sub": func(v1, v2 float64) float64 { return v1 - v2 },
		"Mul": func(v1, v2 float64) float64 { return v1 * v2 },
	}
	_ = Float64Func
	// ;INPLACEGEN_(test)_TO
}
