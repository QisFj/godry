package gen

import (
	"flag"
	"strings"
)

var usualTypes = []string{
	"string", "int", "uint",
}

func ParseFlagTypes() []string {
	pTypes := flag.String("types", "", "types, split by comma")
	flag.Parse()
	if *pTypes != "" {
		return strings.Split(*pTypes, ",")
	}
	return usualTypes
}
