// The following directive is necessary to make the package coherent:

//go:build ignore
// +build ignore

// This program generates:
// - cexps.go
// It can be invoked by running go generate
package main

import (
	"log"
	"strings"

	"github.com/QisFj/godry/gen"
)

func main() {
	types := gen.ParseFlagTypes()
	log.Printf("types: \n- '%s'\n", strings.Join(types, "'\n- '"))
	gen.Gen("cexps", T, types)
}

const T = `// Code generated by go generate; DO NOT EDIT.
package cexp

<< range $t := .>>
type <<$t | Title>>Getter func() <<$t>>

func <<$t | Title>>(condition bool, v1, v2 <<$t>>) <<$t>> {
	if condition {
		return v1
	}
	return v2
}

func <<$t | Title>>ShortCircuit(condition bool, g1, g2 <<$t | Title>>Getter) <<$t>> {
	if condition {
		return g1()
	}
	return g2()
}
<< end >>
`
