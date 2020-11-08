# cexp

Get value by condition. It likes `:?`.

```go
package main

import (
	"log"

	"github.com/QisFj/godry/cexp"
)

func main(){
    // won't panic, would print "false", 
    log.Print(cexp.StringShortCircuit(false,
    	func() string {
    	    panic("! not short circuit !")
    	},
    	func() string {
    		return "false"
    	},
    ))
}
```