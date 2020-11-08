# Package set

set.

## Usage

```go
package main

import (
	"log"

	"github.com/QisFj/godry/set"
)

func main() {
	log.Println("Int Set:")
	intSet()
	log.Println("String Set:")
	stringSet()
}

func intSet() {
	s := set.Int{}
	s.Add(1, 1, 1, 2, 2, 2)
	log.Printf("set list: %v\n", s.List())
}

func stringSet() {
	s := set.String{}
	s.Add("1", "1", "1", "2", "2", "2")
	log.Printf("set list: %v\n", s.List())
}
```