package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/QisFj/godry/slice"
)

func main() {
	is := []int{1, 2, 3, 4, 5, 6}
	ss := slice.Map(is, func(i int, v int) string {
		return strconv.Itoa(v)
	})
	log.Printf("after map: %#v\n", ss)
	ss = slice.Filter(ss, func(_ int, v string) bool {
		return strings.ContainsAny(v, "135")
	})
	log.Printf("after filter: %#v\n", ss)
}
