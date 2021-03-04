package main

import (
	"testing"
)

//go:generate inplacegen -name=string
func Test_example3(t *testing.T) {
	// ;INPLACEGEN_(ignore)_FROM
	// ;1
	// ;[["A"], ["B"], ["C"]]
	// ;{{$1}}
	// ;INPLACEGEN_(ignore)_TO
	_ = `
INPLACEGEN_(string)_FROM
1
[["a"], ["b"], ["c"]]
{{$1}}
a
b
c
INPLACEGEN_(string)_TO
`
}
