# Package slice

do something for a slice.

## Usage

```go
package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/QisFj/godry/slice"
)

func main() {
	is := []int{1, 2, 3, 4, 5, 6}
	ss := slice.MapString(is, func(i int, v interface{}) string {
		return strconv.Itoa(v.(int))
	})
	log.Printf("after map: %#v\n", ss)
	slice.Filter(&ss, func(index int) bool {
		return strings.ContainsAny(ss[index], "135")
	})
	log.Printf("after filter: %#v\n", ss)
}
```

## Map

Map a slice to another slice.

## Filter

Filter some elements of slice.
