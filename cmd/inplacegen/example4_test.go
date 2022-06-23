package main

import "testing"

//go:generate inplacegen
func Test_example4(t *testing.T) {
	// INPLACEGEN_(test)_FROM
	// import inplacegen_import.txt
	// should same as example2 # this line not used, but would be keep
	IntFunc := map[string]func(v1, v2 int) int {
		"Add": func(v1, v2 int) int { return v1 + v2},
		"Sub": func(v1, v2 int) int { return v1 - v2},
		"Mul": func(v1, v2 int) int { return v1 * v2},
	}
	_ = IntFunc
	Float64Func := map[string]func(v1, v2 float64) float64 {
		"Add": func(v1, v2 float64) float64 { return v1 + v2},
		"Sub": func(v1, v2 float64) float64 { return v1 - v2},
		"Mul": func(v1, v2 float64) float64 { return v1 * v2},
	}
	_ = Float64Func
	// INPLACEGEN_(test)_TO
}

