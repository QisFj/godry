// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates:
// - maps.go
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
	gen.Gen("maps", T, types)
}

const T = `// Code generated by go generate; DO NOT EDIT.
package slice

<< range $t := .>>type << $t | Title >>MapFunc func(i int, v interface{}) << $t >>
<< end >>
<< range $t := .>>func Map<< $t | Title >>(slice interface{}, f << $t | Title >>MapFunc) []<< $t >> {
	sv := valueOf(slice)
	if sv.IsNil() {
		return nil
	}
	list := make([]<<$t>>, 0, sv.Len())
	foreach(sv, func(i int, v interface{}) {
		list = append(list, f(i, v))
	})
	return list
}
<< end >>`
