# assert

```go
package main

import (
	"log"

	"github.com/QisFj/godry/assert"
)

func main(){
    defer assert.Catch(func(violation assert.Violation) {
		// Catch Violation: <assertion violation>: xxx.go:123 1 == 2
        log.Printf("Catch Violation: %s", violation)
	})
    assert.Assert(1 == 2, "1 == 2", nil)
}
```