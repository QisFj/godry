package main

import "testing"

//go:generate inplacegen
func Test_main(t *testing.T) {
	// INPLACEGEN_(test)_FROM
	// 2
	// [["Int", "int"], ["Float64", "float64"]]
	// [["Add", "+"], ["Sub", "-"], ["Mul", "*"]]
	// 	{{$1}}{{$2}} := func (v1, v2 {{$[1.2]}}) {{$[1.2]}} {
	// 		return v1 {{$[2.2]}} v2
	// 	}
	// 	_ = {{$1}}{{$2}}
	IntAdd := func(v1, v2 int) int {
		return v1 + v2
	}
	_ = IntAdd
	IntSub := func(v1, v2 int) int {
		return v1 - v2
	}
	_ = IntSub
	IntMul := func(v1, v2 int) int {
		return v1 * v2
	}
	_ = IntMul
	Float64Add := func(v1, v2 float64) float64 {
		return v1 + v2
	}
	_ = Float64Add
	Float64Sub := func(v1, v2 float64) float64 {
		return v1 - v2
	}
	_ = Float64Sub
	Float64Mul := func(v1, v2 float64) float64 {
		return v1 * v2
	}
	_ = Float64Mul
	// INPLACEGEN_(test)_TO
}
